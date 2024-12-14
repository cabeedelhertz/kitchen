package logging

import (
	"context"

	"go.uber.org/zap"
)

// NB(MLH) this type alias provides uniqueness
type key int

const (
	loggerKey key = iota
)

// FromContext gets a Logger from the provided context. If no Logger is present,
// this will return the global logger instance
func FromContext(ctx context.Context) *zap.Logger {
	if t, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return t
	}
	return zap.L()
}

// NewContext puts a logger into a context
func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
