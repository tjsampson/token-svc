package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/internal/log"

	"github.com/tjsampson/token-svc/internal/errors"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	tags "github.com/opentracing/opentracing-go/ext"
)

// Provider is the redis client provider interface
type Provider interface {
	Ping(ctx context.Context) (string, error)
	Close() error
	Set(ctx context.Context, key string, value string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type provider struct {
	client *redis.Client
	cfg    *config.Config
	logger log.Factory
	tracer opentracing.Tracer
}

func addr(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port)
}

// New returns a new Redis Provider
func New(cfg *config.Config, logger log.Factory, tracer opentracing.Tracer) (Provider, error) {
	logs := logger.With(zap.String("package", "redis"))

	client := redis.NewClient(&redis.Options{
		Addr:     addr(cfg),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if client == nil {
		return nil, &errors.RestError{
			Code:          503,
			Message:       "invalid redis client",
			OriginalError: nil,
		}
	}

	return &provider{
		client: client,
		cfg:    cfg,
		logger: logs,
		tracer: tracer,
	}, nil
}

func (p *provider) Get(ctx context.Context, key string) (string, error) {
	p.logger.For(ctx).Info("return cache", zap.String("cache_key", key))
	val, err := p.client.Get(key).Result()
	if err == redis.Nil {
		p.logger.For(ctx).Error("missing cache key", zap.String("cache_key", key), zap.Error(err))
		return "", err
	} else if err != nil {
		p.logger.For(ctx).Error("failed get cache key", zap.String("cache_key", key), zap.Error(err))
		return "", err
	}
	p.logger.For(ctx).Info("cache returned", zap.String("cache_value", val))
	return val, nil

}

func (p *provider) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := p.tracer.StartSpan("CACHE SET", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "redis")
		span.SetTag("param.key", key)
		span.SetTag("param.value", value)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	p.logger.For(ctx).Info("set cache", zap.String("cache_key", key), zap.String("cache_value", value))

	err := p.client.Set(key, value, exp).Err()

	if err != nil {
		p.logger.For(ctx).Error("failed to set cache", zap.Error(err))
	}

	p.logger.For(ctx).Info("cache set", zap.String("cache_key", key), zap.String("cache_value", value))
	return err
}

// Ping issues a ping to the redis server to check health/status
func (p *provider) Ping(ctx context.Context) (string, error) {
	p.logger.For(ctx).Info("cache ping")
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := p.tracer.StartSpan("CACHE PING", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "redis")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	pong, err := p.client.Ping().Result()
	if err != nil {
		p.logger.For(ctx).Error("cache ping error", zap.Error(err))
		return "", &errors.RestError{
			Code:          503,
			Message:       "invalid redis ping",
			OriginalError: err,
		}
	}
	p.logger.For(ctx).Info("cache pong", zap.String("ping", pong))
	return pong, nil
}

func (p *provider) Close() error {
	if p.client == nil {
		return nil
	}

	err := p.Close()
	if err != nil {
		p.logger.Bg().Error("failed to close redis connection", zap.Error(err))
	}
	return nil
}
