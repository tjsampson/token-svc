package authservice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/datastores/redis"
	"github.com/tjsampson/token-svc/internal/errors"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/authmodels"
	"github.com/tjsampson/token-svc/internal/models/tokenmodels"
	"github.com/tjsampson/token-svc/internal/models/usermodels"
	"github.com/tjsampson/token-svc/internal/repos/userrepo"
	"github.com/tjsampson/token-svc/internal/services/cookieservice"
	"github.com/tjsampson/token-svc/internal/services/jwtservice"
	"github.com/tjsampson/token-svc/internal/services/tracingservice"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	lockKeyVal = "1"
)

// Service is the auth service interface
type Service interface {
	Login(ctx context.Context, creds *authmodels.UserCreds) (authmodels.LoginResponse, error)
	Register(ctx context.Context, creds *authmodels.UserRegistration) (usermodels.Record, error)
}

type service struct {
	logger        log.Factory
	cfg           *config.Config
	jwtClient     jwtservice.Provider
	userRepo      userrepo.Store
	tracer        opentracing.Tracer
	traceProvider tracingservice.Provider
	redis         redis.Provider
	cookieOven    cookieservice.Provider
}

// New returns a new Service interface implementation
func New(logger log.Factory, cfg *config.Config, jwtClient jwtservice.Provider, usrRepo userrepo.Store, tracer opentracing.Tracer, traceProvider tracingservice.Provider, redis redis.Provider, cookieOven cookieservice.Provider) Service {
	return &service{
		logger:        logger.With(zap.String("package", "authservice")),
		cfg:           cfg,
		jwtClient:     jwtClient,
		userRepo:      usrRepo,
		tracer:        tracer,
		redis:         redis,
		cookieOven:    cookieOven,
		traceProvider: traceProvider,
	}
}

func (svc *service) Register(ctx context.Context, userReg *authmodels.UserRegistration) (usermodels.Record, error) {
	svc.logger.For(ctx).Info("entering authservice.Register", zap.String("email", userReg.Email))
	var err error
	var result usermodels.Record
	svc.logger.For(ctx).Info("start authservice.Register.hashPassword", zap.String("email", userReg.Email))
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
	svc.logger.For(ctx).Info("stop authservice.Register.hashPassword", zap.String("email", userReg.Email))
	if err != nil {
		svc.logger.For(ctx).Error("failed to GenerateFromPassword", zap.Error(err), zap.String("email", userReg.Email))
		return result, errors.ErrorWrapper(err, "AuthService.Register.hashPassword")
	}

	user, err := svc.userRepo.Insert(ctx, userReg.Email, string(passHashBytes))

	if err != nil {
		svc.logger.For(ctx).Error("failed to insert user", zap.Error(err), zap.String("email", userReg.Email))
		return result, errors.ErrorWrapper(err, fmt.Sprintf("failed to register %s", userReg.Email))
	}
	svc.logger.For(ctx).Info("leaving authservice.Register", zap.String("email", userReg.Email))
	return user, nil
}

func (svc *service) validateUserCreds(ctx context.Context, creds *authmodels.UserCreds) (usermodels.Record, error) {
	svc.logger.For(ctx).Info("entering authservice.validateUserCreds", zap.String("email", creds.Email))
	user := usermodels.Record{}
	var invalidCredsErr = func(originalErr error) error {
		return &errors.RestError{
			Code:          401,
			Message:       "Invalid Credentials",
			OriginalError: originalErr,
		}
	}
	svc.logger.For(ctx).Info("start authservice.validateUserCreds.ReadByEmail", zap.String("email", creds.Email))
	user, err := svc.userRepo.ReadByEmail(ctx, creds.Email)
	svc.logger.For(ctx).Info("stop authservice.validateUserCreds.ReadByEmail", zap.String("email", creds.Email))
	if err != nil {
		svc.logger.For(ctx).Error("failed authservice.validateUserCreds.ReadByEmail", zap.Error(err))
		return user, invalidCredsErr(err)
	}
	svc.logger.For(ctx).Info("start authservice.validateUserCreds.CompareHashAndPassword", zap.String("email", creds.Email))
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	svc.logger.For(ctx).Info("stop authservice.validateUserCreds.CompareHashAndPassword", zap.String("email", creds.Email))
	if err != nil {
		svc.logger.For(ctx).Error("failed authservice.validateUserCreds.CompareHashAndPassword", zap.Error(err))
		return user, invalidCredsErr(err)
	}
	svc.logger.For(ctx).Info("leaving authservice.validateUserCreds", zap.String("email", creds.Email))
	return user, nil
}

