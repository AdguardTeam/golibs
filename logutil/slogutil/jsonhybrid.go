//go:build go1.21

package slogutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/syncutil"
)

// JSONHybridHandler is a hybrid JSON-and-text handler [slog.Handler] more
// suitable for stricter environments.  It guarantees that the only properties
// present in the resulting objects are "level", "msg", "time", and "source",
// depending on the options.  All other attributes are packed into the "msg"
// property using the same format as [slogutil.TextHandler].
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

	textAttrs := (*textAttrsPtr)[:0]
	r.Attrs(func(a slog.Attr) (cont bool) {
		textAttrs = append(textAttrs, a)

		return true
	})

	textRec := slog.NewRecord(time.Time{}, r.Level, "", 0)
	textRec.AddAttrs(h.textAttrs...)
	textRec.AddAttrs(textAttrs...)

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

// bufferedTextHandler is a combination of one bytes buffer and a text handler
// that writes to it.
type bufferedTextHandler struct {
	buffer  *bytes.Buffer
	handler *slog.TextHandler
}

// newBufferedTextHandler returns a new bufferedTextHandler with the given
// buffer length.
func newBufferedTextHandler(l int) (h *bufferedTextHandler) {
	buf := bytes.NewBuffer(make([]byte, 0, l))

	return &bufferedTextHandler{
		buffer:  buf,
		handler: slog.NewTextHandler(buf, textHandlerOpts),
	}
}

// textHandlerOpts are the options used by buffered text handlers of JSON hybrid
// handlers.
var textHandlerOpts = &slog.HandlerOptions{
	ReplaceAttr: removeTopLevel,
}

// removeTopLevel is a [slog.HandlerOptions.ReplaceAttr] function that removes
// "level", "msg", "time", and "source" attributes.
func removeTopLevel(groups []string, a slog.Attr) (res slog.Attr) {
	if len(groups) > 0 {
		return a
	}

	switch a.Key {
	case
		slog.LevelKey,
		slog.MessageKey,
		slog.TimeKey,
		slog.SourceKey:
		return slog.Attr{}
	default:
		return a
	}
}

// reset must be called before using h after retrieving it from a pool.
func (h *bufferedTextHandler) reset() {
	h.buffer.Reset()
}
