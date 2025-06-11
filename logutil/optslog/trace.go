package optslog

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// Trace1 is an optimized version of [slog.Logger.Log] that prevents it from
// allocating when debugging is not necessary.
func Trace1[T1 any](ctx context.Context, l *slog.Logger, msg, name1 string, arg1 T1) {
	if l.Enabled(ctx, slogutil.LevelTrace) {
		l.Log(ctx, slogutil.LevelTrace, msg, name1, arg1)
	}
}

// Trace2 is an optimized version of [slog.Logger.Log] that prevents it from
// allocating when debugging is not necessary.
func Trace2[T1, T2 any](
	ctx context.Context,
	l *slog.Logger,
	msg,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
) {
	if l.Enabled(ctx, slogutil.LevelTrace) {
		l.Log(ctx, slogutil.LevelTrace, msg, name1, arg1, name2, arg2)
	}
}

// Trace3 is an optimized version of [slog.Logger.Log] that prevents it from
// allocating when debugging is not necessary.
func Trace3[T1, T2, T3 any](
	ctx context.Context,
	l *slog.Logger,
	msg,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
	name3 string, arg3 T3,
) {
	if l.Enabled(ctx, slogutil.LevelTrace) {
		l.Log(ctx, slogutil.LevelTrace, msg, name1, arg1, name2, arg2, name3, arg3)
	}
}

// Trace4 is an optimized version of [slog.Logger.Log] that prevents it from
// allocating when debugging is not necessary.
func Trace4[T1, T2, T3, T4 any](
	ctx context.Context,
	l *slog.Logger,
	msg,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
	name3 string, arg3 T3,
	name4 string, arg4 T4,
) {
	if l.Enabled(ctx, slogutil.LevelTrace) {
		l.Log(ctx, slogutil.LevelTrace, msg, name1, arg1, name2, arg2, name3, arg3, name4, arg4)
	}
}