func (svc *service) failedLoginAttemps(ctx context.Context, email string) int {
	failedAttemps := 0
	failedLogins, err := svc.redis.Get(ctx, fmt.Sprintf("%v-%v", svc.cfg.Token.FailedLoginCacheKeyID, email))
	if err != nil {
		// If this is first failure, the key does not exists
		svc.logger.For(ctx).Error("failed login attempt cache retrieval", zap.String("email", email))
	} else {
		svc.logger.For(ctx).Info("failed login attempts", zap.String("email", email), zap.String("failed_attempts", failedLogins))
		i, err := strconv.Atoi(failedLogins)
		if err != nil {
			svc.logger.For(ctx).Error("failed login attempt int conversion", zap.String("email", email), zap.String("failedLogins", failedLogins))
		}
		failedAttemps = i
	}
	failedAttemps = failedAttemps + 1

	// Set the Redis Cache Key for failed login attempts
	err = svc.redis.Set(
		ctx,
		fmt.Sprintf("%v-%v", svc.cfg.Token.FailedLoginCacheKeyID, email),
		strconv.Itoa(failedAttemps),
		time.Duration(svc.cfg.Token.FailedLoginAttemptCacheLifeSpanMins)*time.Minute,
	)
	if err != nil {
		svc.logger.For(ctx).Error("cache set failed: failed login attempt", zap.Error(err), zap.String("email", email))
	}
	return failedAttemps
}

func (svc *service) setUserLock(ctx context.Context, email string) {
	// Set the Redis Cache Key for user account locked
	err := svc.redis.Set(
		ctx,
		fmt.Sprintf("%v-%v", svc.cfg.Cache.UserAccountLockedKeyID, email),
		lockKeyVal,
		time.Duration(svc.cfg.Cache.UserAccountLockedLifeSpanMins)*time.Minute,
	)
	if err != nil {
		svc.logger.For(ctx).Error("cache set failed: setUserLock", zap.Error(err), zap.String("email", email))
	}
}

