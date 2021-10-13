package netutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneHostPort(t *testing.T) {
	t.Parallel()

	assert.Equal(t, (*netutil.HostPort)(nil), (*netutil.HostPort)(nil).Clone())
	assert.Equal(t, &netutil.HostPort{}, (&netutil.HostPort{}).Clone())

	hp := &netutil.HostPort{Host: "example.com", Port: 12345}
	clone := hp.Clone()
	assert.Equal(t, hp, clone)
}

func TestCloneHostPorts(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []*netutil.HostPort(nil), netutil.CloneHostPorts(nil))
	assert.Equal(t, []*netutil.HostPort{}, netutil.CloneHostPorts([]*netutil.HostPort{}))

	hps := []*netutil.HostPort{{Host: "example.com", Port: 12345}}
	clone := netutil.CloneHostPorts(hps)
	assert.Equal(t, hps, clone)

	require.Len(t, clone, len(hps))
	require.Len(t, clone[0].Host, len(hps[0].Host))
}

func TestHostPort_encoding(t *testing.T) {
	t.Parallel()

	v := &netutil.HostPort{
		Host: "example.com",
		Port: 12345,
	}

	testutil.AssertMarshalText(t, "example.com:12345", v)
	testutil.AssertUnmarshalText(t, "example.com:12345", v)
}
