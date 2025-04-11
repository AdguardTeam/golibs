package redisutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultDialer(t *testing.T) {
	d, err := redisutil.NewDefaultDialer(&redisutil.DefaultDialerConfig{
		Addr: &netutil.HostPort{
			Host: "localhost",
			Port: redisutil.DefaultPort,
		},
	})
	require.NoError(t, err)

	assert.NotNil(t, d)
}
