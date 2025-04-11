package redisutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPool_Get_integration(t *testing.T) {
	p := newPool(t)

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	c, err := p.Get(ctx)
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, c.Close)

	const (
		key   = "key"
		value = "value"
	)

	okStr, err := redis.String(c.Do(redisutil.CmdSET, key, value, redisutil.ParamNX))
	require.NoError(t, err)

	assert.Equal(t, "OK", okStr)

	gotVal, err := redis.String(c.Do(redisutil.CmdGET, key))
	require.NoError(t, err)

	assert.Equal(t, value, gotVal)
}
