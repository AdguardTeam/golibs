package redisutil_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testTimeout is the common timeout for tests.
const testTimeout = 1 * time.Second

// testPortEnvVarName is the environment variable name the presence and value of
// which define whether to run depending tests and on which port Redis server is
// running.
const testPortEnvVarName = "TEST_REDIS_PORT"

// Key and value constants.
const (
	testKey   = "key"
	testValue = "value"
)

// Redis pool configuration constants for common tests.
const (
	testIdleTimeout     = 30 * time.Second
	testMaxConnLifetime = 30 * time.Second

	testMaxActive = 10
	textMaxIdle   = 3

	testDBIndex = 15
)

// testLogger is the common logger for tests.
var testLogger = slogutil.NewDiscardLogger()

// newIntegrationDialer returns a *redisutil.DefaultDialer for tests or skips
// the test if [testPortEnvVarName] is not set.  It selects a database at
// [testDBIndex] and flushes it after the test.
func newIntegrationDialer(tb testing.TB) (d *redisutil.DefaultDialer) {
	tb.Helper()

	portStr := os.Getenv(testPortEnvVarName)
	if portStr == "" {
		tb.Skipf("skipping; %s is not set", testPortEnvVarName)
	}

	port64, err := strconv.ParseUint(portStr, 10, 16)
	require.NoError(tb, err)

	d, err = redisutil.NewDefaultDialer(&redisutil.DefaultDialerConfig{
		Addr: &netutil.HostPort{
			Host: "localhost",
			Port: uint16(port64),
		},
		DBIndex: testDBIndex,
	})
	require.NoError(tb, err)

	testutil.CleanupAndRequireSuccess(tb, func() (cleanupErr error) {
		ctx := testutil.ContextWithTimeout(tb, testTimeout)
		c, cleanupErr := d.DialContext(ctx)
		require.NoError(tb, cleanupErr)
		testutil.CleanupAndRequireSuccess(tb, c.Close)

		okStr, cleanupErr := redis.String(c.Do(redisutil.CmdFLUSHDB, redisutil.ParamSYNC))
		require.NoError(tb, cleanupErr)

		assert.Equal(tb, redisutil.RespOK, okStr)

		return cleanupErr
	})

	return d
}

// newIntegrationPool returns a *redisutil.DefaultPool for tests or skips the
// test if [testPortEnvVarName] is not set.  It selects a database at
// [testDBIndex] and flushes it after the test.
func newIntegrationPool(tb testing.TB) (p *redisutil.DefaultPool) {
	tb.Helper()

	dialer := newIntegrationDialer(tb)
	p, err := redisutil.NewDefaultPool(&redisutil.DefaultPoolConfig{
		Logger:          testLogger,
		Dialer:          dialer,
		MaxConnLifetime: testMaxConnLifetime,
		IdleTimeout:     testIdleTimeout,
		MaxActive:       testMaxActive,
		MaxIdle:         textMaxIdle,
		Wait:            true,
	})
	require.NoError(tb, err)

	return p
}
