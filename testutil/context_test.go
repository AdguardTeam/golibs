package testutil_test

import (
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextWithTimeout(t *testing.T) {
	t.Parallel()

	const defaultTimeout = 10 * time.Second

	now := time.Now()

	// If there is no deadline, use the default timeout of 10 seconds.
	// Otherwise, use half of the time undil the deadline as the new timeout.
	parent := t.Context()
	dl, ok := parent.Deadline()
	if !ok {
		dl = now.Add(defaultTimeout)
	}

	timeout := dl.Sub(now) / 2

	wantDL := time.Now().Add(timeout)
	ctx := testutil.ContextWithTimeout(t, timeout)
	require.NotNil(t, ctx)

	gotDL, ok := ctx.Deadline()
	require.True(t, ok)

	assert.WithinDuration(t, wantDL, gotDL, timeout/100)
}
