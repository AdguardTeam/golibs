package netutil_test

import (
	"net/netip"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrefix_encoding(t *testing.T) {
	t.Parallel()

	p := &netutil.Prefix{
		Prefix: netip.MustParsePrefix("1.2.3.4/16"),
	}

	testutil.AssertMarshalText(t, "1.2.3.4/16", p)
	testutil.AssertUnmarshalText(t, "1.2.3.4/16", p)
}

func TestPrefix_UnmarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want       netip.Prefix
		wantErrMsg string
		name       string
		in         string
	}{{
		want:       netip.PrefixFrom(testIPv4Addr, testIPv4Addr.BitLen()),
		wantErrMsg: "",
		name:       "success_ipv4",
		in:         "1.2.3.4",
	}, {
		want:       netip.PrefixFrom(testIPv6Addr, testIPv6Addr.BitLen()),
		wantErrMsg: "",
		name:       "success_ipv6",
		in:         "1234::cdef",
	}, {
		want:       netip.PrefixFrom(testIPv6Addr, 16),
		wantErrMsg: "",
		name:       "success_ipv6",
		in:         "1234::cdef/16",
	}, {
		want:       netip.Prefix{},
		wantErrMsg: `ParseAddr("1.2.3.4.5"): IPv4 address too long`,
		name:       "bad_ipv4",
		in:         "1.2.3.4.5",
	}, {
		want: netip.Prefix{},
		wantErrMsg: `ParseAddr("1234:::cdef"): ` +
			`each colon-separated field must have at least one digit (at ":cdef")`,
		name: "bad_ipv6",
		in:   "1234:::cdef",
	}, {
		want: netip.Prefix{},
		wantErrMsg: `netip.ParsePrefix("1.2.3.4//16"): ParseAddr("1.2.3.4/"): ` +
			`unexpected character (at "/")`,
		name: "bad_cidr",
		in:   "1.2.3.4//16",
	}, {
		want:       netip.PrefixFrom(netip.MustParseAddr("::ffff:1.2.3.4"), 96),
		wantErrMsg: "",
		name:       "success_4_to_6",
		in:         "::ffff:1.2.3.4/96",
	}, {
		want:       netip.PrefixFrom(netip.MustParseAddr("::ffff:1.2.3.4"), 16),
		wantErrMsg: "",
		name:       "success_not_4_to_6",
		in:         "::ffff:1.2.3.4/16",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := &netutil.Prefix{}
			err := p.UnmarshalText([]byte(tc.in))

			assert.Equal(t, tc.want, p.Prefix)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
		})
	}
}

func BenchmarkPrefix_UnmarshalText(b *testing.B) {
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

	p := &netutil.Prefix{}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for b.Loop() {
				err = p.UnmarshalText([]byte(bc.in))
			}

			require.NotNil(b, p)
			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkPrefix_UnmarshalText
	//	BenchmarkPrefix_UnmarshalText/good_cidr4
	//	BenchmarkPrefix_UnmarshalText/good_cidr4-16         	 6845535	       154.8 ns/op	      16 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_ip4
	//	BenchmarkPrefix_UnmarshalText/good_ip4-16           	11513816	       112.1 ns/op	       8 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_cidr6
	//	BenchmarkPrefix_UnmarshalText/good_cidr6-16         	 7655109	       187.7 ns/op	      16 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_ip6
	//	BenchmarkPrefix_UnmarshalText/good_ip6-16           	 8027028	       166.3 ns/op	      16 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_cidr4to6
	//	BenchmarkPrefix_UnmarshalText/good_cidr4to6-16      	 5304142	       213.8 ns/op	      24 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_ip4to6
	//	BenchmarkPrefix_UnmarshalText/good_ip4to6-16        	 6912882	       156.0 ns/op	      16 B/op	       1 allocs/op
	//	BenchmarkPrefix_UnmarshalText/good_cidr_not4to6
	//	BenchmarkPrefix_UnmarshalText/good_cidr_not4to6-16  	 8418352	       208.8 ns/op	      24 B/op	       1 allocs/op
}

