package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type auditLogger struct {
	logger *zap.Logger
}

func (al auditLogger) Audit(msg string, fields ...zapcore.Field) {
	al.logger.Info(msg, fields...)
}

func (al auditLogger) Panic(msg string, fields ...zapcore.Field) {
	al.logger.Panic(msg, fields...)
}

func (al auditLogger) Debug(msg string, fields ...zapcore.Field) {
	al.logger.Debug(msg, fields...)
}

func (al auditLogger) Info(msg string, fields ...zapcore.Field) {
	al.logger.Info(msg, fields...)
}

func (al auditLogger) Error(msg string, fields ...zapcore.Field) {
	al.logger.Error(msg, fields...)
}

func (al auditLogger) Fatal(msg string, fields ...zapcore.Field) {
	al.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
// func (al auditLogger) With(fields ...zapcore.Field) Logger {
// 	return auditLogger{logger: al.logger.With(fields...), span: al.span}
// }
