// Package slogutil contains extensions and utilities for package log/slog from
// the standard library.
package slogutil

import (
	"bytes"
	"cmp"
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
)

// Additional or renamed key constants.
const (
	KeyError  = "err"
	KeyPrefix = "prefix"

	KeyMessage = slog.MessageKey
	KeySource  = slog.SourceKey
	KeyTime    = slog.TimeKey
	keyLevel   = slog.LevelKey
)

// Config contains the configuration for a logger.
type Config struct {
	// Output is the output destination.  If not set, [os.Stdout] is used.
	Output io.Writer

	// Format is the format for the logs.  If not set, [FormatDefault] is used.
	// If set, it must be valid.
	Format Format

	// Level is the minimum record level that will be logged.
	Level slog.Level

	// AddTimestamp, if true, adds a timestamp to every record.
	AddTimestamp bool
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

	lvl := c.Level
	format := cmp.Or(c.Format, FormatDefault)
	output := cmp.Or[io.Writer](c.Output, os.Stdout)
	if format == FormatDefault {
		// Fast path for the default handler.
		return newDefault(output, lvl, c.AddTimestamp)
	}

	replaceAttr := newReplaceAttr(!c.AddTimestamp)

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

// newDefault returns a new default slog logger set up with the given options.
//
// TODO(d.kolyshev): Replace log level name for [LevelTrace].
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

// newReplaceAttr is a function that returns [slog.HandlerOptions.ReplaceAttr]
// function for provided parameters.
func newReplaceAttr(removeTime bool) func(groups []string, a slog.Attr) (res slog.Attr) {
	if !removeTime {
		return ReplaceLevel
	}

	return func(groups []string, a slog.Attr) (res slog.Attr) {
		return ReplaceLevel(groups, RemoveTime(groups, a))
	}
}

// traceAttrValue is a [LevelTrace] value under the [slog.LevelKey] key.
var traceAttrValue = slog.StringValue("TRACE")

// ReplaceLevel is a function for [slog.HandlerOptions.ReplaceAttr] that adds
// [LevelTrace] custom name for level attribute.
func ReplaceLevel(groups []string, a slog.Attr) (res slog.Attr) {
	if len(groups) > 0 {
		return a
	}

	if a.Key == slog.LevelKey {
		lvl := a.Value.Any().(slog.Level)
		if lvl == LevelTrace {
			a.Value = traceAttrValue
		}
	}

	return a
}

// RemoveTime is a function for [slog.HandlerOptions.ReplaceAttr] that removes
// the "time" attribute.
func RemoveTime(groups []string, a slog.Attr) (res slog.Attr) {
	if len(groups) > 0 {
		return a
	}

	if a.Key == KeyTime {
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

// PrintLines splits s and logs the lines to l using the given level, including
// empty lines.
func PrintLines(ctx context.Context, l *slog.Logger, lvl slog.Level, msg, s string) {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		l.Log(ctx, lvl, msg, "line_num", i+1, "line", line)
	}
}

// PrintByteLines splits b and logs the lines to l using the given level,
// including empty lines.
func PrintByteLines(ctx context.Context, l *slog.Logger, lvl slog.Level, msg string, b []byte) {
	lines := bytes.Split(b, []byte{'\n'})
	for i, line := range lines {
		l.Log(ctx, lvl, msg, "line_num", i+1, "line", line)
	}
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
