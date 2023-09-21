package sysresolv_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/netutil/sysresolv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemResolvers_Addrs(t *testing.T) {
	r := sysresolv.NewTestResolvers(t, nil)

	var addrs []string
	require.NotPanics(t, func() {
		addrs = r.Addrs()
	})

	assert.NotEmpty(t, addrs)
}
