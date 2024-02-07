package slogutil

import (
	"context"
	"log/slog"
)

// A LevelHandler wraps a Handler with an Enabled method that returns false for
// levels below a minimum.
//
// See https://cs.opensource.google/go/x/exp/+/master:slog/example_level_handler_test.go.
type LevelHandler struct {
	level   slog.Leveler
	handler slog.Handler
}

// NewLevelHandler returns a LevelHandler with the given level.  All methods
// except Enabled delegate to h.
func NewLevelHandler(level slog.Leveler, h slog.Handler) (lh *LevelHandler) {
	// As an optimization, avoid chains of LevelHandlers.
	lh, ok := h.(*LevelHandler)
	if ok {
		h = lh.Handler()
	}

	return &LevelHandler{level, h}
}

// type check
var _ slog.Handler = (*LevelHandler)(nil)

// Enabled implements the [slog.Handler] interface for *LevelHandler.  It
// reports whether level is as high as h's level.
func (h *LevelHandler) Enabled(_ context.Context, level slog.Level) (ok bool) {
	return level >= h.level.Level()
}

// Handle implements the [slog.Handler] interface for *LevelHandler.
func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements the [slog.Handler] interface for *LevelHandler.
func (h *LevelHandler) WithAttrs(attrs []slog.Attr) (res slog.Handler) {
	return NewLevelHandler(h.level, h.handler.WithAttrs(attrs))
}

// WithGroup implements the [slog.Handler] interface for *LevelHandler.
func (h *LevelHandler) WithGroup(name string) (res slog.Handler) {
	return NewLevelHandler(h.level, h.handler.WithGroup(name))
}

// Handler returns the slog.Handler wrapped by h.
func (h *LevelHandler) Handler() (unwrapped slog.Handler) {
	return h.handler
}
