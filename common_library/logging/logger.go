package logging

import (
	"common_library/ctxdata"
	"context"
	"go.uber.org/zap"
)

type loggerKey struct{}

const (
	traceIdKey = "trace_id"
)

var (
	loggerKeyInstance = loggerKey{}
)

type Logger struct {
	l *zap.Logger
}

func New(l *zap.Logger) *Logger {
	return &Logger{l}
}

func ContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKeyInstance, logger)
}

func GetFromContext(ctx context.Context) (*Logger, bool) {
	logger, ok := ctx.Value(loggerKeyInstance).(*Logger)
	return logger, ok
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Error(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	args = argsWithTraceID(ctx, args)
	l.l.Sugar().Infof(format, args...)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	args = argsWithTraceID(ctx, args)
	l.l.Sugar().Warnf(format, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	args = argsWithTraceID(ctx, args)
	l.l.Sugar().Errorf(format, args...)
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	args = argsWithTraceID(ctx, args)
	l.l.Sugar().Debugf(format, args...)
}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	args = argsWithTraceID(ctx, args)
	l.l.Sugar().Fatalf(format, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Info(msg, append(fields, zap.Any("context", ctx))...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Warn(msg, append(fields, zap.Any("context", ctx))...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Error(msg, append(fields, zap.Any("context", ctx))...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	fields = fieldsWithTraceID(ctx, fields)
	l.l.Debug(msg, append(fields, zap.Any("context", ctx))...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l: l.l.With(fields...)}
}

func fieldsWithTraceID(ctx context.Context, fields []zap.Field) []zap.Field {
	if traceId, ok := ctxdata.GetTraceID(ctx); ok {
		fields = append(fields, zap.String(traceIdKey, traceId))
	}
	return fields
}

func argsWithTraceID(ctx context.Context, args []any) []any {
	if traceId, ok := ctxdata.GetTraceID(ctx); ok {
		args = append(args, traceIdKey, traceId)
	}
	return args
}
