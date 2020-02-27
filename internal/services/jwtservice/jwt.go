package jwtservice

import (
	"context"
	"crypto/rsa"
	"fmt"

	"io/ioutil"
	"time"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/tokenmodels"
	"github.com/tjsampson/token-svc/pkg/metrics"

	"github.com/dgrijalva/jwt-go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// TokenType represents the different types of tokens the jwtservice can mint
type TokenType int

const (
	// Unknown token type
	Unknown TokenType = iota
	// Access token type
	Access
	// Refresh token type
	Refresh
	// EmailValidation token type
	EmailValidation
)

const (
	auditEventJWTError      = "jwt-error"
	auditEventJWTValidation = "jwt-validation"
)

type (
	customClaims struct {
		Roles []string `json:"roles"`
		Name  string   `json:"name"`
	}

	accessTokenClaims struct {
		*jwt.StandardClaims
		customClaims
	}

	refreshTokenClaims struct {
		*jwt.StandardClaims
	}
)

// Provider is the JWT client provider
type Provider interface {
	GenerateAccessToken(ctx context.Context, aTokenChan chan tokenmodels.TokenResult, tokenData map[string]interface{})
	GenerateRefreshToken(ctx context.Context, rTokenChan chan tokenmodels.TokenResult, tokenData map[string]interface{})
	IsValidAccessToken(ctx context.Context, tkn string) (*accessTokenClaims, bool)
}

type provider struct {
	cfg         *config.Config
	signBytes   []byte
	signKey     *rsa.PrivateKey
	verifyBytes []byte
	verifyKey   *rsa.PublicKey
	logger      log.Factory
	metrics     *metrics.Provider
	tracer      opentracing.Tracer
}

// New returns a JWT provider used for Signing and Verifying token
func New(cfg *config.Config, logger log.Factory, tracer opentracing.Tracer, metricProvider *metrics.Provider) (Provider, error) {

	signBytes, err := ioutil.ReadFile(cfg.Token.AuthPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := ioutil.ReadFile(cfg.Token.AuthPublicKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return &provider{
		cfg:         cfg,
		signBytes:   signBytes,
		signKey:     signKey,
		verifyBytes: verifyBytes,
		verifyKey:   verifyKey,
		tracer:      tracer,
		metrics:     metricProvider,
		logger:      logger.With(zap.String("package", "jwt")),
	}, nil
}

func (p *provider) investigateJWTError(ctx context.Context, err error) {
	// Lets audit these JWT Errors
	switch err.Error() {
	case jwt.ErrInvalidKeyType.Error(),
		jwt.ErrECDSAVerification.Error(),
		jwt.ErrHashUnavailable.Error(),
		jwt.ErrInvalidKey.Error(),
		jwt.ErrKeyMustBePEMEncoded.Error(),
		jwt.ErrNotECPrivateKey.Error(),
		jwt.ErrNotECPublicKey.Error(),
		jwt.ErrNotRSAPrivateKey.Error(),
		jwt.ErrNotRSAPublicKey.Error(),
		jwt.ErrSignatureInvalid.Error(),
		jwt.NoneSignatureTypeDisallowedError.Error():
		p.logger.For(ctx).Error(auditEventJWTError, zap.Error(err), zap.Bool("audit", true))
		p.metrics.StatAuditCount.WithLabelValues(auditEventJWTError).Inc()
	default:
		// Trap JWT Validation Errors
		// A "normal" error is an Expired Token (ValidationErrorExpired)
		// unknown/suspicious errors are basically everything else
		// we want to audit these suspicious/unknown errors
		// more than likely a client is messing with the token (i.e. hacker)
		// TODO: implement audit logger (isolate audit logs from application logs)
		if pgerr, ok := err.(*jwt.ValidationError); ok {
			switch pgerr.Errors {
			case jwt.ValidationErrorExpired:
				// We don't care about Expired Tokens
				// break out and move on
				break
			case
				// Someone is messing with the token
				// lets audit these validation errors
				jwt.ValidationErrorAudience,
				jwt.ValidationErrorClaimsInvalid,
				jwt.ValidationErrorId,
				jwt.ValidationErrorIssuedAt,
				jwt.ValidationErrorIssuer,
				jwt.ValidationErrorMalformed,
				jwt.ValidationErrorNotValidYet,
				jwt.ValidationErrorSignatureInvalid,
				jwt.ValidationErrorUnverifiable:
				p.logger.For(ctx).Error(auditEventJWTValidation, zap.Error(err), zap.Bool("audit", true))
				p.metrics.StatAuditCount.WithLabelValues(auditEventJWTValidation).Inc()
			default:
				// Not sure this should ever happen
				// but if it does, we should audit it
				p.logger.For(ctx).Error(fmt.Sprintf("unknown %s", auditEventJWTValidation), zap.Error(err), zap.Bool("audit", true))
				p.metrics.StatAuditCount.WithLabelValues(auditEventJWTValidation).Inc()
			}
		}
	}
}

func (p *provider) IsValidAccessToken(ctx context.Context, tkn string) (*accessTokenClaims, bool) {
	p.logger.For(ctx).Info("entering jwtservice.IsValidAccessToken")
	// Parse the token
	token, err := jwt.ParseWithClaims(tkn, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		p.logger.For(ctx).Info("return token key")
		return p.verifyKey, nil
	})

	// Let's investigate the JWT error
	// these errors could be from a potential hacker or someone misusing the token
	// let's audit these errors (potential IP Block)
	if err != nil {
		p.logger.For(ctx).Error("invalid access token", zap.Error(err))
		p.investigateJWTError(ctx, err)
		return nil, false
	}

	// Token Claim
	tokenClaims := token.Claims.(*accessTokenClaims)

	p.logger.For(ctx).Info("leaving jwtservice.IsValidAccessToken", zap.Bool("is_valid", token.Valid))
	return tokenClaims, token.Valid
}



func (p *provider) GenerateAccessToken(ctx context.Context, aTokenChan chan tokenmodels.TokenResult, tokenData map[string]interface{}) {
	p.logger.For(ctx).Info("entering jwtservice.GenerateAccessToken")
	accessToken := jwt.New(jwt.GetSigningMethod("RS256"))
	accessToken.Claims = &accessTokenClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(p.cfg.Token.AccessTokenLifeSpanMins)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    p.cfg.Token.Issuer,
			Subject:   tokenData["subject"].(string),
			Id:        tokenData["id"].(string),
		},
		customClaims{
			Roles: []string{"test"},
			Name:  tokenData["name"].(string),
		},
	}
	accessTokenSigned, err := accessToken.SignedString(p.signKey)
	aTokenChan <- tokenmodels.TokenResult{Token: accessTokenSigned, Err: err}
	p.logger.For(ctx).Info("leaving jwtservice.GenerateAccessToken")
	close(aTokenChan)
}

func (p *provider) GenerateRefreshToken(ctx context.Context, rTokenChan chan tokenmodels.TokenResult, tokenData map[string]interface{}) {
	p.logger.For(ctx).Info("entering jwtservice.GenerateRefreshToken")
	refreshToken := jwt.New(jwt.GetSigningMethod("RS256"))
	refreshToken.Claims = &refreshTokenClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(p.cfg.Token.RefreshTokenLifeSpanMins)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    p.cfg.Token.Issuer,
			Subject:   tokenData["subject"].(string),
			Id:        tokenData["id"].(string),
		},
	}
	refreshTokenSigned, err := refreshToken.SignedString(p.signKey)
	rTokenChan <- tokenmodels.TokenResult{Token: refreshTokenSigned, Err: err}
	p.logger.For(ctx).Info("leaving jwtservice.GenerateRefreshToken")
	close(rTokenChan)
}
