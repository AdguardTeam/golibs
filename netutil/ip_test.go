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

func TestSingleIPSubnet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want *net.IPNet
		name string
		in   net.IP
	}{{
		want: &net.IPNet{
			IP:   testIPv4,
			Mask: net.CIDRMask(32, 32),
		},
		name: "ipv4",
		in:   testIPv4,
	}, {
		want: &net.IPNet{
			IP:   testIPv6,
			Mask: net.CIDRMask(128, 128),
		},
		name: "ipv6",
		in:   testIPv6,
	}, {
		want: nil,
		name: "nil",
		in:   nil,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := netutil.SingleIPSubnet(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
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
		wantPort int
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

func TestCloneIPNet(t *testing.T) {
	t.Parallel()

	var (
		mask4 = net.CIDRMask(16, netutil.IPv4BitLen)
		mask6 = net.CIDRMask(32, netutil.IPv6BitLen)
	)

	testCases := []struct {
		n    *net.IPNet
		name string
	}{{
		n:    &net.IPNet{IP: testIPv4, Mask: mask4},
		name: "common_v4",
	}, {
		n:    &net.IPNet{IP: nil, Mask: mask4},
		name: "nil_ip_v4",
	}, {
		n:    &net.IPNet{IP: testIPv4, Mask: nil},
		name: "nil_mask_v4",
	}, {
		n:    &net.IPNet{IP: testIPv6, Mask: mask6},
		name: "common_v6",
	}, {
		n:    &net.IPNet{IP: nil, Mask: mask6},
		name: "nil_ip_v6",
	}, {
		n:    &net.IPNet{IP: testIPv6, Mask: nil},
		name: "nil_mask_v6",
	}, {
		n:    &net.IPNet{IP: nil, Mask: nil},
		name: "empty",
	}, {
		n:    nil,
		name: "nil",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clone := netutil.CloneIPNet(tc.n)
			assert.Equal(t, tc.n, clone)

			if tc.n == nil {
				return
			}

			assert.Len(t, clone.IP, len(tc.n.IP))
			assert.Len(t, clone.Mask, len(tc.n.Mask))

			assert.NotSame(t, tc.n, clone)
			if tc.n.IP != nil {
				assert.NotSame(t, &tc.n.IP[0], &clone.IP[0])
			}
			if tc.n.Mask != nil {
				assert.NotSame(t, &tc.n.Mask[0], &clone.Mask[0])
			}
		})
	}
}

func TestParseSubnet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want       *net.IPNet
		wantErrMsg string
		name       string
		in         string
	}{{
		want:       netutil.SingleIPSubnet(testIPv4),
		wantErrMsg: "",
		name:       "success_ipv4",
		in:         "1.2.3.4",
	}, {
		want:       netutil.SingleIPSubnet(testIPv6),
		wantErrMsg: "",
		name:       "success_ipv6",
		in:         "1234::cdef",
	}, {
		want:       &net.IPNet{IP: testIPv6, Mask: net.CIDRMask(16, netutil.IPv6BitLen)},
		wantErrMsg: "",
		name:       "success_ipv6",
		in:         "1234::cdef/16",
	}, {
		want:       nil,
		wantErrMsg: `bad cidr address "1.2.3.4.5": bad ip address "1.2.3.4.5"`,
		name:       "bad_ipv4",
		in:         "1.2.3.4.5",
	}, {
		want:       nil,
		wantErrMsg: `bad cidr address "1234:::cdef": bad ip address "1234:::cdef"`,
		name:       "bad_ipv6",
		in:         "1234:::cdef",
	}, {
		want:       nil,
		wantErrMsg: `bad cidr address "1.2.3.4//16"`,
		name:       "bad_cidr",
		in:         "1.2.3.4//16",
	}, {
		want:       &net.IPNet{IP: testIPv4, Mask: net.CIDRMask(0, netutil.IPv4BitLen)},
		wantErrMsg: "",
		name:       "success_4_to_6",
		in:         "::ffff:1.2.3.4/96",
	}, {
		want:       &net.IPNet{IP: testIPv4, Mask: net.CIDRMask(16, netutil.IPv6BitLen)},
		wantErrMsg: "",
		name:       "success_not_4_to_6",
		in:         "::ffff:1.2.3.4/16",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			n, err := netutil.ParseSubnet(tc.in)
			assert.Equal(t, tc.want, n)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if err != nil {
				assert.ErrorAs(t, err, new(*netutil.AddrError))
			}
		})
	}
}

func TestParseSubnet_equivalence(t *testing.T) {
	t.Parallel()

	const netStr = "::ffff:1.2.3.4/16"

	parsedIP, netN, err := net.ParseCIDR(netStr)
	require.NoError(t, err)

	netutilN, err := netutil.ParseSubnet(netStr)
	require.NoError(t, err)

	testCases := []struct {
		want assert.ComparisonAssertionFunc
		name string
		ip   net.IP
	}{{
		want: assert.NotEqual,
		name: "ip4",
		ip:   testIPv4,
	}, {
		want: assert.Equal,
		name: "ip6",
		ip:   testIPv6,
	}, {
		want: assert.NotEqual,
		name: "ip4_zero",
		ip:   netutil.IPv4Zero(),
	}, {
		want: assert.NotEqual,
		name: "ip6_zero",
		ip:   netutil.IPv6Zero(),
	}, {
		want: assert.NotEqual,
		name: "ip_from_str",
		ip:   parsedIP,
	}, {
		want: assert.Equal,
		name: "invalid",
		ip:   net.IP{1, 2, 3, 4, 5},
	}, {
		want: assert.Equal,
		name: "nil",
		ip:   nil,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.want(t, netN.Contains(tc.ip), netutilN.Contains(tc.ip))
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

func BenchmarkParseSubnet(b *testing.B) {
	benchCases := []struct {
		name string
		in   string
	}{{
		name: "good_cidr4",
		in:   "1.2.3.4/16",
	}, {
		name: "good_ip4",
		in:   "1.2.3.4",
	}, {
		name: "good_cidr6",
		in:   "abcd::1234/96",
	}, {
		name: "good_ip6",
		in:   "abcd::1234",
	}, {
		name: "good_cidr4to6",
		in:   "::ffff:1.2.3.4/97",
	}, {
		name: "good_ip4to6",
		in:   "::ffff:1.2.3.4",
	}, {
		name: "good_cidr_not4to6",
		in:   "::ffff:1.2.3.4/16",
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ipNetSink, errSink = netutil.ParseSubnet(bc.in)
			}

			require.NotNil(b, ipNetSink)
			require.NoError(b, errSink)
		})
	}

	benchErrCases := []struct {
		name string
		in   string
	}{{
		name: "bad_cidr",
		in:   "1.2.3.4//567",
	}, {
		name: "bad_ip",
		in:   "1.2.3.4.5",
	}}

	for _, bc := range benchErrCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, errSink = netutil.ParseSubnet(bc.in)
			}

			require.Error(b, errSink)
		})
	}
}