// Login accepts UserCreds and generates the following...
//  - JWT Access Token
// 	- JWT Refresh Token
//	- Secure Cookie
func (svc *service) Login(ctx context.Context, creds *authmodels.UserCreds) (authmodels.LoginResponse, error) {
	svc.logger.For(ctx).Info("entering authservice.Login", zap.String("email", creds.Email))

	locked, _ := svc.redis.Get(ctx, fmt.Sprintf("%v-%v", svc.cfg.Cache.UserAccountLockedKeyID, creds.Email))

	if locked == lockKeyVal {
		return authmodels.LoginResponse{}, &errors.RestError{
			Code:    403, // Forbidden - Account is currently Locked
			Message: fmt.Sprintf("user account locked (%s)", creds.Email),
		}
	}

	user, err := svc.validateUserCreds(ctx, creds)
	if err != nil {
		svc.logger.For(ctx).Error("failed user cred validation", zap.Error(err), zap.String("email", creds.Email))
		if svc.failedLoginAttemps(ctx, creds.Email) > int(svc.cfg.Token.FailedLoginAttemptsMax) {
			svc.logger.For(ctx).Info("user exceeded failed login attempts", zap.String("email", creds.Email), zap.Uint16("failed_attempts_max", svc.cfg.Token.FailedLoginAttemptsMax))
			svc.setUserLock(ctx, creds.Email)
			return authmodels.LoginResponse{}, &errors.RestError{
				Code:          418, // I'm a teapot HTTP 218 (people could be trying to break in)
				Message:       fmt.Sprintf("user account locked (%s)", creds.Email),
				OriginalError: err,
			}
		}
		return authmodels.LoginResponse{}, errors.ErrorWrapper(err, "HealthService.GetDatabaseHealth")
	}

	// Setup our Channels for concurrent calls
	accessTokenChan := make(chan tokenmodels.TokenResult, 1)
	refreshTokenChan := make(chan tokenmodels.TokenResult, 1)
	cookieDataChan := make(chan *http.Cookie, 1)

	// Generate the New JWT IDs (GUIDs)
	accessTokenID := uuid.NewV4().String()
	refreshTokenID := uuid.NewV4().String()

	// Establish the cookie data
	cookieData := map[string]string{
		svc.cfg.Cookie.KeyUserID:       string(user.ID),
		svc.cfg.Cookie.KeyEmail:        user.Email,
		svc.cfg.Cookie.KeyJWTAccessID:  accessTokenID,
		svc.cfg.Cookie.KeyJWTRefreshID: refreshTokenID}

	// Establish accesssTokenData
	accessTokenData := map[string]interface{}{
		"subject": strconv.Itoa(user.ID),
		"id": accessTokenID,
		"name": user.Email,
	}

	// Establish refreshTokenData
	refreshTokenData := map[string]interface{}{
		"subject": strconv.Itoa(user.ID),
		"id": refreshTokenID,
		"name": user.Email,
	}

	go svc.cookieOven.BakeCookie(ctx, cookieDataChan, cookieData)
	go svc.jwtClient.GenerateAccessToken(ctx, accessTokenChan, accessTokenData)
	go svc.jwtClient.GenerateRefreshToken(ctx, refreshTokenChan, refreshTokenData)

	result := authmodels.LoginResponse{}
	for {
		accessTokenResult, accessTokenOK := <-accessTokenChan
		refreshTokenResult, refreshTokenOK := <-refreshTokenChan
		cookieResult, cookieOK := <-cookieDataChan

		if accessTokenOK == false && refreshTokenOK == false && cookieOK == false {
			break
		}
		if accessTokenResult.Err != nil {
			svc.logger.For(ctx).Error("access token failed", zap.Error(accessTokenResult.Err))
			err = accessTokenResult.Err
			break
		}
		if refreshTokenResult.Err != nil {
			svc.logger.For(ctx).Error("refresh token failed", zap.Error(refreshTokenResult.Err))
			err = refreshTokenResult.Err
			break
		}
		if cookieResult.Name == "" {
			// if the cookie service fails to encode the cookie
			// we log the error, return an empty cookie
			err = fmt.Errorf("failed to bake cookie")
			break
		}

		result.AccessToken = accessTokenResult.Token
		result.RefreshToken = refreshTokenResult.Token
		result.HTTPCookie = cookieResult
	}

	if err != nil {
		svc.logger.For(ctx).Error("failed token/cookie generation", zap.Error(err))
		return result, errors.ErrorWrapper(err, "AuthService.Login")
	}

	// Set the Redis Cache Key
	if err = svc.redis.Set(ctx, fmt.Sprintf("%v-%v", svc.cfg.Token.AccessCacheKeyID, user.ID), accessTokenID, time.Duration(svc.cfg.Token.AccessTokenLifeSpanMins)*time.Minute); err != nil {
		svc.logger.For(ctx).Error("failed set token cache", zap.Error(err))
		return result, errors.ErrorWrapper(err, "AuthService.Login")
	}
	svc.logger.For(ctx).Info("leaving authservice.Login", zap.String("email", creds.Email))
	return result, nil
}
