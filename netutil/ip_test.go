package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneIPs(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []net.IP(nil), netutil.CloneIPs(nil))
	assert.Equal(t, []net.IP{}, netutil.CloneIPs([]net.IP{}))

	ips := []net.IP{testIPv4}
	clone := netutil.CloneIPs(ips)
	assert.Equal(t, ips, clone)

	require.Len(t, clone, len(ips))
	require.Len(t, clone[0], len(ips[0]))

	assert.NotSame(t, &ips[0], &clone[0])
	assert.NotSame(t, &ips[0][0], &clone[0][0])
}

func TestSpecialAddrs(t *testing.T) {
	t.Parallel()

	assert.NotSame(t, netutil.IPv4bcast(), netutil.IPv4bcast())
	assert.NotSame(t, netutil.IPv4allsys(), netutil.IPv4allsys())
	assert.NotSame(t, netutil.IPv4allrouter(), netutil.IPv4allrouter())

	assert.NotSame(t, netutil.IPv4Zero(), netutil.IPv4Zero())
	assert.NotSame(t, netutil.IPv6Zero(), netutil.IPv6Zero())
}

func TestIPAndPortFromAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		in       net.Addr
		wantIP   net.IP
		wantPort uint16
	}{{
		name:     "nil",
		in:       nil,
		wantIP:   nil,
		wantPort: 0,
	}, {
		name:     "tcp",
		in:       &net.TCPAddr{IP: testIPv4, Port: 12345},
		wantIP:   testIPv4,
		wantPort: 12345,
	}, {
		name:     "udp",
		in:       &net.UDPAddr{IP: testIPv4, Port: 12345},
		wantIP:   testIPv4,
		wantPort: 12345,
	}, {
		name:     "custom",
		in:       struct{ net.Addr }{},
		wantIP:   nil,
		wantPort: 0,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotIP, gotPort := netutil.IPAndPortFromAddr(tc.in)
			assert.Equal(t, tc.wantIP, gotIP)
			assert.Equal(t, tc.wantPort, gotPort)
		})
	}
}

func TestValidateIP(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		wantErrMsg string
		wantErrAs  any
		in         net.IP
	}{{
		name:       "success_ipv4",
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         testIPv4,
	}, {
		name:       "success_ipv6",
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         testIPv6,
	}, {
		name:       "error_nil",
		wantErrMsg: `bad ip address "<nil>": ip address is empty`,
		wantErrAs:  new(*netutil.LengthError),
		in:         nil,
	}, {
		name:       "error_empty",
		wantErrMsg: `bad ip address "<nil>": ip address is empty`,
		wantErrAs:  new(*netutil.LengthError),
		in:         net.IP{},
	}, {
		name: "error_bad",
		wantErrMsg: `bad ip address "?010203": ` +
			`bad ip address length 3, allowed: [4 16]`,
		wantErrAs: new(*netutil.LengthError),
		in:        net.IP{1, 2, 3},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateIP(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}
