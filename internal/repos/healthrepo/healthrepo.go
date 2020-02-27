package healthrepo

import (
	"context"
	"database/sql"

	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/healthmodels"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// Store interface
type Store interface {
	HealthStatus(ctx context.Context) (healthmodels.DBHealth, error)
}

// New returns a conrete implementation of the Store interface
func New(dbConn *sql.DB, logger log.Factory, tracer opentracing.Tracer) Store {
	return &store{
		db:     dbConn,
		logger: logger.With(zap.String("package", "healthrepo")),
		tracer: tracer,
	}
}

type store struct {
	logger log.Factory
	db     *sql.DB
	tracer opentracing.Tracer
}

func (s *store) HealthStatus(ctx context.Context) (healthmodels.DBHealth, error) {
	s.logger.For(ctx).Info("health status")
	var err error
	var sqlVersion []uint8
	result := healthmodels.DBHealth{}

	query := "SELECT version();"

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "postgres")
		span.SetTag("sql.query", query)
		defer span.Finish()
		// ctx = opentracing.ContextWithSpan(ctx, span)
	}

	err = s.db.QueryRow(query).Scan(&sqlVersion)
	if err != nil {
		return result, err
	}
	result.Status = "OK"
	result.Version = string(sqlVersion)
	return result, err
}
