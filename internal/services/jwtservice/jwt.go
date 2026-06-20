package jwtservice

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/tokenmodels"
	"github.com/tjsampson/token-svc/pkg/metrics"

	"github.com/golang-jwt/jwt/v5"
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
		jwt.RegisteredClaims
		customClaims
	}

	refreshTokenClaims struct {
		jwt.RegisteredClaims
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

	signBytes, err := os.ReadFile(cfg.Token.AuthPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := os.ReadFile(cfg.Token.AuthPublicKeyPath)
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
	// Key/configuration errors are server-side issues — audit immediately.
	if errors.Is(err, jwt.ErrInvalidKeyType) ||
		errors.Is(err, jwt.ErrECDSAVerification) ||
		errors.Is(err, jwt.ErrHashUnavailable) ||
		errors.Is(err, jwt.ErrInvalidKey) ||
		errors.Is(err, jwt.ErrKeyMustBePEMEncoded) ||
		errors.Is(err, jwt.ErrNotECPrivateKey) ||
		errors.Is(err, jwt.ErrNotECPublicKey) ||
		errors.Is(err, jwt.ErrNotRSAPrivateKey) ||
		errors.Is(err, jwt.ErrNotRSAPublicKey) ||
		errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		p.logger.For(ctx).Error(auditEventJWTError, zap.Error(err), zap.Bool("audit", true))
		p.metrics.StatAuditCount.WithLabelValues(auditEventJWTError).Inc()
		return
	}

	// Expired tokens are normal — clients will refresh; no audit needed.
	if errors.Is(err, jwt.ErrTokenExpired) {
		return
	}

	// Everything else (malformed, wrong issuer/audience, tampered claims, etc.)
	// is suspicious and should be audited as a potential attack.
	if errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenUnverifiable) ||
		errors.Is(err, jwt.ErrTokenInvalidAudience) ||
		errors.Is(err, jwt.ErrTokenInvalidClaims) ||
		errors.Is(err, jwt.ErrTokenInvalidId) ||
		errors.Is(err, jwt.ErrTokenInvalidIssuer) ||
		errors.Is(err, jwt.ErrTokenInvalidSubject) ||
		errors.Is(err, jwt.ErrTokenNotValidYet) {
		p.logger.For(ctx).Error(auditEventJWTValidation, zap.Error(err), zap.Bool("audit", true))
		p.metrics.StatAuditCount.WithLabelValues(auditEventJWTValidation).Inc()
		return
	}

	// Catch-all for any unclassified error — audit it.
	p.logger.For(ctx).Error(fmt.Sprintf("unknown %s", auditEventJWTValidation), zap.Error(err), zap.Bool("audit", true))
	p.metrics.StatAuditCount.WithLabelValues(auditEventJWTValidation).Inc()
}

func (p *provider) IsValidAccessToken(ctx context.Context, tkn string) (*accessTokenClaims, bool) {
	p.logger.For(ctx).Info("entering jwtservice.IsValidAccessToken")
	token, err := jwt.ParseWithClaims(tkn, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		p.logger.For(ctx).Info("return token key")
		return p.verifyKey, nil
	})

	if err != nil {
		p.logger.For(ctx).Error("invalid access token", zap.Error(err))
		p.investigateJWTError(ctx, err)
		return nil, false
	}

	tokenClaims := token.Claims.(*accessTokenClaims)

	p.logger.For(ctx).Info("leaving jwtservice.IsValidAccessToken", zap.Bool("is_valid", token.Valid))
	return tokenClaims, token.Valid
}

func (p *provider) GenerateAccessToken(ctx context.Context, aTokenChan chan tokenmodels.TokenResult, tokenData map[string]interface{}) {
	p.logger.For(ctx).Info("entering jwtservice.GenerateAccessToken")
	accessToken := jwt.New(jwt.SigningMethodRS256)
	accessToken.Claims = &accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(p.cfg.Token.AccessTokenLifeSpanMins))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    p.cfg.Token.Issuer,
			Subject:   tokenData["subject"].(string),
			ID:        tokenData["id"].(string),
		},
		customClaims: customClaims{
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
	refreshToken := jwt.New(jwt.SigningMethodRS256)
	refreshToken.Claims = &refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(p.cfg.Token.RefreshTokenLifeSpanMins))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    p.cfg.Token.Issuer,
			Subject:   tokenData["subject"].(string),
			ID:        tokenData["id"].(string),
		},
	}
	refreshTokenSigned, err := refreshToken.SignedString(p.signKey)
	rTokenChan <- tokenmodels.TokenResult{Token: refreshTokenSigned, Err: err}
	p.logger.For(ctx).Info("leaving jwtservice.GenerateRefreshToken")
	close(rTokenChan)
}
