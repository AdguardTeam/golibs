package slogutil

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/errors"
)

// ctxKey is the type for context keys.
type ctxKey int

// Context key values.
const (
	ctxKeyLogger ctxKey = iota
)

// ContextWithLogger returns a new context with the given logger.
func ContextWithLogger(parent context.Context, l *slog.Logger) (ctx context.Context) {
	return context.WithValue(parent, ctxKeyLogger, l)
}

// LoggerFromContext returns a logger for this request, if any.
func LoggerFromContext(ctx context.Context) (l *slog.Logger, ok bool) {
	v := ctx.Value(ctxKeyLogger)
	if v == nil {
		return nil, false
	}

	return v.(*slog.Logger), true
}

// MustLoggerFromContext returns a logger for this request and panics if there
// is no logger.
func MustLoggerFromContext(ctx context.Context) (l *slog.Logger) {
	l, ok := LoggerFromContext(ctx)
	if !ok {
		panic(errors.Error("slogutil: no logger in context"))
	}

	return l
}
