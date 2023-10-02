package sysresolv_test

import (
	"net/netip"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/sysresolv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemResolvers_Addrs(t *testing.T) {
	t.Parallel()

	r, err := sysresolv.NewSystemResolvers(nil, 53)
	require.NoError(t, err)

	var addrs []netip.AddrPort
	require.NotPanics(t, func() {
		addrs = r.Addrs()
	})

	assert.NotEmpty(t, addrs)
}
