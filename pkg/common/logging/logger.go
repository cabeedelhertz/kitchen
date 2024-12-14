package logging

import (
	"kitchen/pkg/common/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LogFieldEnv
	LogFieldEnv = "dd.env"
	// LogFieldServiceName the log field holding the service name
	LogFieldServiceName = "dd.service"
	// LogFieldServiceVersion the log field holding the service version
	LogFieldServiceVersion = "dd.version"
	// LogFieldTraceID the logging field holding the trace id
	LogFieldTraceID = "dd.trace_id"
	// LogFieldSpanID the logging field holding the span id
	LogFieldSpanID = "dd.span_id"
)

var zcfg zap.Config

// Configure configures the logging system
func Configure(cfg config.Base) error {

	// Create zap configuration based on the env
	if cfg.Env != config.Local {
		zcfg = zap.NewProductionConfig()
	} else {
		zcfg = zap.NewDevelopmentConfig()
	}

	// If we are not in production, turn Development mode on
	if cfg.Env != config.Production {
		zcfg.Development = true
	}

	// Always disable sampling in production, we can sample downstream
	if cfg.Env == config.Production || cfg.DisableLogSampling {
		zcfg.Sampling = nil
	}

	// Override the log level if one is provided
	if cfg.LogLevel != "" {
		var level zapcore.Level
		if err := level.Set(cfg.LogLevel); err == nil {
			zcfg.Level.SetLevel(level)
		}
	}
	if cfg.DisableStackTraces {
		zcfg.DisableStacktrace = true
	}

	// Set the initial logger fields
	zcfg.InitialFields = map[string]interface{}{
		LogFieldEnv:            cfg.Env.String(),
		LogFieldServiceName:    cfg.ServiceName,
		LogFieldServiceVersion: cfg.ServiceVersion,
	}

	// Build a logger from our config
	logger, err := zcfg.Build()
	if err != nil {
		return err
	}

	// Replace the global logger with our logger
	zap.ReplaceGlobals(logger)

	return nil
}

// Config returns the zap global configuration
func Config() zap.Config {
	return zcfg
}

// NewLogger returns a new logger with the supplied name
func NewLogger(name string) *zap.Logger {
	return zap.L().Named(name)
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return zap.L().WithOptions(opts...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return zap.L().With(fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	zap.L().Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zap.Field) {
	zap.L().DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	zap.L().Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(msg, fields...)
}
