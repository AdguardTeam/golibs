package slogutil

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/syncutil"
)

// JSONHybridHandler is a hybrid JSON-and-text [slog.Handler] more suitable for
// stricter environments.  It guarantees that the only properties present in the
// resulting objects are "severity" and "message".  All other attributes are
// packed into the "message" property using the same format as
// [slog.TextHandler].
//
// NOTE: [JSONHybridHandler.WithGroup] is not currently supported and panics.
//
// Example of output:
//
//	{"severity":"NORMAL","message":"time=2024-10-22T12:09:59.525+03:00 level=INFO msg=listening prefix=websvc server=http://127.0.0.1:8181"}
type JSONHybridHandler struct {
	level       slog.Leveler
	encoder     *json.Encoder
	bufTextPool *syncutil.Pool[bufferedTextHandler]

	// mu protects encoder.
	mu *sync.Mutex

	textAttrs []slog.Attr
}

// initLineLenEst is the estimation used to set the initial sizes of log-line
// buffers.
const initLineLenEst = 256

// NewJSONHybridHandler creates a new properly initialized *JSONHybridHandler.
// opts are used for the underlying text handler.
func NewJSONHybridHandler(w io.Writer, opts *slog.HandlerOptions) (h *JSONHybridHandler) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	lvl := slog.LevelInfo
	if opts != nil && opts.Level != nil {
		lvl = opts.Level.Level()
	}

	return &JSONHybridHandler{
		level:   lvl,
		encoder: enc,
		bufTextPool: syncutil.NewPool(func() (bufTextHdlr *bufferedTextHandler) {
			return newBufferedTextHandler(initLineLenEst, opts)
		}),
		mu:        &sync.Mutex{},
		textAttrs: nil,
	}
}

// type check
var _ slog.Handler = (*JSONHybridHandler)(nil)

// Enabled implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) Enabled(_ context.Context, level slog.Level) (ok bool) {
	return level >= h.level.Level()
}

// Handle implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	bufTextHdlr := h.bufTextPool.Get()
	defer h.bufTextPool.Put(bufTextHdlr)

	bufTextHdlr.reset()

	r.AddAttrs(h.textAttrs...)

	err = bufTextHdlr.handler.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("handling text for data: %w", err)
	}

	msg := byteString(bufTextHdlr.buffer.Bytes())

	// Remove newline.
	msg = msg[:len(msg)-1]
	data := newJSONHybridMessage(r.Level, msg)

	h.mu.Lock()
	defer h.mu.Unlock()

	return h.encoder.Encode(data)
}

// byteString optimizes memory allocations in [JSONHybridHandler.Handle].
type byteString []byte

// type check
var _ encoding.TextMarshaler = byteString{}

// MarshalText implements the [encoding.TextMarshaler] interface for byteString.
func (b byteString) MarshalText() (res []byte, err error) {
	return b, nil
}

// jsonHybridMessage represents the data structure for *JSONHybridHandler.
type jsonHybridMessage = struct {
	Severity string     `json:"severity"`
	Message  byteString `json:"message"`
}

// newJSONHybridMessage returns new properly initialized message.
func newJSONHybridMessage(lvl slog.Level, msg byteString) (m *jsonHybridMessage) {
	severity := "NORMAL"
	if lvl >= slog.LevelError {
		severity = "ERROR"
	}

	return &jsonHybridMessage{
		Severity: severity,
		Message:  msg,
	}
}

// WithAttrs implements the [slog.Handler] interface for *JSONHybridHandler.
func (h *JSONHybridHandler) WithAttrs(attrs []slog.Attr) (res slog.Handler) {
	return &JSONHybridHandler{
		level:       h.level,
		encoder:     h.encoder,
		bufTextPool: h.bufTextPool,
		mu:          h.mu,
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
