package redisutil_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func TestRoleChecker_TestConnection_integration(t *testing.T) {
	d := newDialer(t)

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	c, err := d.DialContext(ctx)
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, c.Close)

	require.True(t, t.Run("match", func(t *testing.T) {
		rc, testErr := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
		})
		require.NoError(t, testErr)

		ctx = testutil.ContextWithTimeout(t, testTimeout)
		testErr = rc.TestConnection(ctx, c, time.Now())
		require.NoError(t, testErr)
	}))

	require.True(t, t.Run("no_match", func(t *testing.T) {
		rc, testErr := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
			Role:   redisutil.RoleSlave,
		})
		require.NoError(t, testErr)

		ctx = testutil.ContextWithTimeout(t, testTimeout)
		testErr = rc.TestConnection(ctx, c, time.Now())

		wantErrMsg := fmt.Sprintf(
			"testing conn: want role %q, got %q",
			redisutil.RoleStringSlave,
			redisutil.RoleStringMaster,
		)
		testutil.AssertErrorMsg(t, wantErrMsg, testErr)
	}))
}
