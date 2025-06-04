package testutil

import (
	"context"
	"testing"
	"time"
)

// ContextWithTimeout is a helper that creates a new context with timeout and
// registers ctx's cancellation with [testing.TB.Cleanup].
//
// TODO(a.garipov): Consider using [testing.TB.Context].
func ContextWithTimeout(tb testing.TB, timeout time.Duration) (ctx context.Context) {
	tb.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	tb.Cleanup(cancel)

	return ctx
}
