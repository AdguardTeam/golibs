package netutil_test

import (
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ipv4Suffix     = `.in-addr.arpa`
	ipv4RevGood    = `4.3.2.1` + ipv4Suffix
	ipv4RevGoodUp  = `4.3.2.1.In-Addr.Arpa`
	ipv4NetRevGood = `10` + ipv4Suffix

	ipv4RevGoodUnspecified = `0.0.0.0` + ipv4Suffix

	ipv4Missing = `.0.0.127` + ipv4Suffix
	ipv4Char    = `1.0.z.127` + ipv4Suffix
)

const (
	ipv6Suffix    = `.ip6.arpa`
	ipv6RevZeroes = `0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0`

	ipv6RevGoodSuffix = `0.0.0.0.0.0.0.0.0.0.0.0.4.3.2.1.ip6.arpa`
	ipv6RevGood       = `f.e.d.c.0.0.0.0.0.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevGoodUp     = `F.E.D.C.0.0.0.0.0.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6NetRevGood    = `1.` + ipv6RevGoodSuffix

	ipv6RevGoodUnspecified = ipv6RevZeroes + "." + ipv6RevZeroes + ipv6Suffix

	ipv6RevCharHi  = `4.3.2.1.d.c.b.a.0.z.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevCharLo  = `4.3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevDots    = `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6RevZeroes + "." + ipv6Suffix
	ipv6RevLen     = `3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevTooLong = `5.4.3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevMany    = `4.3.2.1.dbc.b.a.0.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevMissing = `.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6RevSpace   = `4.3.2.1.d.c.b.a. .0.0.0.0.0.0.0.` + ipv6RevGoodSuffix
	ipv6NetRevHex  = `10.` + ipv6RevGoodSuffix
	ipv6NetRevChar = `z.` + ipv6RevGoodSuffix
)

func TestIPFromReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantErrMsg string
		wantErrAs  any
		want       net.IP
	}{{
		name:       "good_ipv4",
		in:         ipv4RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4,
	}, {
		name:       "good_ipv4_fqdn",
		in:         ipv4RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4,
	}, {
		name:       "good_ipv4_case",
		in:         ipv4RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4,
	}, {
		name: "bad_ipv4_missing",
		in:   ipv4Missing,
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      nil,
	}, {
		name: "bad_ipv4_char",
		in:   ipv4Char,
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`bad ipv4 address "1.0.z.127"`,
		wantErrAs: new(*netutil.AddrError),
		want:      nil,
	}, {
		name:       "good_ipv6",
		in:         ipv6RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6,
	}, {
		name:       "good_ipv6_fqdn",
		in:         ipv6RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6,
	}, {
		name:       "good_ipv6_case",
		in:         ipv6RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6,
	}, {
		name: "bad_ipv6_many",
		in:   ipv6RevMany,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMany + `": ` +
			`not a full reversed ip address`,
		wantErrAs: new(*netutil.AddrError),
		want:      nil,
	}, {
		name: "bad_ipv6_missing",
		in:   ipv6RevMissing,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMissing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      nil,
	}, {
		name: "bad_ipv6_char_lo",
		in:   ipv6RevCharLo,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_char_hi",
		in:   ipv6RevCharHi,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_dots",
		in:   ipv6RevDots,
		wantErrMsg: `bad arpa domain name "` + ipv6RevDots + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      nil,
	}, {
		name: "bad_ipv6_len",
		in:   ipv6RevLen,
		wantErrMsg: `bad arpa domain name "` + ipv6RevLen + `": ` +
			`bad arpa domain name length 70, allowed: 72`,
		wantErrAs: new(*netutil.LengthError),
		want:      nil,
	}, {
		name: "bad_ipv6_space",
		in:   ipv6RevSpace,
		wantErrMsg: `bad arpa domain name "` + ipv6RevSpace + `": ` +
			`bad arpa domain name rune ' '`,
		wantErrAs: new(*netutil.RuneError),
		want:      nil,
	}, {
		name: "not_a_reversed_ip",
		in:   testIPv4.String(),
		wantErrMsg: `bad arpa domain name "` + testIPv4.String() + `": ` +
			`bad top-level domain name label "4": all octets are numeric`,
		wantErrAs: new(errors.Error),
		want:      nil,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ip, err := netutil.IPFromReversedAddr(tc.in)
			assert.Equal(t, tc.want.To16(), ip.To16())
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestIPToReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		want       string
		wantErrMsg string
		wantErrAs  any
		in         net.IP
	}{{
		name:       "good_ipv4",
		want:       ipv4RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         testIPv4,
	}, {
		name:       "good_ipv6",
		want:       ipv6RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         testIPv6,
	}, {
		name:       "nil_ip",
		want:       "",
		wantErrMsg: `bad ip address "<nil>"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         nil,
	}, {
		name:       "empty_ip",
		want:       "",
		wantErrMsg: `bad ip address "<nil>"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         net.IP{},
	}, {
		name:       "unspecified_ipv4",
		want:       ipv4RevGoodUnspecified,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.IPv4zero,
	}, {
		name:       "unspecified_ipv6",
		want:       ipv6RevGoodUnspecified,
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.IPv6unspecified,
	}, {
		name:       "wrong_length_ip",
		want:       "",
		wantErrMsg: `bad ip address "?0102030405"`,
		wantErrAs:  new(*netutil.AddrError),
		in:         net.IP{1, 2, 3, 4, 5},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			arpa, err := netutil.IPToReversedAddr(tc.in)
			assert.Equal(t, tc.want, arpa)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

// newIPNet returns an IP network to use in test cases.  It doesn't validate
// anything.
func newIPNet(ip net.IP, ones int) (n *net.IPNet) {
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(ones, len(ip)*8),
	}
}

func TestSubnetFromReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want       *net.IPNet
		wantErrAs  any
		wantErrMsg string
		in         string
		name       string
	}{{
		want:       newIPNet(testIPv4, netutil.IPv4BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGood,
		name:       "good_ipv4_single_addr",
	}, {
		want:       newIPNet(testIPv4, netutil.IPv4BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGood + ".",
		name:       "good_ipv4_single_addr_fqdn",
	}, {
		want:       newIPNet(testIPv4, netutil.IPv4BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGoodUp,
		name:       "good_ipv4_single_addr_case",
	}, {
		want:       newIPNet(testIPv4ZeroTail, netutil.IPv4BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         `0.0.0.` + ipv4NetRevGood,
		name:       "good_ipv4_single_addr_leading_zero",
	}, {
		want:       newIPNet(testIPv4ZeroTail, 16),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         "0." + ipv4NetRevGood,
		name:       "good_ipv4_subnet_leading_zero",
	}, {
		want:       newIPNet(testIPv4ZeroTail, 8),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4NetRevGood,
		name:       "good_ipv4_subnet",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv4Missing,
		name: "bad_ipv4_missing",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`bad ipv4 address "1.0.z.127"`,
		in:   ipv4Char,
		name: "bad_ipv4_char",
	}, {
		want:      nil,
		wantErrAs: new(*strconv.NumError),
		wantErrMsg: `bad arpa domain name "x.` + ipv4NetRevGood + `": ` +
			`strconv.ParseUint: parsing "x": invalid syntax`,
		in:   `x.` + ipv4NetRevGood,
		name: "bad_ipv4_subnet_char",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "05.` + ipv4NetRevGood + `": ` +
			`bad domain name label "05": leading zero is forbidden at this position`,
		in:   `05.` + ipv4NetRevGood,
		name: "bad_ipv4_subnet_unexpected_zero",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "5.` + ipv4RevGood + `": ` +
			`not a reversed ip network`,
		in:   `5.` + ipv4RevGood,
		name: "bad_ipv4_too_long",
	}, {
		want:       newIPNet(testIPv6, netutil.IPv6BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGood,
		name:       "good_ipv6_single_addr",
	}, {
		want:       newIPNet(testIPv6, netutil.IPv6BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGood + ".",
		name:       "good_ipv6_single_addr_fqdn",
	}, {
		want:       newIPNet(testIPv6, netutil.IPv6BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGoodUp,
		name:       "good_ipv6_single_addr_case",
	}, {
		want:       newIPNet(testIPv6ZeroTail, netutil.IPv6BitLen),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevZeroes + "." + ipv6RevGoodSuffix,
		name:       "good_ipv6_single_addr_leading_zeroes",
	}, {
		want:       newIPNet(testIPv6ZeroTail, 68),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         "0." + ipv6RevGoodSuffix,
		name:       "good_ipv6_subnet_leading_zeroes",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevMany + `": ` +
			`not a full reversed ip address`,
		in:   ipv6RevMany,
		name: "bad_ipv6_single_addr_many",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + strings.TrimPrefix(ipv6RevMany, "4.3.2.1.") + `": ` +
			`not a reversed ip network`,
		in:   strings.TrimPrefix(ipv6RevMany, "4.3.2.1."),
		name: "bad_ipv6_many",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevMissing + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv6RevMissing,
		name: "bad_ipv6_missing",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6RevCharLo,
		name: "bad_ipv6_single_addr_char_lo",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6RevCharHi,
		name: "bad_ipv6_single_addr_char_hi",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6NetRevChar + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6NetRevChar,
		name: "bad_ipv6_char",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevDots + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv6RevDots,
		name: "bad_ipv6_dots",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevTooLong + `": ` +
			`arpa domain name is too long: got 74, max 72`,
		in:   ipv6RevTooLong,
		name: "bad_ipv6_len",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevSpace + `": ` +
			`bad arpa domain name rune ' '`,
		in:   ipv6RevSpace,
		name: "bad_ipv6_space",
	}, {
		want:      nil,
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv6NetRevHex + `": ` +
			`not a reversed ip network`,
		in:   ipv6NetRevHex,
		name: "bad_ipv6_hex",
	}, {
		want:      nil,
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad arpa domain name "` + testIPv4.String() + `": ` +
			`bad top-level domain name label "4": all octets are numeric`,
		in:   testIPv4.String(),
		name: "not_a_reversed_subnet",
	}, {
		want:      nil,
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad arpa domain name "` + ipv4Suffix[1:] + `": ` +
			`not a reversed ip network`,
		in:   ipv4Suffix[1:],
		name: "root_arpa_v4",
	}, {
		want:      nil,
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad arpa domain name "` + ipv6Suffix[1:] + `": ` +
			`not a reversed ip network`,
		in:   ipv6Suffix[1:],
		name: "root_arpa_v6",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			subnet, err := netutil.SubnetFromReversedAddr(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			} else {
				require.NotNil(t, subnet)

				assert.Equal(t, tc.want.IP.To16(), subnet.IP.To16())
				assert.Equal(t, tc.want.Mask, subnet.Mask)
			}
		})
	}
}

func BenchmarkSubnetFromReversedAddr(b *testing.B) {
	benchCases := []struct {
		name string
		in   string
	}{{
		name: "ipv4_single_addr",
		in:   ipv4RevGood,
	}, {
		name: "ipv4_subnet",
		in:   ipv4NetRevGood,
	}, {
		name: "ipv6_single_addr",
		in:   ipv6RevGood,
	}, {
		name: "ipv6_subnet",
		in:   ipv6NetRevGood,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ipNetSink, errSink = netutil.SubnetFromReversedAddr(bc.in)
			}

			require.NotNil(b, ipNetSink)
			require.NoError(b, errSink)
		})
	}
}
