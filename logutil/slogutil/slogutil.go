//go:build go1.21

// Package slogutil contains extensions and utilities for package log/slog from
// the standard library.
package slogutil

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"runtime/debug"
)

// Additional key constants.
const (
	KeyPrefix = "prefix"
	KeyError  = "err"
)

// Config contains the configuration for a logger.
type Config struct {
	// Output is the output destination.  It must not be nil.
	Output io.Writer

	// Format is the format for the logs.  It must be valid.
	Format Format

	// AddTimestamp, if true, adds a timestamp to every record.
	AddTimestamp bool

	// Verbose, if true, enables verbose logging.
	Verbose bool
}

// New creates a slog logger with the given parameters.  c must not be nil and
// its fields must be valid.
func New(c *Config) (l *slog.Logger) {
	lvl := slog.LevelInfo
	if c.Verbose {
		lvl = slog.LevelDebug
	}

	if c.Format == FormatDefault {
		// Fast path for the default handler.
		return newDefault(c.Output, lvl, c.AddTimestamp)
	}

	var replaceAttr func(groups []string, a slog.Attr) (res slog.Attr)
	if !c.AddTimestamp {
		replaceAttr = RemoveTime
	}

	var h slog.Handler
	switch c.Format {
	case FormatJSON:
		h = slog.NewJSONHandler(c.Output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	case FormatJSONHybrid:
		h = NewJSONHybridHandler(c.Output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	case FormatText:
		h = slog.NewTextHandler(c.Output, &slog.HandlerOptions{
			Level:       lvl,
			ReplaceAttr: replaceAttr,
		})
	default:
		panic(&BadFormatError{
			Format: string(c.Format),
		})
	}

	return slog.New(h)
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
