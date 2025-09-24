package testutil

import (
	"context"
	"testing"
	"time"
)

// ContextWithTimeout is a helper that creates a new context with timeout and
// registers ctx's cancellation with [testing.TB.Cleanup].
//
// TODO(a.garipov):  Consider creating a separate version using
// [testing.TB.Context].  Great care is necessary, as testing.TB.Context is
// cancelled before cleanup functions, so the context resulting from this new
// function should not be used in them.
func ContextWithTimeout(tb testing.TB, timeout time.Duration) (ctx context.Context) {
	tb.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	tb.Cleanup(cancel)

	return ctx
}
