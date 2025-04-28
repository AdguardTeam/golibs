package netutil_test

import (
	"net"
	"net/netip"
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

const nonARPADomain = "valid.domain.example"

var (
	testIPv4Pref     = netip.PrefixFrom(testIPv4Addr, testIPv4Addr.BitLen())
	testIPv4PartPref = netip.MustParsePrefix("10.0.0.0/8")

	emptyV4Pref = netip.PrefixFrom(netip.IPv4Unspecified(), 0)

	testIPv6Pref     = netip.PrefixFrom(testIPv6Addr, testIPv6Addr.BitLen())
	testIPv6PartPref = netip.MustParsePrefix("1234::/64")

	emptyV6Pref = netip.PrefixFrom(netip.IPv6Unspecified(), 0)
)

func TestIPFromReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantErrAs  any
		want       netip.Addr
		wantErrMsg string
	}{{
		name:       "good_ipv4",
		in:         ipv4RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4Addr,
	}, {
		name:       "good_ipv4_fqdn",
		in:         ipv4RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4Addr,
	}, {
		name:       "good_ipv4_case",
		in:         ipv4RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv4Addr,
	}, {
		name: "bad_ipv4_missing",
		in:   ipv4Missing,
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv4_char",
		in:   ipv4Char,
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`ParseAddr("1.0.z.127"): unexpected character (at "z.127")`,
		wantErrAs: new(*netutil.AddrError),
		want:      netip.Addr{},
	}, {
		name: "ipv6_arpa_v4",
		in:   testIPv6Addr.String() + ipv4Suffix,
		wantErrMsg: `bad arpa domain name "` + testIPv6Addr.String() +
			ipv4Suffix + `": bad ipv4 address "` + testIPv6Addr.String() + `"`,
		wantErrAs: new(*netutil.AddrError),
		want:      netip.Addr{},
	}, {
		name:       "good_ipv6",
		in:         ipv6RevGood,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6Addr,
	}, {
		name:       "good_ipv6_fqdn",
		in:         ipv6RevGood + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6Addr,
	}, {
		name:       "good_ipv6_case",
		in:         ipv6RevGoodUp,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       testIPv6Addr,
	}, {
		name: "bad_ipv6_many",
		in:   ipv6RevMany,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMany + `": ` +
			`not a full reversed ip address`,
		wantErrAs: new(*netutil.AddrError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_missing",
		in:   ipv6RevMissing,
		wantErrMsg: `bad arpa domain name "` + ipv6RevMissing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_char_lo",
		in:   ipv6RevCharLo,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_char_hi",
		in:   ipv6RevCharHi,
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.RuneError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_dots",
		in:   ipv6RevDots,
		wantErrMsg: `bad arpa domain name "` + ipv6RevDots + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.LengthError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_len",
		in:   ipv6RevLen,
		wantErrMsg: `bad arpa domain name "` + ipv6RevLen + `": ` +
			`bad arpa domain name length 70, allowed: 72`,
		wantErrAs: new(*netutil.LengthError),
		want:      netip.Addr{},
	}, {
		name: "bad_ipv6_space",
		in:   ipv6RevSpace,
		wantErrMsg: `bad arpa domain name "` + ipv6RevSpace + `": ` +
			`bad arpa domain name rune ' '`,
		wantErrAs: new(*netutil.RuneError),
		want:      netip.Addr{},
	}, {
		name: "ipv4_arpa_v6",
		in:   testIPv4Addr.String() + ipv6Suffix,
		wantErrMsg: `bad arpa domain name "` + testIPv4Addr.String() +
			ipv6Suffix + `": bad arpa domain name length 16, allowed: 72`,
		wantErrAs: new(*netutil.AddrError),
		want:      netip.Addr{},
	}, {
		name: "not_a_reversed_ip",
		in:   testIPv4.String(),
		wantErrMsg: `bad arpa domain name "` + testIPv4.String() + `": ` +
			`bad top-level domain name label "4": all octets are numeric`,
		wantErrAs: new(errors.Error),
		want:      netip.Addr{},
	}, {
		name: "not_arpa",
		in:   nonARPADomain,
		wantErrMsg: `bad arpa domain name "` + nonARPADomain + `": ` +
			`not a full reversed ip address`,
		wantErrAs: new(*netutil.AddrError),
		want:      netip.Addr{},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			addr, err := netutil.IPFromReversedAddr(tc.in)
			assert.Equal(t, tc.want, addr)
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

func TestPrefixFromReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want       netip.Prefix
		wantErrAs  any
		wantErrMsg string
		in         string
		name       string
	}{{
		want:       netip.MustParsePrefix("1.2.3.4/32"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGood,
		name:       "good_ipv4_single_addr",
	}, {
		want:       netip.MustParsePrefix("1.2.3.4/32"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGood + ".",
		name:       "good_ipv4_single_addr_fqdn",
	}, {
		want:       netip.MustParsePrefix("1.2.3.4/32"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4RevGoodUp,
		name:       "good_ipv4_single_addr_case",
	}, {
		want:       netip.MustParsePrefix("10.0.0.0/32"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         `0.0.0.` + ipv4NetRevGood,
		name:       "good_ipv4_single_addr_leading_zero",
	}, {
		want:       netip.MustParsePrefix("10.0.0.0/16"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         "0." + ipv4NetRevGood,
		name:       "good_ipv4_subnet_leading_zero",
	}, {
		want:       testIPv4PartPref,
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv4NetRevGood,
		name:       "good_ipv4_subnet",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "4.3.2.1.not-in-addr.arpa": ` +
			`not a reversed ip network`,
		in:   "4.3.2.1.not-in-addr.arpa",
		name: "almost_arpa_v4",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv4Missing,
		name: "bad_ipv4_missing",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`ParseAddr("1.0.z.127"): unexpected character (at "z.127")`,
		in:   ipv4Char,
		name: "bad_ipv4_char",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*strconv.NumError),
		wantErrMsg: `bad arpa domain name "x.` + ipv4NetRevGood + `": ` +
			`strconv.ParseUint: parsing "x": invalid syntax`,
		in:   `x.` + ipv4NetRevGood,
		name: "bad_ipv4_subnet_char",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "05.` + ipv4NetRevGood + `": ` +
			`bad domain name label "05": leading zero is forbidden at this position`,
		in:   `05.` + ipv4NetRevGood,
		name: "bad_ipv4_subnet_unexpected_zero",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "5.` + ipv4RevGood + `": ` +
			`not a reversed ip network`,
		in:   `5.` + ipv4RevGood,
		name: "bad_ipv4_too_long",
	}, {
		want:       netip.MustParsePrefix("1234::cdef/128"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGood,
		name:       "good_ipv6_single_addr",
	}, {
		want:       netip.MustParsePrefix("1234::cdef/128"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGood + ".",
		name:       "good_ipv6_single_addr_fqdn",
	}, {
		want:       netip.MustParsePrefix("1234::cdef/128"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevGoodUp,
		name:       "good_ipv6_single_addr_case",
	}, {
		want:       netip.MustParsePrefix("1234::0000/128"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         ipv6RevZeroes + "." + ipv6RevGoodSuffix,
		name:       "good_ipv6_single_addr_leading_zeroes",
	}, {
		want:       netip.MustParsePrefix("1234::0000/68"),
		wantErrAs:  nil,
		wantErrMsg: "",
		in:         "0." + ipv6RevGoodSuffix,
		name:       "good_ipv6_subnet_leading_zeroes",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "a.b.c.not-ip6.arpa": ` +
			`not a reversed ip network`,
		in:   "a.b.c.not-ip6.arpa",
		name: "almost_arpa_v6",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevMany + `": ` +
			`not a full reversed ip address`,
		in:   ipv6RevMany,
		name: "bad_ipv6_single_addr_many",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + strings.TrimPrefix(ipv6RevMany, "4.3.2.1.") + `": ` +
			`not a reversed ip network`,
		in:   strings.TrimPrefix(ipv6RevMany, "4.3.2.1."),
		name: "bad_ipv6_many",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevMissing + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv6RevMissing,
		name: "bad_ipv6_missing",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6RevCharLo,
		name: "bad_ipv6_single_addr_char_lo",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevCharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6RevCharHi,
		name: "bad_ipv6_single_addr_char_hi",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6NetRevChar + `": ` +
			`bad arpa domain name rune 'z'`,
		in:   ipv6NetRevChar,
		name: "bad_ipv6_char",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevDots + `": ` +
			`bad domain name label "": domain name label is empty`,
		in:   ipv6RevDots,
		name: "bad_ipv6_dots",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevTooLong + `": ` +
			`arpa domain name is too long: got 74, max 72`,
		in:   ipv6RevTooLong,
		name: "bad_ipv6_len",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad arpa domain name "` + ipv6RevSpace + `": ` +
			`bad arpa domain name rune ' '`,
		in:   ipv6RevSpace,
		name: "bad_ipv6_space",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + ipv6NetRevHex + `": ` +
			`not a reversed ip network`,
		in:   ipv6NetRevHex,
		name: "bad_ipv6_hex",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad arpa domain name "` + testIPv4.String() + `": ` +
			`bad top-level domain name label "4": all octets are numeric`,
		in:   testIPv4.String(),
		name: "not_a_reversed_subnet",
	}, {
		want:       emptyV4Pref,
		wantErrAs:  nil,
		wantErrMsg: ``,
		in:         ipv4Suffix[1:],
		name:       "root_arpa_v4",
	}, {
		want:       emptyV6Pref,
		wantErrAs:  nil,
		wantErrMsg: ``,
		in:         ipv6Suffix[1:],
		name:       "root_arpa_v6",
	}, {
		want:      netip.Prefix{},
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad arpa domain name "` + nonARPADomain + `": ` +
			`not a reversed ip network`,
		in:   nonARPADomain,
		name: "not_arpa",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			subnet, err := netutil.PrefixFromReversedAddr(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
			assert.Equal(t, tc.want, subnet)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestExtractReversedAddr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want    netip.Prefix
		name    string
		domain  string
		wantErr string
	}{{
		want:   netip.Prefix{},
		name:   "not_an_arpa",
		domain: "some.domain.name.",
		wantErr: `bad arpa domain name "some.domain.name": ` +
			`not a reversed ip network`,
	}, {
		want:   netip.Prefix{},
		name:   "bad_domain_name",
		domain: "abc.123.",
		wantErr: `bad arpa domain name "abc.123": ` +
			`bad top-level domain name label "123": all octets are numeric`,
	}, {
		want:   netip.Prefix{},
		name:   "almost_arpa_v4",
		domain: "4.3.2.1.not-in-addr.arpa",
		wantErr: `bad arpa domain name "4.3.2.1.not-in-addr.arpa": ` +
			`not a reversed ip network`,
	}, {
		want:    testIPv4Pref,
		name:    "whole_v4",
		domain:  ipv4RevGood,
		wantErr: "",
	}, {
		want:    testIPv4PartPref,
		name:    "partial_v4",
		domain:  ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    netip.PrefixFrom(testIPv4PartPref.Addr(), testIPv4PartPref.Bits()+8),
		name:    "partial_v4_zero_label",
		domain:  "0." + ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    testIPv4Pref,
		name:    "whole_v4_within_domain",
		domain:  "a." + ipv4RevGood,
		wantErr: "",
	}, {
		want:    testIPv4Pref,
		name:    "whole_v4_additional_label",
		domain:  "5." + ipv4RevGood,
		wantErr: "",
	}, {
		want:    testIPv4PartPref,
		name:    "partial_v4_within_domain",
		domain:  "abc." + ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    testIPv4PartPref,
		name:    "overflow_v4",
		domain:  "256." + ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    testIPv4PartPref,
		name:    "overflow_v4_within_domain",
		domain:  "a.256." + ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    testIPv4PartPref,
		name:    "partial_v4_leading_zero_label",
		domain:  "05." + ipv4NetRevGood,
		wantErr: "",
	}, {
		want:    emptyV4Pref,
		name:    "empty_v4",
		domain:  ipv4Suffix[1:],
		wantErr: ``,
	}, {
		want:    emptyV4Pref,
		name:    "empty_v4_within_domain",
		domain:  "a" + ipv4Suffix,
		wantErr: ``,
	}, {
		want:   netip.Prefix{},
		name:   "almost_arpa_v6",
		domain: "a.b.c.not-ip6.arpa",
		wantErr: `bad arpa domain name "a.b.c.not-ip6.arpa": ` +
			`not a reversed ip network`,
	}, {
		want:    testIPv6Pref,
		name:    "whole_v6",
		domain:  ipv6RevGood,
		wantErr: "",
	}, {
		want:    testIPv6PartPref,
		name:    "partial_v6",
		domain:  ipv6RevGoodSuffix,
		wantErr: "",
	}, {
		want:    testIPv6Pref,
		name:    "whole_v6_within_domain",
		domain:  "g." + ipv6RevGood,
		wantErr: "",
	}, {
		want:    testIPv6Pref,
		name:    "whole_v6_additional_label",
		domain:  "1." + ipv6RevGood,
		wantErr: "",
	}, {
		want:    testIPv6PartPref,
		name:    "partial_v6_within_domain",
		domain:  "l." + ipv6RevGoodSuffix,
		wantErr: "",
	}, {
		want:    emptyV6Pref,
		name:    "empty_v6",
		domain:  ipv6Suffix[1:],
		wantErr: ``,
	}, {
		want:    emptyV6Pref,
		name:    "empty_v6_within_domain",
		domain:  "g" + ipv6Suffix,
		wantErr: ``,
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			subnet, err := netutil.ExtractReversedAddr(tc.domain)
			if tc.wantErr != "" {
				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.EqualError(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, subnet)
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
			var (
				prefix netip.Prefix
				err    error
			)

			b.ReportAllocs()
			for b.Loop() {
				prefix, err = netutil.PrefixFromReversedAddr(bc.in)
			}

			require.NotNil(b, prefix)
			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkSubnetFromReversedAddr
	//	BenchmarkSubnetFromReversedAddr/ipv4_single_addr
	//	BenchmarkSubnetFromReversedAddr/ipv4_single_addr-16         	 3429879	       321.1 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSubnetFromReversedAddr/ipv4_subnet
	//	BenchmarkSubnetFromReversedAddr/ipv4_subnet-16              	 6412760	       186.3 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSubnetFromReversedAddr/ipv6_single_addr
	//	BenchmarkSubnetFromReversedAddr/ipv6_single_addr-16         	  790614	      1449 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSubnetFromReversedAddr/ipv6_subnet
	//	BenchmarkSubnetFromReversedAddr/ipv6_subnet-16              	 1381117	       817.3 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkExtractReversedAddr(b *testing.B) {
	const serviceLabel = "_srv."

	benchCases := []struct {
		name string
		in   string
	}{{
		name: "ipv4_root",
		in:   ipv4Suffix[len("."):],
	}, {
		name: "ipv4_subnet",
		in:   ipv4NetRevGood,
	}, {
		name: "ipv4_subnet_within_domain",
		in:   serviceLabel + ipv4NetRevGood,
	}, {
		name: "ipv6_root",
		in:   ipv6Suffix[len("."):],
	}, {
		name: "ipv6_subnet",
		in:   ipv6NetRevGood,
	}, {
		name: "ipv6_subnet_within_domain",
		in:   serviceLabel + ipv6NetRevGood,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var (
				prefix netip.Prefix
				err    error
			)

			b.ReportAllocs()
			for b.Loop() {
				prefix, err = netutil.ExtractReversedAddr(bc.in)
			}

			require.NotNil(b, prefix)
			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: Apple M1 Pro
	//	BenchmarkExtractReversedAddr/ipv4_root-8         	13283544	        89.83 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkExtractReversedAddr/ipv4_subnet-8       	 7445521	       160.6 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkExtractReversedAddr/ipv4_subnet_within_domain-8         	 5402799	       220.6 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkExtractReversedAddr/ipv6_root-8                         	13672136	        87.12 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkExtractReversedAddr/ipv6_subnet-8                       	 1000000	      1019 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkExtractReversedAddr/ipv6_subnet_within_domain-8         	 1000000	      1038 ns/op	       0 B/op	       0 allocs/op
}
