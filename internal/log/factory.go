package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tjsampson/token-svc/internal/config"
	"github.com/tjsampson/token-svc/pkg/metrics"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

func logLevel(cfgLevel string) zap.AtomicLevel {
	var logLevel zap.AtomicLevel
	switch cfgLevel {
	case "debug":
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error", "err":
		logLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		logLevel = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		logLevel = zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return logLevel
}

func PrometheusHook(e zapcore.Entry) error {
	return nil
}

// NewFactory creates a new Factory.
func NewFactory(appCfg *config.Config, logger *zap.Logger) Factory {
	cfg := zap.Config{
		Level:            logLevel(appCfg.Logger.Level),
		Encoding:         appCfg.Logger.Encoding,
		OutputPaths:      appCfg.Logger.OutputPaths,
		ErrorOutputPaths: appCfg.Logger.ErrorOutputPaths,
		InitialFields: map[string]interface{}{
			"service": appCfg.API.ServiceName,
		},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zapcore.RegisterHooks(logger.Core(), metrics.LogLevelPrometheusHook().Fire)

	return Factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// TODO for Jaeger span extract trace/span IDs as fields
		return spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: b.logger.With(fields...)}
}
