package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 日志器
type Logger struct {
	*zap.Logger
}

// New 创建日志器
func New(service string) *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 输出到 stdout
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	l, _ := config.Build(
		zap.Fields(
			zap.String("service", service),
		),
	)
	return &Logger{l}
}

// WithContext 从 Context 中提取 Trace 信息并附加到日志
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return &Logger{l.Logger.With(
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)}
	}
	return l
}

// WithTenant 添加租户 ID
func (l *Logger) WithTenant(tenantID int64) *Logger {
	return &Logger{l.Logger.With(zap.Int64("tenant_id", tenantID))}
}

// WithBrand 添加品牌 ID
func (l *Logger) WithBrand(brandID int64) *Logger {
	return &Logger{l.Logger.With(zap.Int64("brand_id", brandID))}
}

// WithUser 添加用户 ID
func (l *Logger) WithUser(userID int64) *Logger {
	return &Logger{l.Logger.With(zap.Int64("user_id", userID))}
}

// WithError 添加错误信息
func (l *Logger) WithError(err error) *Logger {
	return &Logger{l.Logger.With(zap.Error(err))}
}

// InfoContext 记录 Info 级别日志（带 Context）
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Logger.Info(msg, fields...)
}

// ErrorContext 记录 Error 级别日志（带 Context）
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Logger.Error(msg, fields...)
}

// WarnContext 记录 Warn 级别日志（带 Context）
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Logger.Warn(msg, fields...)
}

// DebugContext 记录 Debug 级别日志（带 Context）
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Logger.Debug(msg, fields...)
}

// 全局日志器
var defaultLogger *Logger

func init() {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "opengeo"
	}
	defaultLogger = New(serviceName)
}

// Default 获取默认日志器
func Default() *Logger {
	return defaultLogger
}

// SetDefault 设置默认日志器
func SetDefault(l *Logger) {
	defaultLogger = l
}
