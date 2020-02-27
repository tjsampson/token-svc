package userrepo

import (
	"context"
	"database/sql"

	"github.com/tjsampson/token-svc/internal/datastores/postgres"
	"github.com/tjsampson/token-svc/internal/log"
	"github.com/tjsampson/token-svc/internal/models/usermodels"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// Store is the database store Interface
// typically maps to a table in the DB, but not a requirement
type Store interface {
	ReadByEmail(ctx context.Context, email string) (usermodels.Record, error)
	Insert(ctx context.Context, email, passHash string) (usermodels.Record, error)
	List(ctx context.Context) ([]usermodels.Record, error)
}

// New returns a conrete implementation of the Store interface
func New(dbConn *sql.DB, logger log.Factory, tracer opentracing.Tracer) Store {
	return &store{
		db:     dbConn,
		logger: logger.With(zap.String("package", "userrepo")),
		tracer: tracer,
	}
}

type store struct {
	db     *sql.DB
	tracer opentracing.Tracer
	logger log.Factory
}

func (s *store) Insert(ctx context.Context, email, passHash string) (usermodels.Record, error) {
	s.logger.For(ctx).Info("entering userrepo.Insert", zap.String("email", email))

	query := `
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2)
	RETURNING id`

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("SQL INSERT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "postgres")
		span.SetTag("sql.query", query)
		span.SetTag("param.email", email)
		defer span.Finish()
		// ctx = opentracing.ContextWithSpan(ctx, span)
	}

	_, err := s.db.Exec(query, email, passHash)
	if err != nil {
		s.logger.For(ctx).Error("failed userrepo.Insert.Exec", zap.Error(err), zap.String("email", email))
		return usermodels.Record{}, err
	}

	user, err := s.ReadByEmail(ctx, email)
	if err != nil {
		s.logger.For(ctx).Error("failed userrepo.Insert.ReadByEmail", zap.Error(err), zap.String("email", email))
		return usermodels.Record{}, postgres.ErrorCheck(err)
	}
	s.logger.For(ctx).Info("leaving userrepo.Insert", zap.String("email", email))
	return user, nil
}

func (s *store) ReadByEmail(ctx context.Context, email string) (usermodels.Record, error) {
	s.logger.For(ctx).Info("entering userrepo.ReadByEmail", zap.String("email", email))
	defer s.logger.For(ctx).Info("leaving userrepo.ReadByEmail", zap.String("email", email))
	userData := usermodels.Record{}
	query := "SELECT id, uid, email, email_verified, password_hash, created_at, updated_at FROM users WHERE email=$1"
	queryStmt, err := s.db.Prepare(query)

	if err != nil {
		s.logger.For(ctx).Error("failed to prepare read user by email", zap.Error(err))
		return userData, err
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "postgres")
		span.SetTag("sql.query", query)
		span.SetTag("param.email", email)
		defer span.Finish()
		// ctx = opentracing.ContextWithSpan(ctx, span)
	}

	err = queryStmt.QueryRow(email).Scan(&userData.ID, &userData.UID, &userData.Email, &userData.EmailVerified, &userData.PasswordHash, &userData.CreatedAt, &userData.UpdatedAt)
	if err != nil {
		s.logger.For(ctx).Error("failed userrepo.ReadByEmail.QueryRow", zap.Error(err), zap.String("email", email))
		return usermodels.Record{}, postgres.ErrorCheck(err)
	}

	return userData, nil

}

func (s *store) List(ctx context.Context) ([]usermodels.Record, error) {
	s.logger.For(ctx).Info("entering userrepo.List")
	defer s.logger.For(ctx).Info("leaving userrepo.List")
	users := []usermodels.Record{}
	query := "SELECT id, uid, email, email_verified, password_hash, created_at, updated_at FROM users"

	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "postgres")
		span.SetTag("sql.query", query)
		defer span.Finish()
		// ctx = opentracing.ContextWithSpan(ctx, span)
	}

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.For(ctx).Error("failed userrepo.List.Query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user usermodels.Record

		err = rows.Scan(&user.ID, &user.UID, &user.Email, &user.EmailVerified, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			s.logger.For(ctx).Error("failed userrepo.List.Query.Rows.Scan")
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