func BenchmarkPrefix_UnmarshalText_errors(b *testing.B) {
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

	p := &netutil.Prefix{}

	for _, bc := range benchErrCases {
		b.Run(bc.name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for b.Loop() {
				err = p.UnmarshalText([]byte(bc.in))
			}

			require.Error(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkPrefix_UnmarshalText_errors
	//	BenchmarkPrefix_UnmarshalText_errors/bad_cidr
	//	BenchmarkPrefix_UnmarshalText_errors/bad_cidr-16         	 1468476	       804.8 ns/op	     192 B/op	       7 allocs/op
	//	BenchmarkPrefix_UnmarshalText_errors/bad_ip
	//	BenchmarkPrefix_UnmarshalText_errors/bad_ip-16           	 5339067	       199.9 ns/op	      64 B/op	       2 allocs/op
}

func TestIsValidIPPrefixString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want assert.BoolAssertionFunc
		name string
		in   string
	}{{
		want: assert.True,
		name: "good_ipv4",
		in:   testIPv4Prefix.String(),
	}, {
		want: assert.True,
		name: "good_ipv6",
		in:   testIPv6Prefix.String(),
	}, {
		want: assert.True,
		name: "good_ipv4_zero",
		in:   testIPv4.String() + "/0",
	}, {
		want: assert.False,
		name: "bad_ip",
		in:   "1.2.3/8",
	}, {
		want: assert.False,
		name: "bad_empty",
		in:   testIPv4.String() + "",
	}, {
		want: assert.False,
		name: "bad_slash",
		in:   testIPv4.String() + "/",
	}, {
		want: assert.False,
		name: "bad_long",
		in:   testIPv4.String() + "/1111",
	}, {
		want: assert.False,
		name: "bad_invalid",
		in:   testIPv4.String() + "/!",
	}, {
		want: assert.False,
		name: "bad_ipv4_bits",
		in:   testIPv4.String() + "/33",
	}, {
		want: assert.False,
		name: "bad_ipv6_bits",
		in:   testIPv6.String() + "/129",
	}, {
		want: assert.False,
		name: "bad_bits_leading_zeroes",
		in:   testIPv6.String() + "/012",
	}, {
		want: assert.False,
		name: "bad_ipv6_zone",
		in:   testIPv6.String() + "%eth0/12",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.want(t, netutil.IsValidIPPrefixString(tc.in))
		})
	}
}

func BenchmarkIsValidIPPrefixString(b *testing.B) {
	benchCases := []struct {
		want require.BoolAssertionFunc
		name string
		in   string
	}{{
		want: require.True,
		name: "good_ipv4",
		in:   testIPv4Prefix.String(),
	}, {
		want: require.True,
		name: "good_ipv6",
		in:   testIPv6Prefix.String(),
	}, {
		want: require.False,
		name: "bad_ip",
		in:   "1.2.3/32",
	}, {
		want: require.False,
		name: "bad_empty",
		in:   testIPv4.String() + "",
	}, {
		want: require.False,
		name: "bad_slash",
		in:   testIPv4.String() + "/",
	}, {
		want: require.False,
		name: "bad_long",
		in:   testIPv4.String() + "/1111",
	}, {
		want: require.False,
		name: "bad_invalid",
		in:   testIPv4.String() + "/!",
	}, {
		want: require.False,
		name: "bad_overflow",
		in:   testIPv4.String() + "/129",
	}, {
		want: require.False,
		name: "bad_ipv6_zone",
		in:   testIPv6.String() + "%eth0/12",
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var got bool
			b.ReportAllocs()
			for b.Loop() {
				got = netutil.IsValidIPPrefixString(bc.in)
			}

			bc.want(b, got)
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: Apple M1 Pro
	//	BenchmarkIsValidIPPrefixString/good_ipv4-8         	25224192	        41.01 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/good_ipv6-8         	34011074	        35.29 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_ip-8            	43467367	        26.75 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_empty-8         	250258578	         4.809 ns/op       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_slash-8         	30908748	        38.80 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_long-8          	28030369	        42.61 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_invalid-8       	30037076	        39.82 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_overflow-8      	27973848	        41.42 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidIPPrefixString/bad_ipv6_zone-8     	39114808	        30.37 ns/op	       0 B/op	       0 allocs/op
}

func FuzzIsValidIPPrefixString(f *testing.F) {
	for _, seed := range []string{
		"",
		" ",
		"192.0.2.1",
		"192.0.2.1/0",
		"2001:db8::68",
		"1.2.3.4/",
		"1.2.3.4/1",
		"1.2.3.4/12",
		"1.2.3.4/128",
		"1.2.3.4/!",
		"1.2.3.4/012",
		"2001:db8::68/32",
		"2001:db8::68/256",
		"2001:db8::68/1024",
		"2001:db8::68%eth0/32",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ok := netutil.IsValidIPPrefixString(input)
		_, err := netip.ParsePrefix(input)

		require.Equalf(t, err == nil, ok, "input: %q, error: %v", input, err)
	})
}
