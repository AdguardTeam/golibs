// Package fakeslog contains fake implementations of interfaces from package
// log/slog.
package fakeslog

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/testutil"
)

// Handler is a [slog.Handler] for tests.
type Handler struct {
	OnEnabled   func(ctx context.Context, level slog.Level) (ok bool)
	OnHandle    func(ctx context.Context, record slog.Record) (err error)
	OnWithAttrs func(attrs []slog.Attr) (handler slog.Handler)
	OnWithGroup func(name string) (handler slog.Handler)
}

// NewHandler returns a new *Handler all methods of which panic.
func NewHandler() (h *Handler) {
	return &Handler{
		OnEnabled: func(ctx context.Context, level slog.Level) (ok bool) {
			panic(testutil.UnexpectedCall(ctx, level))
		},
		OnHandle: func(ctx context.Context, record slog.Record) (err error) {
			panic(testutil.UnexpectedCall(ctx, record))
		},
		OnWithAttrs: func(attrs []slog.Attr) (handler slog.Handler) {
			panic(testutil.UnexpectedCall(attrs))
		},
		OnWithGroup: func(name string) (handler slog.Handler) {
			panic(testutil.UnexpectedCall(name))
		},
	}
}

// type check
var _ slog.Handler = (*Handler)(nil)

// Enabled implements the [slog.Handler] interface for *Handler.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) (ok bool) {
	return h.OnEnabled(ctx, level)
}

// Handle implements the [slog.Handler] interface for *Handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) (err error) {
	return h.OnHandle(ctx, record)
}

// WithAttrs implements the [slog.Handler] interface for *Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) (handler slog.Handler) {
	return h.OnWithAttrs(attrs)
}

// WithGroup implements the [slog.Handler] interface for *Handler.
func (h *Handler) WithGroup(name string) (handler slog.Handler) {
	return h.OnWithGroup(name)
}
