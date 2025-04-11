package redisutil_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeredis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleChecker_TestConnection(t *testing.T) {
	respOthers := []any{
		[]any{[]byte("127.0.0.1"), []byte("1234"), []byte("1")},
		[]any{[]byte("127.0.0.2"), []byte("1234"), []byte("1")},
	}

	require.True(t, t.Run("match", func(t *testing.T) {
		conn := fakeredis.NewConn()

		conn.OnDo = func(cmdName string, args ...any) (reply any, err error) {
			assert.Equal(t, redisutil.CmdROLE, cmdName)

			resp := []any{[]byte(redisutil.RoleStringMaster), 1, respOthers}

			return resp, nil
		}

		rc, err := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
		})
		require.NoError(t, err)

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err = rc.TestConnection(ctx, conn, time.Now())
		require.NoError(t, err)
	}))

	require.True(t, t.Run("no_match", func(t *testing.T) {
		conn := fakeredis.NewConn()

		conn.OnDo = func(cmdName string, args ...any) (reply any, err error) {
			assert.Equal(t, redisutil.CmdROLE, cmdName)

			resp := []any{[]byte(redisutil.RoleStringSlave), 1, respOthers}

			return resp, nil
		}

		rc, err := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
		})
		require.NoError(t, err)

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err = rc.TestConnection(ctx, conn, time.Now())

		wantErrMsg := fmt.Sprintf(
			"testing conn: want role %q, got %q",
			redisutil.RoleStringMaster,
			redisutil.RoleStringSlave,
		)
		testutil.AssertErrorMsg(t, wantErrMsg, err)
	}))
}

func TestRoleChecker_TestConnection_integration(t *testing.T) {
	d := newIntegrationDialer(t)

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	c, err := d.DialContext(ctx)
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, c.Close)

	require.True(t, t.Run("match", func(t *testing.T) {
		rc, testErr := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
		})
		require.NoError(t, testErr)

		testCtx := testutil.ContextWithTimeout(t, testTimeout)
		testErr = rc.TestConnection(testCtx, c, time.Now())
		require.NoError(t, testErr)
	}))

	require.True(t, t.Run("no_match", func(t *testing.T) {
		rc, testErr := redisutil.NewRoleChecker(&redisutil.RoleCheckerConfig{
			Logger: testLogger,
			Role:   redisutil.RoleSlave,
		})
		require.NoError(t, testErr)

		testCtx := testutil.ContextWithTimeout(t, testTimeout)
		testErr = rc.TestConnection(testCtx, c, time.Now())

		wantErrMsg := fmt.Sprintf(
			"testing conn: want role %q, got %q",
			redisutil.RoleStringSlave,
			redisutil.RoleStringMaster,
		)
		testutil.AssertErrorMsg(t, wantErrMsg, testErr)
	}))
}
