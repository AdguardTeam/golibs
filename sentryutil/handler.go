package sentryutil

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// traceMsgHandler is the [slog.Handler] that prints unnecessary messages at
// [slogutil.LevelTrace].
//
// TODO(a.garipov):  Consider ways to generalize and move to slogutil.
type traceMsgHandler struct {
	handler slog.Handler
}

// type check
var _ slog.Handler = (*traceMsgHandler)(nil)

// Enabled implements the [slog.Handler] interface for *traceMsgHandler.
func (h *traceMsgHandler) Enabled(ctx context.Context, lvl slog.Level) (ok bool) {
	return h.handler.Enabled(ctx, lvl)
}

// traceMsg is a message that is commonly issued by Sentry clients.  Since it
// doesn't contain any useful information in most cases, it should be printed at
// [slogutil.LevelTrace].
const traceMsg = "Dropping transaction: EnableTracing is set to false"

// Handle implements the [slog.Handler] interface for *traceMsgHandler.
func (h *traceMsgHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	if r.Message == traceMsg {
		if !h.handler.Enabled(ctx, slogutil.LevelTrace) {
			return nil
		}

		r.Level = slogutil.LevelTrace
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs implements the [slog.Handler] interface for *traceMsgHandler.
func (h *traceMsgHandler) WithAttrs(attrs []slog.Attr) (res slog.Handler) {
	return h.handler.WithAttrs(attrs)
}

// WithGroup implements the [slog.Handler] interface for *traceMsgHandler.
func (h *traceMsgHandler) WithGroup(g string) (res slog.Handler) {
	return h.handler.WithGroup(g)
}
