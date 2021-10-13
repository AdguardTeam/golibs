package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPPortFromAddr(t *testing.T) {
	ip4 := net.IP{1, 2, 3, 4}
	ipp := &netutil.IPPort{IP: ip4, Port: 12345}

	testCases := []struct {
		in         net.Addr
		wantIPPort *netutil.IPPort
		name       string
	}{{
		in:         nil,
		wantIPPort: nil,
		name:       "nil",
	}, {
		in:         &net.TCPAddr{IP: ip4, Port: 12345},
		wantIPPort: ipp,
		name:       "tcp",
	}, {
		in:         &net.UDPAddr{IP: ip4, Port: 12345},
		wantIPPort: ipp,
		name:       "udp",
	}, {
		in:         struct{ net.Addr }{},
		wantIPPort: nil,
		name:       "custom",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotIPPort := netutil.IPPortFromAddr(tc.in)
			assert.Equal(t, tc.wantIPPort, gotIPPort)
		})
	}
}

func TestCloneIPPort(t *testing.T) {
	assert.Equal(t, (*netutil.IPPort)(nil), (*netutil.IPPort)(nil).Clone())
	assert.Equal(t, &netutil.IPPort{}, (&netutil.IPPort{}).Clone())

	ipp := &netutil.IPPort{IP: net.IP{1, 2, 3, 4}, Port: 12345}
	clone := ipp.Clone()
	assert.Equal(t, ipp, clone)

	require.Len(t, clone.IP, len(ipp.IP))

	assert.NotSame(t, &ipp.IP[0], &clone.IP[0])
}

func TestCloneIPPorts(t *testing.T) {
	assert.Equal(t, []*netutil.IPPort(nil), netutil.CloneIPPorts(nil))
	assert.Equal(t, []*netutil.IPPort{}, netutil.CloneIPPorts([]*netutil.IPPort{}))

	ipps := []*netutil.IPPort{{IP: net.IP{1, 2, 3, 4}, Port: 12345}}
	clone := netutil.CloneIPPorts(ipps)
	assert.Equal(t, ipps, clone)

	require.Len(t, clone, len(ipps))
	require.Len(t, clone[0].IP, len(ipps[0].IP))

	assert.NotSame(t, &ipps[0], &clone[0])
	assert.NotSame(t, &ipps[0].IP[0], &clone[0].IP[0])
}

func TestIPPort_encoding(t *testing.T) {
	v := &netutil.IPPort{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}

	testutil.AssertMarshalText(t, "1.2.3.4:12345", v)
	testutil.AssertUnmarshalText(t, "1.2.3.4:12345", v)
}
