package optslog

import (
	"context"
	"log/slog"
)

// Debug1 is an optimized version of [slog.Logger.DebugContext] that prevents it
// from allocating when debugging is not necessary.
func Debug1[T1 any](ctx context.Context, l *slog.Logger, msg, name1 string, arg1 T1) {
	if l.Enabled(ctx, slog.LevelDebug) {
		l.DebugContext(ctx, msg, name1, arg1)
	}
}

// Debug2 is an optimized version of [slog.Logger.DebugContext] that prevents it
// from allocating when debugging is not necessary.
func Debug2[T1, T2 any](
	ctx context.Context,
	l *slog.Logger,
	msg string,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
) {
	if l.Enabled(ctx, slog.LevelDebug) {
		l.DebugContext(ctx, msg, name1, arg1, name2, arg2)
	}
}

// Debug3 is an optimized version of [slog.Logger.DebugContext] that prevents it
// from allocating when debugging is not necessary.
func Debug3[T1, T2, T3 any](
	ctx context.Context,
	l *slog.Logger,
	msg string,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
	name3 string, arg3 T3,
) {
	if l.Enabled(ctx, slog.LevelDebug) {
		l.DebugContext(ctx, msg, name1, arg1, name2, arg2, name3, arg3)
	}
}

// Debug4 is an optimized version of [slog.Logger.DebugContext] that prevents it
// from allocating when debugging is not necessary.
func Debug4[T1, T2, T3, T4 any](
	ctx context.Context,
	l *slog.Logger,
	msg string,
	name1 string, arg1 T1,
	name2 string, arg2 T2,
	name3 string, arg3 T3,
	name4 string, arg4 T4,
) {
	if l.Enabled(ctx, slog.LevelDebug) {
		l.DebugContext(ctx, msg, name1, arg1, name2, arg2, name3, arg3, name4, arg4)
	}
}
