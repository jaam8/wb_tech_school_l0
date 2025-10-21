package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key string

const (
	KeyForLogger    key    = "logger"
	KeyForRequestID key    = "request_id"
	KeyForLogLevel  key    = "log_level"
	DebugLvl        string = "debug"
	InfoLvl         string = "info"
)

type Logger struct {
	l *zap.Logger
}

func NewLogger(logLevel string) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	})

	var level zapcore.Level
	switch logLevel {
	case DebugLvl:
		level = zap.DebugLevel
	case InfoLvl:
		level = zap.InfoLevel
	}
	config.Level.SetLevel(level)

	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	loggerStruct := &Logger{l: logger}

	return loggerStruct, nil
}

func New(ctx context.Context) (context.Context, error) {
	logLevel, ok := ctx.Value(KeyForLogLevel).(string)
	if !ok {
		logLevel = "debug"
	}

	loggerStruct, err := NewLogger(logLevel)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, KeyForLogger, loggerStruct)

	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(KeyForLogger).(*Logger)
}

func TryAppendRequestIDFromContext(ctx context.Context, fields []zap.Field) []zap.Field {
	if ctx.Value(KeyForRequestID) != nil {
		fields = append(fields, zap.String(string(KeyForRequestID), ctx.Value(KeyForRequestID).(string)))
	}

	return fields
}

func GetOrCreateLoggerFromCtx(ctx context.Context) *Logger {
	logger := GetLoggerFromCtx(ctx)
	if logger == nil {
		logLevel, ok := ctx.Value(KeyForLogLevel).(string)
		if !ok {
			logLevel = "debug"
		}
		logger, _ = NewLogger(logLevel)
	}

	return logger
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	GetLoggerFromCtx(ctx).l.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	GetLoggerFromCtx(ctx).l.Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	GetLoggerFromCtx(ctx).l.Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	GetLoggerFromCtx(ctx).l.Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = TryAppendRequestIDFromContext(ctx, fields)
	GetLoggerFromCtx(ctx).l.Fatal(msg, fields...)
}

func With(ctx context.Context, fields ...zap.Field) context.Context {
	currentLogger := GetLoggerFromCtx(ctx)
	newZapLogger := currentLogger.l.With(fields...)
	currentLogger.l = newZapLogger

	return ctx
}
