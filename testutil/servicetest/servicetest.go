// Package servicetest contains test utilities for package service.
package servicetest

import (
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/service"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

// RequireRun is a helper that starts a service and adds its shutdown to tb's
// cleanup.
func RequireRun(tb testing.TB, s service.Interface, timeout time.Duration) {
	tb.Helper()

	ctx := testutil.ContextWithTimeout(tb, timeout)
	require.NoError(tb, s.Start(ctx))

	testutil.CleanupAndRequireSuccess(tb, func() (err error) {
		ctx = testutil.ContextWithTimeout(tb, timeout)

		return s.Shutdown(ctx)
	})
}
