package userservice

import (
	"context"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/datastores/redis"
	"github.com/tjsampson/token-svc/internal/errors"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/usermodels"
	"github.com/tjsampson/token-svc/internal/repos/userrepo"
	"github.com/tjsampson/token-svc/internal/services/jwtservice"
	"github.com/tjsampson/token-svc/internal/services/tracingservice"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// Service is the User Service Contract interface
type Service interface {
	List(ctx context.Context) ([]usermodels.Record, error)
}

type service struct {
	logger        log.Factory
	cfg           *config.Config
	jwtClient     jwtservice.Provider
	userRepo      userrepo.Store
	tracer        opentracing.Tracer
	traceProvider tracingservice.Provider
	redis         redis.Provider
}

// New returns a new Service interface implementation
func New(logger log.Factory, cfg *config.Config, jwtClient jwtservice.Provider, usrRepo userrepo.Store, tracer opentracing.Tracer, traceProvider tracingservice.Provider, redis redis.Provider) Service {
	return &service{
		logger:        logger.With(zap.String("package", "userservice")),
		cfg:           cfg,
		jwtClient:     jwtClient,
		userRepo:      usrRepo,
		tracer:        tracer,
		redis:         redis,
		traceProvider: traceProvider,
	}
}

// List returns a full list of customers
// TODO: Add Admin role/group check
// TODO: Add pagination
func (svc *service) List(ctx context.Context) ([]usermodels.Record, error) {
	svc.logger.For(ctx).Info("entering userservice List")
	users, err := svc.userRepo.List(ctx)
	svc.logger.For(ctx).Info("leaving userservice List")
	return users, errors.ErrorWrapper(err, "UserService.List")
}
