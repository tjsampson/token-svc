package serviceprovider

import (
	"database/sql"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/datastores/postgres"
	"github.com/tjsampson/token-svc/internal/datastores/redis"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/validation"

	"github.com/tjsampson/token-svc/internal/repos/healthrepo"
	"github.com/tjsampson/token-svc/internal/repos/userrepo"
	"github.com/tjsampson/token-svc/internal/services/authservice"
	"github.com/tjsampson/token-svc/internal/services/cookieservice"
	"github.com/tjsampson/token-svc/internal/services/healthservice"
	"github.com/tjsampson/token-svc/internal/services/jwtservice"
	"github.com/tjsampson/token-svc/internal/services/tracingservice"
	"github.com/tjsampson/token-svc/internal/services/userservice"
	"github.com/tjsampson/token-svc/pkg/metrics"
	"github.com/tjsampson/token-svc/pkg/version"

	"gopkg.in/go-playground/validator.v9"

	"go.uber.org/zap"
)

// Context is the service provider context
// which provides access to the service objects
// this is main container that holds our dependencies
type Context struct {
	DB            *sql.DB
	VersionInfo   version.Info
	Config        *config.Config
	Metrics       *metrics.Provider
	CookieOven    cookieservice.Provider
	RedisClient   redis.Provider
	TraceProvider tracingservice.Provider
	Validator     validation.Provider
	Logger        log.Factory
	AuthService   authservice.Service
	HealthService healthservice.Service
	JwtClient     jwtservice.Provider
	UserRepo      userrepo.Store
	UserService	  userservice.Service
}

var zlogger *zap.Logger

// Initialize is the initializer for the Provider
// if this fails, we should fail hard (panic/fatal etc.)
// this should only be called once, on app startup
// NEW is glue and this is where that glue is applied
// create your deps here and add them to the Context
// then your services can use those deps via (interface injection)
func Initialize(vInfo version.Info) *Context {
	cfg := config.Read()

	zlogger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}

	logger := log.NewFactory(cfg, zlogger.With(zap.String("package", "appcontext")))

	tracingProvider := tracingservice.New(cfg.API.ServiceName, logger, true)

	dbConn, err := postgres.Database(cfg)

	if err != nil {
		logger.Bg().Fatal("failed db connection", zap.Error(err))
	}

	metricProvider := metrics.New()

	jwtProvider, err := jwtservice.New(cfg, logger, tracingProvider.Tracer, metricProvider)

	if err != nil {
		logger.Bg().Fatal("failed jwt clien", zap.Error(err))
	}

	redisProvider, err := redis.New(cfg, logger, tracingservice.New("redis", logger, false).Tracer)

	if err != nil {

		logger.Bg().Fatal("failed redis client", zap.Error(err))
	}

	cookieOven := cookieservice.New(cfg, logger)

	healthRepo := healthrepo.New(dbConn, logger, tracingservice.New("postgres", logger, false).Tracer)

	healthSvc := healthservice.New(healthRepo, redisProvider, vInfo, metricProvider, logger)

	userRepo := userrepo.New(dbConn, logger, tracingservice.New("postgres", logger, false).Tracer)

	authSvc := authservice.New(logger, cfg, jwtProvider, userRepo, tracingProvider.Tracer, tracingProvider, redisProvider, cookieOven)

	validator := validation.New(validator.New())

	userSvc := userservice.New(logger,cfg,jwtProvider,userRepo,tracingProvider.Tracer, tracingProvider, redisProvider)

	return &Context{
		DB:            dbConn,
		Logger:        logger,
		CookieOven:    cookieOven,
		HealthService: healthSvc,
		AuthService:   authSvc,
		VersionInfo:   vInfo,
		Config:        cfg,
		RedisClient:   redisProvider,
		Metrics:       metricProvider,
		TraceProvider: tracingProvider,
		Validator:     validator,
		JwtClient:     jwtProvider,
		UserRepo:      userRepo,
		UserService: userSvc,
	}
}
