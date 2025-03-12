package netutil_test

import (
	"net/netip"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
)

func TestSubnetSet_optimized(t *testing.T) {
	t.Parallel()

	spPurpSet := netutil.SubnetSetFunc(netutil.IsSpecialPurpose)
	locSrvSet := netutil.SubnetSetFunc(netutil.IsLocallyServed)

	testCases := []struct {
		wantSpecialPurpose assert.BoolAssertionFunc
		wantLocallyServed  assert.BoolAssertionFunc
		ip                 netip.Addr
		name               string
	}{{
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "public",
		ip:                 netip.MustParseAddr("8.8.8.8"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "unspecified_v4",
		ip:                 netip.MustParseAddr("0.0.0.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 netip.MustParseAddr("10.0.0.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "shared_address_space",
		ip:                 netip.MustParseAddr("100.64.0.1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "loopback",
		ip:                 netip.MustParseAddr("127.0.0.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "link_local",
		ip:                 netip.MustParseAddr("169.254.0.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 netip.MustParseAddr("172.16.0.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-1)",
		ip:                 netip.MustParseAddr("192.0.2.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "reserved",
		ip:                 netip.MustParseAddr("192.88.99.1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 netip.MustParseAddr("192.168.0.0"),
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "non-local_v4",
		ip:                 netip.MustParseAddr("192.169.0.1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-2)",
		ip:                 netip.MustParseAddr("198.51.100.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-3)",
		ip:                 netip.MustParseAddr("203.0.113.0"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "limited_broadcast",
		ip:                 netip.MustParseAddr("255.255.255.255"),
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "public_v6",
		ip:                 netip.MustParseAddr("::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "loopback_address",
		ip:                 netip.MustParseAddr("::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "unspecified_v6",
		ip:                 netip.MustParseAddr("::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		// "translat" stands for "translation" as per RFC 6890.
		name: "v4-v6_translat",
		ip:   netip.MustParseAddr("64:ff9b::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "v4-v6_translat",
		ip:                 netip.MustParseAddr("64:ff9b:1::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "discard-only",
		ip:                 netip.MustParseAddr("100::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "ietf_protocol_assignments",
		ip:                 netip.MustParseAddr("2001::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "teredo",
		ip:                 netip.MustParseAddr("2001::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "port_control_protocol_anycast",
		ip:                 netip.MustParseAddr("2001:1::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "traversal_using_relays_around_nat_anycast",
		ip:                 netip.MustParseAddr("2001:1::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "benchmarking",
		ip:                 netip.MustParseAddr("2001:2::10"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "amt",
		ip:                 netip.MustParseAddr("2001:3::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "as112-v6",
		ip:                 netip.MustParseAddr("2001:4:112::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "orchid",
		ip:                 netip.MustParseAddr("2001:10::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "orchid_v2",
		ip:                 netip.MustParseAddr("2001:20::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation",
		ip:                 netip.MustParseAddr("2001:db8::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "6to4",
		ip:                 netip.MustParseAddr("2002::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "direct_delegation_as112_service",
		ip:                 netip.MustParseAddr("2620:4f:8000::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "link-local",
		ip:                 netip.MustParseAddr("fd00::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "unique-local",
		ip:                 netip.MustParseAddr("fc00::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "linked-scoped_unicast",
		ip:                 netip.MustParseAddr("fe80::12"),
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "zero",
		ip:                 netip.Addr{},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name+"_is_special-purpose", func(t *testing.T) {
			t.Parallel()

			tc.wantSpecialPurpose(t, spPurpSet.Contains(tc.ip))
		})
		t.Run(tc.name+"_is_locally-served", func(t *testing.T) {
			t.Parallel()

			tc.wantLocallyServed(t, locSrvSet.Contains(tc.ip))
		})
	}
}

func BenchmarkSliceSubnetSet_comparison(b *testing.B) {
	rawNets := []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.0.0/29",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"240.0.0.0/4",
		"255.255.255.255/32",
		"::1/128",
		"::/128",
		"64:ff9b::/96",
		"100::/64",
		"2001::/23",
		"2001::/32",
		"2001:2::/48",
		"2001:db8::/32",
		"2001:10::/28",
		"2002::/16",
		"fc00::/7",
		"fe80::/10",
	}
	nets := make([]netip.Prefix, 0, len(rawNets))
	for _, s := range rawNets {
		nets = append(nets, netip.MustParsePrefix(s))
	}

	rawIPs := []string{
		"8.8.8.8",
		"0.0.0.0",
		"10.0.0.0",
		"127.0.0.0",
		"169.254.0.0",
		"172.16.0.0",
		"192.0.2.0",
		"192.88.99.1",
		"192.168.0.0",
		"192.169.0.1",
		"198.51.100.0",
		"203.0.113.0",
		"224.0.0.1",
		"255.255.255.255",
		"::2",
		"::1",
		"::",
		"64:ff9b::1",
		"64:ff9b:1::1",
		"100::1",
		"2001::1",
		"2001::2",
		"2001:1::1",
		"2001:1::2",
		"2001:2::10",
		"2001:3::1",
		"2001:4:112::1",
		"2001:10::1",
		"2001:20::1",
		"2001:db8::1",
		"2002::1",
		"2620:4f:8000::",
		"fd00::",
		"fe80::12",
	}
	ipsLen := len(rawIPs)
	ips := make([]netip.Addr, 0, ipsLen)
	for _, s := range rawIPs {
		ips = append(ips, netip.MustParseAddr(s))
	}

	general := netutil.SliceSubnetSet(nets)
	optimized := netutil.SubnetSetFunc(netutil.IsSpecialPurpose)

	b.Run("general_set", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			boolSink = general.Contains(ips[i%ipsLen])
		}
	})

	b.Run("optimized_set", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := range b.N {
			boolSink = optimized.Contains(ips[i%ipsLen])
		}
	})

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkSliceSubnetSet_comparison
	//	BenchmarkSliceSubnetSet_comparison/general_set
	//	BenchmarkSliceSubnetSet_comparison/general_set-16         	13402231	        92.23 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSliceSubnetSet_comparison/optimized_set
	//	BenchmarkSliceSubnetSet_comparison/optimized_set-16       	119661958	         9.272 ns/op	       0 B/op	       0 allocs/op
}
