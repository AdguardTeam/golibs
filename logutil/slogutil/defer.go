package slogutil

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/AdguardTeam/golibs/osutil"
)

// CloseAndLog is a convenient helper to log errors returned by closer.  The
// point is to not lose information from deferred Close calls.  The error is
// logged with the specified logging level.  l must not be nil.
//
// Instead of:
//
//	defer f.Close()
//
// You can now write:
//
//	defer slogutil.CloseAndLog(ctx, l, f, slog.LevelDebug)
//
// Note that if closer is nil, it is simply ignored.
func CloseAndLog(ctx context.Context, l *slog.Logger, closer io.Closer, lvl slog.Level) {
	if closer == nil {
		return
	}

	err := closer.Close()
	if err == nil {
		return
	}

	l.Log(ctx, lvl, "deferred closing", KeyError, err)
}

// RecoverAndLog is a deferred helper that recovers from a panic and logs the
// panic value into l along with the stacktrace.  l must not be nil.
func RecoverAndLog(ctx context.Context, l *slog.Logger) {
	PrintRecovered(ctx, l, recover())
}

// RecoverAndLogDefault is like [RecoverAndLog] but tries to get the logger from
// ctx using [LoggerFromContext] and, if there is none, uses [slog.Default].
func RecoverAndLogDefault(ctx context.Context) {
	v := recover()
	if v == nil {
		return
	}

	l, ok := LoggerFromContext(ctx)
	if !ok {
		l = slog.Default()
	}

	PrintRecovered(ctx, l, v)
}

// RecoverAndExit recovers a panic and, if there is one, logs it using l and
// exits with the given exit code.  l must not be nil.
func RecoverAndExit(ctx context.Context, l *slog.Logger, code osutil.ExitCode) {
	v := recover()
	if v == nil {
		return
	}

	PrintRecovered(ctx, l, v)

	os.Exit(code)
}

// PrintRecovered prints a message with the recovered value, if there is any, as
// well as the stack.  If v is nil, PrintRecovered does nothing.  Otherwise, l
// must not be nil.
func PrintRecovered(ctx context.Context, l *slog.Logger, v any) {
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
