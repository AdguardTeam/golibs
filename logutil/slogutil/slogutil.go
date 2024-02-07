// Package slogutil contains extensions and utilities for package log/slog from
// the standard library.
package slogutil

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
)

// Additional key constants.
const (
	KeyPrefix = "prefix"
	KeyError  = "err"
)

// Config contains the configuration for a logger.
type Config struct {
	// Output is the output destination.  If not set, [os.Stdout] is used.
	Output io.Writer

	// Format is the format for the logs.  If not set, [FormatDefault] is used.
	// If set, it must be valid.
	Format Format

	// AddTimestamp, if true, adds a timestamp to every record.
	AddTimestamp bool

	// Verbose, if true, enables verbose logging.
	Verbose bool
}

// New creates a slog logger with the given parameters.  If c is nil, the
// defaults are used.
//
// NOTE: If c.Format is [FormatAdGuardLegacy], the legacy logger parameters,
// such as output, should be set separately.
func New(c *Config) (l *slog.Logger) {
	if c == nil {
		c = &Config{
			Output: os.Stdout,
			Format: FormatDefault,
		}
	}

	lvl := slog.LevelInfo
	if c.Verbose {
		lvl = slog.LevelDebug
	}

	format := or(c.Format, FormatDefault)
	output := or[io.Writer](c.Output, os.Stdout)
	if format == FormatDefault {
		// Fast path for the default handler.
		return newDefault(output, lvl, c.AddTimestamp)
	}

	var replaceAttr func(groups []string, a slog.Attr) (res slog.Attr)
	if !c.AddTimestamp {
		replaceAttr = RemoveTime
	}

	var h slog.Handler
	switch format {
	case FormatAdGuardLegacy:
		h = NewAdGuardLegacyHandler(lvl)
	case FormatJSON:
		h = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	case FormatJSONHybrid:
		h = NewJSONHybridHandler(output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	case FormatText:
		h = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	default:
		panic(&BadFormatError{
			Format: string(format),
		})
	}

	return slog.New(h)
}

// or returns defVal if val is its zero value.
//
// TODO(a.garipov): Replace with [cmp.Or] in Go 1.22.
func or[T comparable](val, defVal T) (res T) {
	var zero T
	if val == zero {
		return defVal
	}

	return val
}

// newDefault returns a new default slog logger set up with the given options.
func newDefault(output io.Writer, lvl slog.Level, addTimestamp bool) (l *slog.Logger) {
	h := NewLevelHandler(lvl, slog.Default().Handler())
	log.SetOutput(output)
	if addTimestamp {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	} else {
		log.SetFlags(0)
	}

	return slog.New(h)
}

// RemoveTime is a function for [slog.HandlerOptions.ReplaceAttr] that removes
// the "time" attribute.
func RemoveTime(groups []string, a slog.Attr) (res slog.Attr) {
	if len(groups) > 0 {
		return a
	}

	if a.Key == "time" {
		return slog.Attr{}
	}

	return a
}

// PrintStack logs the stacktrace into l on the given level.
func PrintStack(ctx context.Context, l *slog.Logger, lvl slog.Level) {
	stack := bytes.Split(debug.Stack(), []byte{'\n'})
	for i, line := range stack {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			l.Log(ctx, lvl, "stack", "i", i, "line", line)
		}
	}
}

// RecoverAndLog is a deferred helper that recovers from a panic and logs the
// panic value into l along with the stacktrace.
func RecoverAndLog(ctx context.Context, l *slog.Logger) {
	v := recover()
	if v == nil {
		return
	}

	var args []any
	if err, ok := v.(error); ok {
		args = []any{KeyError, err}
	} else {
		args = []any{"value", v}
	}

	l.ErrorContext(ctx, "recovered from panic", args...)
	PrintStack(ctx, l, slog.LevelError)
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
