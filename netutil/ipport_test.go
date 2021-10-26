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
	t.Parallel()

	ipp := &netutil.IPPort{IP: testIPv4, Port: 12345}

	testCases := []struct {
		in         net.Addr
		wantIPPort *netutil.IPPort
		name       string
	}{{
		in:         nil,
		wantIPPort: nil,
		name:       "nil",
	}, {
		in:         &net.TCPAddr{IP: testIPv4, Port: 12345},
		wantIPPort: ipp,
		name:       "tcp",
	}, {
		in:         &net.UDPAddr{IP: testIPv4, Port: 12345},
		wantIPPort: ipp,
		name:       "udp",
	}, {
		in:         struct{ net.Addr }{},
		wantIPPort: nil,
		name:       "custom",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotIPPort := netutil.IPPortFromAddr(tc.in)
			assert.Equal(t, tc.wantIPPort, gotIPPort)
		})
	}
}

func TestCloneIPPort(t *testing.T) {
	t.Parallel()

	assert.Equal(t, (*netutil.IPPort)(nil), (*netutil.IPPort)(nil).Clone())
	assert.Equal(t, &netutil.IPPort{}, (&netutil.IPPort{}).Clone())

	ipp := &netutil.IPPort{IP: testIPv4, Port: 12345}
	clone := ipp.Clone()
	assert.Equal(t, ipp, clone)

	require.Len(t, clone.IP, len(ipp.IP))

	assert.NotSame(t, &ipp.IP[0], &clone.IP[0])
}

func TestCloneIPPorts(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []*netutil.IPPort(nil), netutil.CloneIPPorts(nil))
	assert.Equal(t, []*netutil.IPPort{}, netutil.CloneIPPorts([]*netutil.IPPort{}))

	ipps := []*netutil.IPPort{{IP: testIPv4, Port: 12345}}
	clone := netutil.CloneIPPorts(ipps)
	assert.Equal(t, ipps, clone)

	require.Len(t, clone, len(ipps))
	require.Len(t, clone[0].IP, len(ipps[0].IP))

	assert.NotSame(t, &ipps[0], &clone[0])
	assert.NotSame(t, &ipps[0].IP[0], &clone[0].IP[0])
}

func TestIPPort_encoding(t *testing.T) {
	t.Parallel()

	v := &netutil.IPPort{
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 12345,
	}

	testutil.AssertMarshalText(t, "1.2.3.4:12345", v)
	testutil.AssertUnmarshalText(t, "1.2.3.4:12345", v)
}
