package slogutil

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/syncutil"
)

// JSONHybridHandler is a hybrid JSON-and-text [slog.Handler] more suitable for
// stricter environments.  It guarantees that the only properties present in the
// resulting objects are "level", "msg", "time", and "source", depending on the
// options.  All other attributes are packed into the "msg" property using the
// same format as [slogutil.TextHandler].
//
// NOTE: [JSONHybridHandler.WithGroup] is not currently supported and panics.
//
// Example of output:
//
//	{"time":"2023-12-01T12:34:56.789Z","level":"INFO","msg":"listening; attrs: prefix=websvc url=http://127.0.0.1:8080"}
type JSONHybridHandler struct {
	json        *slog.JSONHandler
	attrPool    *syncutil.Pool[[]slog.Attr]
	bufTextPool *syncutil.Pool[bufferedTextHandler]
	textAttrs   []slog.Attr
}

const (
	// initAttrsLenEst is the estimation used to set the initial length of
	// attribute slices.
	initAttrsLenEst = 2

	// initLineLenEst is the estimation used to set the initial sizes of
	// log-line buffers.
	initLineLenEst = 256
)

// NewJSONHybridHandler creates a new properly initialized *JSONHybridHandler.
// opts are used for the underlying JSON handler.
func NewJSONHybridHandler(w io.Writer, opts *slog.HandlerOptions) (h *JSONHybridHandler) {
	return &JSONHybridHandler{
		json:     slog.NewJSONHandler(w, opts),
		attrPool: syncutil.NewSlicePool[slog.Attr](initAttrsLenEst),
		bufTextPool: syncutil.NewPool(func() (bufTextHdlr *bufferedTextHandler) {
			return newBufferedTextHandler(initLineLenEst)
		}),
		textAttrs: nil,
	}
}

// type check
var _ slog.Handler = (*JSONHybridHandler)(nil)

// Enabled implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) Enabled(ctx context.Context, level slog.Level) (ok bool) {
	return h.json.Enabled(ctx, level)
}

// Handle implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	bufTextHdlr := h.bufTextPool.Get()
	defer h.bufTextPool.Put(bufTextHdlr)

	bufTextHdlr.reset()

	_, _ = bufTextHdlr.buffer.WriteString(r.Message)

	numAttrs := r.NumAttrs() + len(h.textAttrs)
	if numAttrs > 0 {
		_, _ = bufTextHdlr.buffer.WriteString("; attrs: ")
	}

	textAttrsPtr := h.attrPool.Get()
	defer h.attrPool.Put(textAttrsPtr)

	*textAttrsPtr = (*textAttrsPtr)[:0]
	r.Attrs(func(a slog.Attr) (cont bool) {
		*textAttrsPtr = append(*textAttrsPtr, a)

		return true
	})

	textRec := slog.NewRecord(time.Time{}, r.Level, "", 0)
	textRec.AddAttrs(h.textAttrs...)
	textRec.AddAttrs(*textAttrsPtr...)

	err = bufTextHdlr.handler.Handle(ctx, textRec)
	if err != nil {
		return fmt.Errorf("handling text for msg: %w", err)
	}

	msgForJSON := bufTextHdlr.buffer.String()

	// Remove newline.
	msgForJSON = msgForJSON[:len(msgForJSON)-1]

	return h.json.Handle(ctx, slog.NewRecord(r.Time, r.Level, msgForJSON, r.PC))
}

// WithAttrs implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) WithAttrs(attrs []slog.Attr) (res slog.Handler) {
	return &JSONHybridHandler{
		json:        h.json,
		attrPool:    h.attrPool,
		bufTextPool: h.bufTextPool,
		textAttrs:   append(slices.Clip(h.textAttrs), attrs...),
	}
}

// WithGroup implements the [slog.Handler] interface for *JSONHybridHandler.
//
// NOTE: It is currently not supported and panics.
//
// TODO(a.garipov): Support groups, see
// https://github.com/golang/example/blob/master/slog-handler-guide/README.md.
func (h *JSONHybridHandler) WithGroup(g string) (res slog.Handler) {
	panic(errors.Error("slogutil: JSONHybridHandler.WithGroup is not supported"))
}
