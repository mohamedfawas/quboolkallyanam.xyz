package logger

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	constants "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitLogger(service string) {
	env := os.Getenv(constants.ENVIRONMENT)

	cfg := zap.NewProductionConfig()
	if env != constants.EnvProduction {
		cfg = zap.NewDevelopmentConfig()
	}

	// Use ISO8601 time format
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	// Disable stack traces in production
	if env == constants.EnvProduction {
		cfg.DisableStacktrace = true
		// Also disable caller info in production for cleaner logs
		cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}

	sugar := logger.With(
		zap.String(constants.Service, service),
		zap.String(constants.Env, env),
	).Sugar()

	Log = sugar
}

func Sync() {
	if Log != nil {
		_ = Log.Sync() // flush logs from buffer to stdout
	}
}

func LogWithContext(ctx context.Context, level, msg string, data map[string]interface{}) {
	reqID, _ := ctx.Value(constants.ContextKeyRequestID).(string)
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	fields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("trace_id", sc.TraceID().String()),
		zap.String("span_id", sc.SpanID().String()),
	}

	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case "debug":
		Log.Desugar().Debug(msg, fields...)
	case "warn":
		Log.Desugar().Warn(msg, fields...)
	case "error":
		Log.Desugar().Error(msg, fields...)
	case "fatal":
		Log.Desugar().Fatal(msg, fields...)
	default:
		Log.Desugar().Info(msg, fields...)
	}
}

// Debug logs a debug-level message
func Debug(ctx context.Context, msg string, data map[string]interface{}) {
	LogWithContext(ctx, "debug", msg, data)
}

// Info logs an info-level message
func Info(ctx context.Context, msg string, data map[string]interface{}) {
	LogWithContext(ctx, "info", msg, data)
}

// Warn logs a warn-level message
func Warn(ctx context.Context, msg string, data map[string]interface{}) {
	LogWithContext(ctx, "warn", msg, data)
}

// Error logs an error-level message
func Error(ctx context.Context, msg string, data map[string]interface{}) {
	LogWithContext(ctx, "error", msg, data)
}

// Fatal logs a fatal-level message and exits the application
func Fatal(ctx context.Context, msg string, data map[string]interface{}) {
	LogWithContext(ctx, "fatal", msg, data)
}

func ErrorWithStack(ctx context.Context, msg string, err error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["error"] = err.Error()
	data["error_type"] = fmt.Sprintf("%T", err)

	// Only include stack trace in development
	if os.Getenv("ENVIRONMENT") != "production" {
		data["stack_trace"] = string(debug.Stack())
	}

	Error(ctx, msg, data)
}
