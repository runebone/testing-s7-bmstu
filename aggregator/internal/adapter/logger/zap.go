package logger

import (
	"aggregator/internal/common/logger"
	"context"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	zapLogger, _ := zap.NewProduction()
	return &ZapLogger{logger: zapLogger}
}

func (l *ZapLogger) zapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return &ZapLogger{logger: l.logger.With(l.zapFields(fields)...)}
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Debug(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Info(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Warn(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Error(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Fatal(msg, zap.Any("context", fields))
}
