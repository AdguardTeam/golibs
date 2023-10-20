//go:build go1.21

package slogutil

import (
	"context"
	"log/slog"
)

// DiscardHandler ignores all messages.
type DiscardHandler struct{}

// type check
var _ slog.Handler = DiscardHandler{}

// Enabled implements the [slog.Handler] interface for DiscardHandler.  It
// always returns false.
func (h DiscardHandler) Enabled(_ context.Context, _ slog.Level) (ok bool) { return false }

// Handle implements the [slog.Handler] interface for DiscardHandler.  It always
// returns nil.
func (h DiscardHandler) Handle(_ context.Context, _ slog.Record) (err error) {
	return nil
}

// WithAttrs implements the [slog.Handler] interface for DiscardHandler.  It
// always returns h.
func (h DiscardHandler) WithAttrs(_ []slog.Attr) (res slog.Handler) {
	return h
}

// WithGroup implements the [slog.Handler] interface for DiscardHandler.  It
// always returns h.
func (h DiscardHandler) WithGroup(_ string) (res slog.Handler) {
	return h
}

// NewDiscardLogger returns a new logger that uses [DiscardHandler].
func NewDiscardLogger() (l *slog.Logger) {
	return slog.New(DiscardHandler{})
}
