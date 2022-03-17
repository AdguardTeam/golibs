package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubnetSet_optimized(t *testing.T) {
	t.Parallel()

	spPurpSet := netutil.SubnetSetFunc(netutil.IsSpecialPurpose)
	locSrvSet := netutil.SubnetSetFunc(netutil.IsLocallyServed)

	testCases := []struct {
		wantSpecialPurpose assert.BoolAssertionFunc
		wantLocallyServed  assert.BoolAssertionFunc
		name               string
		ip                 net.IP
	}{{
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "public",
		ip:                 net.IP{8, 8, 8, 8},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "unspecified_v4",
		ip:                 netutil.IPv4Zero(),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 net.IP{10, 0, 0, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "shared_address_space",
		ip:                 net.IP{100, 64, 0, 1},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "loopback",
		ip:                 net.IP{127, 0, 0, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "link_local",
		ip:                 net.IP{169, 254, 0, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 net.IP{172, 16, 0, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-1)",
		ip:                 net.IP{192, 0, 2, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "reserved",
		ip:                 net.IP{192, 88, 99, 1},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "private-use",
		ip:                 net.IP{192, 168, 0, 0},
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "non-local_v4",
		ip:                 net.IP{192, 169, 0, 1},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-2)",
		ip:                 net.IP{198, 51, 100, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation_(test-net-3)",
		ip:                 net.IP{203, 0, 113, 0},
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "limited_broadcast",
		ip:                 net.IP{255, 255, 255, 255},
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "public_v6",
		ip:                 net.ParseIP("::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "loopback_address",
		ip:                 net.ParseIP("::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "unspecified_v6",
		ip:                 net.ParseIP("::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		// "translat" stands for "translation" as per RFC 6890.
		name: "v4-v6_translat",
		ip:   net.ParseIP("64:ff9b::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "v4-v6_translat",
		ip:                 net.ParseIP("64:ff9b:1::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "discard-only",
		ip:                 net.ParseIP("100::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "ietf_protocol_assignments",
		ip:                 net.ParseIP("2001::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "teredo",
		ip:                 net.ParseIP("2001::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "port_control_protocol_anycast",
		ip:                 net.ParseIP("2001:1::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "traversal_using_relays_around_nat_anycast",
		ip:                 net.ParseIP("2001:1::2"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "benchmarking",
		ip:                 net.ParseIP("2001:2::10"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "amt",
		ip:                 net.ParseIP("2001:3::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "as112-v6",
		ip:                 net.ParseIP("2001:4:112::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "orchid",
		ip:                 net.ParseIP("2001:10::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "orchid_v2",
		ip:                 net.ParseIP("2001:20::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "documentation",
		ip:                 net.ParseIP("2001:db8::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "6to4",
		ip:                 net.ParseIP("2002::1"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "direct_delegation_as112_service",
		ip:                 net.ParseIP("2620:4f:8000::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "link-local",
		ip:                 net.ParseIP("fd00::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.False,
		name:               "unique-local",
		ip:                 net.ParseIP("fc00::"),
	}, {
		wantSpecialPurpose: assert.True,
		wantLocallyServed:  assert.True,
		name:               "linked-scoped_unicast",
		ip:                 net.ParseIP("fe80::12"),
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "invalid",
		ip:                 net.IP{1, 2, 3, 4, 5},
	}, {
		wantSpecialPurpose: assert.False,
		wantLocallyServed:  assert.False,
		name:               "nil",
		ip:                 nil,
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
	optimized := netutil.SubnetSetFunc(netutil.IsSpecialPurpose)
	nets, err := netutil.ParseSubnets(
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
	)
	require.NoError(b, err)
	general := netutil.SliceSubnetSet(nets)

	ips := []net.IP{
		net.ParseIP("8.8.8.8"),
		net.ParseIP("0.0.0.0"),
		net.ParseIP("10.0.0.0"),
		net.ParseIP("127.0.0.0"),
		net.ParseIP("169.254.0.0"),
		net.ParseIP("172.16.0.0"),
		net.ParseIP("192.0.2.0"),
		net.ParseIP("192.88.99.1"),
		net.ParseIP("192.168.0.0"),
		net.ParseIP("192.169.0.1"),
		net.ParseIP("198.51.100.0"),
		net.ParseIP("203.0.113.0"),
		net.ParseIP("224.0.0.1"),
		net.ParseIP("255.255.255.255"),
		net.ParseIP("::2"),
		net.ParseIP("::1"),
		net.ParseIP("::"),
		net.ParseIP("64:ff9b::1"),
		net.ParseIP("64:ff9b:1::1"),
		net.ParseIP("100::1"),
		net.ParseIP("2001::1"),
		net.ParseIP("2001::2"),
		net.ParseIP("2001:1::1"),
		net.ParseIP("2001:1::2"),
		net.ParseIP("2001:2::10"),
		net.ParseIP("2001:3::1"),
		net.ParseIP("2001:4:112::1"),
		net.ParseIP("2001:10::1"),
		net.ParseIP("2001:20::1"),
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2002::1"),
		net.ParseIP("2620:4f:8000::"),
		net.ParseIP("fd00::"),
		net.ParseIP("fe80::12"),
	}
	ipsLen := len(ips)

	b.Run("general_set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			boolSink = general.Contains(ips[i%ipsLen])
		}
	})

	b.Run("optimized_set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			boolSink = optimized.Contains(ips[i%ipsLen])
		}
	})
}
