package slogutil

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	aglog "github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/golibs/mathutil"
	"github.com/AdguardTeam/golibs/syncutil"
)

// AdGuardLegacyHandler is a text [slog.Handler] that uses package
// github.com/AdguardTeam/golibs/log for output.  It is a legacy handler that
// will be removed in a future version.
//
// The attribute with the name [KeyPrefix] is handled separately.  Like in the
// legacy logging conventions, it is prepended to the message.  There should
// only be one such attribute between the logger and the message.
//
// Example of output:
//
//	12#34 [debug] debug with no attributes
//	12#34 [info] hdlr: info with prefix and attrs number=123
type AdGuardLegacyHandler struct {
	level       slog.Leveler
	attrPool    *syncutil.Pool[[]slog.Attr]
	bufTextPool *syncutil.Pool[bufferedTextHandler]
	attrs       []slog.Attr
}

// NewAdGuardLegacyHandler creates a new properly initialized
// *AdGuardLegacyHandler.  Output, level, and other flags should be set in the
// legacy logging package.  lvl is used for [AdGuardLegacyHandler.Enabled].
func NewAdGuardLegacyHandler(lvl slog.Leveler) (h *AdGuardLegacyHandler) {
	return &AdGuardLegacyHandler{
		level:    lvl,
		attrPool: syncutil.NewSlicePool[slog.Attr](initAttrsLenEst),
		bufTextPool: syncutil.NewPool(func() (bufTextHdlr *bufferedTextHandler) {
			return newBufferedTextHandler(initLineLenEst)
		}),
		attrs: nil,
	}
}

// type check
var _ slog.Handler = (*AdGuardLegacyHandler)(nil)

// Enabled implements the [slog.Handler] interface for *AdGuardLegacyHandler.
func (h *AdGuardLegacyHandler) Enabled(ctx context.Context, lvl slog.Level) (ok bool) {
	return lvl >= h.level.Level()
}

// Handle implements the [slog.Handler] interface for *AdGuardLegacyHandler.
func (h *AdGuardLegacyHandler) Handle(ctx context.Context, r slog.Record) (err error) {
	bufTextHdlr := h.bufTextPool.Get()
	defer h.bufTextPool.Put(bufTextHdlr)

	bufTextHdlr.reset()

	_, _ = bufTextHdlr.buffer.WriteString(r.Message)

	numAttrs := r.NumAttrs() + len(h.attrs)

	textAttrsPtr := h.attrPool.Get()
	defer h.attrPool.Put(textAttrsPtr)

	*textAttrsPtr = (*textAttrsPtr)[:0]

	var prefix string
	*textAttrsPtr, prefix = h.appendAttrs(*textAttrsPtr, r)

	textRec := slog.NewRecord(time.Time{}, r.Level, "", 0)
	textRec.AddAttrs(*textAttrsPtr...)

	// Append a space to separate the message from the key-value attributes, but
	// do not count the prefix attribute if there is one to prevent an
	// unnecessary space at the end of line.
	numAttrs -= mathutil.BoolToNumber[int](prefix != "")
	if numAttrs > 0 {
		_, _ = bufTextHdlr.buffer.WriteString(" ")
	}

	err = bufTextHdlr.handler.Handle(ctx, textRec)
	if err != nil {
		return fmt.Errorf("handling text for msg: %w", err)
	}

	logFunc, warnStr, err := logFuncForLevel(r.Level)
	if err != nil {
		return fmt.Errorf("legacy handler: handling: %w", err)
	}

	// Remove newline.
	msgForLegacy := bufTextHdlr.buffer.String()
	msgForLegacy = msgForLegacy[:len(msgForLegacy)-1]

	if prefix == "" {
		logFunc("%s%s", warnStr, msgForLegacy)
	} else {
		logFunc("%s: %s%s", prefix, warnStr, msgForLegacy)
	}

	return nil
}

// appendAttrs appends both the handler's and the record's attributes to attrs
// and returns it.  Additionally, it looks up the prefix in the attributes.
func (h *AdGuardLegacyHandler) appendAttrs(
	attrs []slog.Attr,
	r slog.Record,
) (res []slog.Attr, prefix string) {
	res = attrs

	for _, a := range h.attrs {
		if a.Key == KeyPrefix {
			prefix = a.Value.String()
		} else {
			res = append(res, a)
		}
	}

	var recPrefix string
	r.Attrs(func(a slog.Attr) (cont bool) {
		if a.Key == KeyPrefix {
			recPrefix = a.Value.String()
		} else {
			res = append(res, a)
		}

		return true
	})

	if recPrefix == "" {
		return res, prefix
	} else if prefix == "" {
		return res, recPrefix
	}

	aglog.Debug(
		"legacy logger: got prefix %q in record for logger with prefix %q",
		recPrefix,
		prefix,
	)

	return res, prefix
}

// logFunction is a function for printf-like logging functions.
type logFunction func(format string, args ...any)

// logFuncForLevel returns the function for logging, and an optional warning
// string for the given level.  If the level is not recognized, an error is
// returned.
func logFuncForLevel(lvl slog.Level) (logFunc logFunction, warnStr string, err error) {
	switch lvl {
	case LevelTrace:
		return aglog.Debug, "trace: ", nil
	case slog.LevelDebug:
		return aglog.Debug, "", nil
	case slog.LevelInfo:
		return aglog.Info, "", nil
	case slog.LevelWarn:
		return aglog.Info, "warning: ", nil
	case slog.LevelError:
		return aglog.Error, "", nil
	default:
		return nil, "", fmt.Errorf("unsupported level %v", lvl)
	}
}

// WithAttrs implements the [slog.Handler] interface for *AdGuardLegacyHandler.
func (h *AdGuardLegacyHandler) WithAttrs(attrs []slog.Attr) (res slog.Handler) {
	return &AdGuardLegacyHandler{
		level:       h.level,
		attrPool:    h.attrPool,
		bufTextPool: h.bufTextPool,
		attrs:       append(slices.Clip(h.attrs), attrs...),
	}
}

// WithGroup implements the [slog.Handler] interface for *AdGuardLegacyHandler.
//
// NOTE: It is currently not supported and panics.
//
// TODO(a.garipov): Support groups, see
// https://github.com/golang/example/blob/master/slog-handler-guide/README.md.
func (h *AdGuardLegacyHandler) WithGroup(g string) (res slog.Handler) {
	panic(errors.Error("slogutil: AdGuardLegacyHandler.WithGroup is not supported"))
}
