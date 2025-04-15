package redisutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeredis"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPool_Get(t *testing.T) {
	conn := fakeredis.NewConn()

	var isClosed bool
	conn.OnDo = func(cmdName string, args ...any) (reply any, err error) {
		// NOTE:  [*redis.Pool] calls conn.Do("") before putting a conn back
		// into the pool.
		if cmdName == "" {
			require.Empty(t, args)
			require.False(t, isClosed)

			isClosed = true

			return redisutil.RespOK, nil
		}

		assert.Equal(t, redisutil.CmdSET, cmdName)

		require.Len(t, args, 2)

		assert.Equal(t, testKey, args[0])
		assert.Equal(t, testValue, args[1])

		return redisutil.RespOK, nil
	}

	conn.OnErr = func() (err error) { return nil }

	dialer := &fakeredis.Dialer{
		OnDialContext: func(ctx context.Context) (c redis.Conn, err error) {
			return conn, nil
		},
	}

	connTester := &fakeredis.ConnectionTester{
		OnTestConnection: func(ctx context.Context, c redis.Conn, _ time.Time) (err error) {
			assert.Equal(t, conn, c)

			return nil
		},
	}

	p, err := redisutil.NewDefaultPool(&redisutil.DefaultPoolConfig{
		Logger:           testLogger,
		Dialer:           dialer,
		ConnectionTester: connTester,
		MaxConnLifetime:  testMaxConnLifetime,
		IdleTimeout:      testIdleTimeout,
		MaxActive:        testMaxActive,
		MaxIdle:          textMaxIdle,
		Wait:             true,
	})
	require.NoError(t, err)

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	gotConn, err := p.Get(ctx)
	require.NoError(t, err)

	okStr, err := redis.String(gotConn.Do(redisutil.CmdSET, testKey, testValue))
	require.NoError(t, err)

	assert.Equal(t, "OK", okStr)

	err = gotConn.Close()
	require.NoError(t, err)

	assert.True(t, isClosed)
}

func TestDefaultPool_Get_integration(t *testing.T) {
	p := newIntegrationPool(t)

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	c, err := p.Get(ctx)
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, c.Close)

	okStr, err := redis.String(c.Do(redisutil.CmdSET, testKey, testValue, redisutil.ParamNX))
	require.NoError(t, err)

	assert.Equal(t, "OK", okStr)

	gotVal, err := redis.String(c.Do(redisutil.CmdGET, testKey))
	require.NoError(t, err)

	assert.Equal(t, testValue, gotVal)
}
