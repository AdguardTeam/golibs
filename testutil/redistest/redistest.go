// Package redistest contains utilities for testing Redis.
package redistest

import (
	"cmp"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

// EnvPort is the environment variable name the presence and value of which
// define whether to run depending tests and on which port Redis server is
// running.
const EnvPort = "TEST_REDIS_PORT"

// Redis pool configuration constants for common tests.
const (
	DBIndex = 15

	IdleTimeout     = 30 * time.Second
	MaxConnLifetime = 30 * time.Second

	MaxActive = 10
	MaxIdle   = 3
)

// Common constants for Redis tests.
const (
	Timeout = 1 * time.Second
)

// NewDialer returns a *redisutil.DefaultDialer for tests or skips the test if
// [EnvPort] is not set.  If c is nil, a default config is used.  If c is not
// nil:
//   - c.Addr is ignored and set to the value of [EnvPort] on host "localhost".
//   - If c.DBIndex is zero, it is set to [DBIndex].
func NewDialer(tb testing.TB, c *redisutil.DefaultDialerConfig) (d *redisutil.DefaultDialer) {
	tb.Helper()

	portStr := os.Getenv(EnvPort)
	if portStr == "" {
		tb.Skipf("skipping; %s is not set", EnvPort)
	}

	port64, err := strconv.ParseUint(portStr, 10, 16)
	require.NoError(tb, err)

	if c == nil {
		c = &redisutil.DefaultDialerConfig{}
	}

	c.Addr = &netutil.HostPort{
		Host: "localhost",
		Port: uint16(port64),
	}
	c.DBIndex = cmp.Or(c.DBIndex, DBIndex)

	d, err = redisutil.NewDefaultDialer(c)
	require.NoError(tb, err)

	return d
}

// FlushOnCleanup registers a cleanup function that uses the provided dialer
// to flush the database after a test.  The timeout used is [Timeout].
//
// TODO(a.garipov):  Consider ways of making Redis tests independent and
// parallel e.g. by using test names as key prefixes.
func FlushOnCleanup(tb testing.TB, d redisutil.Dialer) {
	tb.Helper()

	testutil.CleanupAndRequireSuccess(tb, func() (err error) {
		ctx := testutil.ContextWithTimeout(tb, Timeout)
		c, err := d.DialContext(ctx)
		if err != nil {
			return fmt.Errorf("dialing: %w", err)
		}

		defer func() { err = errors.WithDeferred(err, c.Close()) }()

		resp, err := redis.String(c.Do(redisutil.CmdFLUSHDB, redisutil.ParamSYNC))
		if err != nil {
			return fmt.Errorf("flushing: %w", err)
		}

		return validate.Equal("flush response", resp, redisutil.RespOK)
	})
}

// NewPool returns a *redisutil.DefaultPool for tests or skips the test if
// [EnvPort] is not set.  If c is nil, a default config is used.  If c is not
// nil:
//   - c.Dialer is ignored and set to a result of calling [NewDialer].
//   - c.Wait is ignored and set to true.
//
// NOTE:  It also uses [FlushOnCleanup] to flush the database after a test.
func NewPool(tb testing.TB, c *redisutil.DefaultPoolConfig) (p *redisutil.DefaultPool) {
	tb.Helper()

	// Create a dialer early to also skip early if it is not necessary.
	dialer := NewDialer(tb, nil)
	FlushOnCleanup(tb, dialer)

	if c == nil {
		c = &redisutil.DefaultPoolConfig{}
	}

	c.Logger = cmp.Or(c.Logger, slogutil.NewDiscardLogger())
	c.Dialer = dialer
	c.MaxConnLifetime = cmp.Or(c.MaxConnLifetime, MaxConnLifetime)
	c.IdleTimeout = cmp.Or(c.IdleTimeout, IdleTimeout)
	c.MaxActive = cmp.Or(c.MaxActive, MaxActive)
	c.MaxIdle = cmp.Or(c.MaxIdle, MaxIdle)
	c.Wait = true

	p, err := redisutil.NewDefaultPool(c)
	require.NoError(tb, err)

	return p
}
