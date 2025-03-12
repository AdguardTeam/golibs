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
		want:       netip.Prefix{},
		wantErrMsg: `ParseAddr("1234:::cdef"): each colon-separated field must have at least one digit (at ":cdef")`,
		name:       "bad_ipv6",
		in:         "1234:::cdef",
	}, {
		want:       netip.Prefix{},
		wantErrMsg: `netip.ParsePrefix("1.2.3.4//16"): ParseAddr("1.2.3.4/"): unexpected character (at "/")`,
		name:       "bad_cidr",
		in:         "1.2.3.4//16",
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
		tc := tc
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
