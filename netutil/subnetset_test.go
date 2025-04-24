package netutil_test

import (
	"encoding/binary"
	"net"
	"net/netip"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
)

var (
	// locallyServedPrefixes is a slice subnet set of locally served networks,
	// as per [RFC 6303].  It's used as a reference for testing the optimized
	// implementations of the [netutil.SubnetSet] interface.  It must not be
	// modified.
	//
	// [RFC 6303]: https://datatracker.ietf.org/doc/html/rfc6303
	locallyServedPrefixes = netutil.SliceSubnetSet{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("127.0.0.0/8"),
		netip.MustParsePrefix("169.254.0.0/16"),
		netip.MustParsePrefix("172.16.0.0/12"),
		netip.MustParsePrefix("192.0.2.0/24"),
		netip.MustParsePrefix("192.168.0.0/16"),
		netip.MustParsePrefix("198.51.100.0/24"),
		netip.MustParsePrefix("203.0.113.0/24"),
		netip.MustParsePrefix("255.255.255.255/32"),
		netip.MustParsePrefix("::/128"),
		netip.MustParsePrefix("::1/128"),
		netip.MustParsePrefix("2001:db8::/32"),
		netip.MustParsePrefix("fd00::/8"),
		netip.MustParsePrefix("fe80::/10"),
	}

	// specialPurposePrefixes is a slice subnet set of special purpose networks,
	// as per IANA Special-Purpose Address Registry.  It's used as a reference
	// for testing the optimized implementations of the [netutil.SubnetSet]
	// interface.  It must not be modified.
	//
	// See https://www.iana.org/assignments/iana-ipv4-special-registry and
	// https://www.iana.org/assignments/iana-ipv6-special-registry.
	specialPurposePrefixes = append(
		locallyServedPrefixes,
		netip.MustParsePrefix("0.0.0.0/8"),
		netip.MustParsePrefix("100.64.0.0/10"),
		netip.MustParsePrefix("169.254.0.0/16"),
		netip.MustParsePrefix("172.16.0.0/12"),
		netip.MustParsePrefix("192.0.0.0/24"),
		netip.MustParsePrefix("192.0.0.0/29"),
		netip.MustParsePrefix("192.88.99.0/24"),
		netip.MustParsePrefix("198.18.0.0/15"),
		netip.MustParsePrefix("240.0.0.0/4"),
		netip.MustParsePrefix("64:ff9b::/96"),
		netip.MustParsePrefix("64:ff9b:1::/48"),
		netip.MustParsePrefix("100::/64"),
		netip.MustParsePrefix("2001::/23"),
		netip.MustParsePrefix("2001::/32"),
		netip.MustParsePrefix("2001:1::1/128"),
		netip.MustParsePrefix("2001:1::2/128"),
		netip.MustParsePrefix("2001:2::/48"),
		netip.MustParsePrefix("2001:3::/32"),
		netip.MustParsePrefix("2001:4:112::/48"),
		netip.MustParsePrefix("2001:10::/28"),
		netip.MustParsePrefix("2001:20::/28"),
		netip.MustParsePrefix("2002::/16"),
		netip.MustParsePrefix("2620:4f:8000::/48"),
		netip.MustParsePrefix("fc00::/7"),
	)
)

var (
	// testIPv4s is a bunch of IPv4 addresses to test the [netutil.SubnetSet]
	// interface.
	testIPv4s = []netip.Addr{
		netip.MustParseAddr("8.8.8.8"),
		netip.MustParseAddr("0.0.0.0"),
		netip.MustParseAddr("10.0.0.0"),
		netip.MustParseAddr("127.0.0.0"),
		netip.MustParseAddr("169.254.0.0"),
		netip.MustParseAddr("172.16.0.0"),
		netip.MustParseAddr("192.0.2.0"),
		netip.MustParseAddr("192.88.99.1"),
		netip.MustParseAddr("192.168.0.0"),
		netip.MustParseAddr("192.169.0.1"),
		netip.MustParseAddr("198.51.100.0"),
		netip.MustParseAddr("203.0.113.0"),
		netip.MustParseAddr("224.0.0.1"),
		netip.MustParseAddr("255.255.255.255"),
	}

	// testIPv6s is a bunch of IPv6 addresses to test the [netutil.SubnetSet]
	// interface.
	testIPv6s = []netip.Addr{
		netip.MustParseAddr("::2"),
		netip.MustParseAddr("::1"),
		netip.MustParseAddr("::"),
		netip.MustParseAddr("64:ff9b::1"),
		netip.MustParseAddr("64:ff9b:1::1"),
		netip.MustParseAddr("100::1"),
		netip.MustParseAddr("2001::1"),
		netip.MustParseAddr("2001::2"),
		netip.MustParseAddr("2001:1::1"),
		netip.MustParseAddr("2001:1::2"),
		netip.MustParseAddr("2001:2::10"),
		netip.MustParseAddr("2001:3::1"),
		netip.MustParseAddr("2001:4:112::1"),
		netip.MustParseAddr("2001:10::1"),
		netip.MustParseAddr("2001:20::1"),
		netip.MustParseAddr("2001:db8::1"),
		netip.MustParseAddr("2002::1"),
		netip.MustParseAddr("2620:4f:8000::"),
		netip.MustParseAddr("fd00::"),
		netip.MustParseAddr("fe80::12"),
	}

	// testSubnetSetCases is a set of cases for fuzzing the [SubnetSet.Contains]
	// method.
	testSubnetSetCases = []struct {
		reference netutil.SubnetSet
		set       netutil.SubnetSet
		name      string
	}{{
		reference: locallyServedPrefixes,
		set:       netutil.SubnetSetFunc(netutil.IsLocallyServed),
		name:      "is_locally_served",
	}, {
		reference: specialPurposePrefixes,
		set:       netutil.SubnetSetFunc(netutil.IsSpecialPurpose),
		name:      "is_special_purpose",
	}}
)

// checkSubnetSet is a helper function that checks the result of calling
// [SubnetSet.Contains] method of the given tested set for the given IP address
// and compares it with the reference.
func checkSubnetSet(t *testing.T, ip netip.Addr, tested netutil.SubnetSet, want bool) {
	t.Helper()

	if want {
		assert.Truef(t, tested.Contains(ip), "%s", ip)
	} else {
		assert.Falsef(t, tested.Contains(ip), "%s", ip)
	}
}

func TestSubnetSet_optimized(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ip   netip.Addr
		name string
	}{{
		name: "public",
		ip:   netip.MustParseAddr("8.8.8.8"),
	}, {
		name: "unspecified_v4",
		ip:   netip.MustParseAddr("0.0.0.0"),
	}, {
		name: "private-use",
		ip:   netip.MustParseAddr("10.0.0.0"),
	}, {
		name: "shared_address_space",
		ip:   netip.MustParseAddr("100.64.0.1"),
	}, {
		name: "loopback",
		ip:   netip.MustParseAddr("127.0.0.0"),
	}, {
		name: "link_local",
		ip:   netip.MustParseAddr("169.254.0.0"),
	}, {
		name: "private-use",
		ip:   netip.MustParseAddr("172.16.0.0"),
	}, {
		name: "documentation_(test-net-1)",
		ip:   netip.MustParseAddr("192.0.2.0"),
	}, {
		name: "reserved",
		ip:   netip.MustParseAddr("192.88.99.1"),
	}, {
		name: "private-use",
		ip:   netip.MustParseAddr("192.168.0.0"),
	}, {
		name: "non-local_v4",
		ip:   netip.MustParseAddr("192.169.0.1"),
	}, {
		name: "documentation_(test-net-2)",
		ip:   netip.MustParseAddr("198.51.100.0"),
	}, {
		name: "documentation_(test-net-3)",
		ip:   netip.MustParseAddr("203.0.113.0"),
	}, {
		name: "limited_broadcast",
		ip:   netip.MustParseAddr("255.255.255.255"),
	}, {
		name: "public_v6",
		ip:   netip.MustParseAddr("::2"),
	}, {
		name: "loopback_address",
		ip:   netip.MustParseAddr("::1"),
	}, {
		name: "unspecified_v6",
		ip:   netip.MustParseAddr("::"),
	}, {
		// "translat" stands for "translation" as per RFC 6890.
		name: "v4-v6_translat",
		ip:   netip.MustParseAddr("64:ff9b::1"),
	}, {
		name: "v4-v6_translat",
		ip:   netip.MustParseAddr("64:ff9b:1::1"),
	}, {
		name: "discard-only",
		ip:   netip.MustParseAddr("100::1"),
	}, {
		name: "ietf_protocol_assignments",
		ip:   netip.MustParseAddr("2001::1"),
	}, {
		name: "teredo",
		ip:   netip.MustParseAddr("2001::2"),
	}, {
		name: "port_control_protocol_anycast",
		ip:   netip.MustParseAddr("2001:1::1"),
	}, {
		name: "traversal_using_relays_around_nat_anycast",
		ip:   netip.MustParseAddr("2001:1::2"),
	}, {
		name: "benchmarking",
		ip:   netip.MustParseAddr("2001:2::10"),
	}, {
		name: "amt",
		ip:   netip.MustParseAddr("2001:3::1"),
	}, {
		name: "as112-v6",
		ip:   netip.MustParseAddr("2001:4:112::1"),
	}, {
		name: "orchid",
		ip:   netip.MustParseAddr("2001:10::1"),
	}, {
		name: "orchid_v2",
		ip:   netip.MustParseAddr("2001:20::1"),
	}, {
		name: "documentation",
		ip:   netip.MustParseAddr("2001:db8::1"),
	}, {
		name: "6to4",
		ip:   netip.MustParseAddr("2002::1"),
	}, {
		name: "direct_delegation_as112_service",
		ip:   netip.MustParseAddr("2620:4f:8000::"),
	}, {
		name: "link-local",
		ip:   netip.MustParseAddr("fd00::"),
	}, {
		name: "unique-local",
		ip:   netip.MustParseAddr("fc00::"),
	}, {
		name: "linked-scoped_unicast",
		ip:   netip.MustParseAddr("fe80::12"),
	}, {
		name: "zero",
		ip:   netip.Addr{},
	}}

	for _, sc := range testSubnetSetCases {
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()

			for _, tc := range testCases {
				want := sc.reference.Contains(tc.ip)

				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					checkSubnetSet(t, tc.ip, sc.set, want)
				})
			}
		})
	}
}

func BenchmarkSliceSubnetSet_comparison(b *testing.B) {
	optimized := netutil.SubnetSetFunc(netutil.IsSpecialPurpose)

	benchCases := []struct {
		set  netutil.SubnetSet
		name string
	}{{
		set:  specialPurposePrefixes,
		name: "slice_set",
	}, {
		set:  optimized,
		name: "func_set",
	}}

	versionCases := []struct {
		name string
		ips  []netip.Addr
	}{{
		name: "v4",
		ips:  testIPv4s,
	}, {
		name: "v6",
		ips:  testIPv6s,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			for _, vc := range versionCases {
				b.Run(vc.name, func(b *testing.B) {
					b.ReportAllocs()

					for i, l := 0, len(vc.ips); b.Loop(); i++ {
						_ = bc.set.Contains(vc.ips[i%l])
					}
				})
			}
		})
	}

	// Most recent results:
	//  goos: darwin
	//  goarch: amd64
	//  pkg: github.com/AdguardTeam/golibs/netutil
	//  cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	//  BenchmarkSliceSubnetSet_comparison/slice_set/v4-12      	 6617144	       169.8 ns/op	       0 B/op	       0 allocs/op
	//  BenchmarkSliceSubnetSet_comparison/slice_set/v6-12      	 6276302	       190.3 ns/op	       0 B/op	       0 allocs/op
	//  BenchmarkSliceSubnetSet_comparison/func_set/v4-12       	72217934	        16.45 ns/op	       0 B/op	       0 allocs/op
	//  BenchmarkSliceSubnetSet_comparison/func_set/v6-12       	64475127	        18.99 ns/op	       0 B/op	       0 allocs/op
}

func FuzzSubnetSet_Contains_v4(f *testing.F) {
	for _, ip := range testIPv4s {
		ipData := ip.As4()
		f.Add(binary.BigEndian.Uint32(ipData[:]))
	}

	f.Fuzz(func(t *testing.T, ipUint32 uint32) {
		var ipData [net.IPv4len]byte
		binary.BigEndian.PutUint32(ipData[:], ipUint32)
		ip := netip.AddrFrom4(ipData)

		for _, tc := range testSubnetSetCases {
			want := tc.reference.Contains(ip)

			t.Run(tc.name, func(t *testing.T) {
				checkSubnetSet(t, ip, tc.set, want)
			})
		}
	})
}

func FuzzSubnetSet_Contains_v6(f *testing.F) {
	for _, ip := range testIPv6s {
		ipData := ip.As16()
		f.Add(
			binary.BigEndian.Uint64(ipData[:8]),
			binary.BigEndian.Uint64(ipData[8:]),
		)
	}

	f.Fuzz(func(t *testing.T, ipHi, ipLo uint64) {
		var ipData [net.IPv6len]byte
		binary.BigEndian.PutUint64(ipData[:8], ipHi)
		binary.BigEndian.PutUint64(ipData[8:], ipLo)
		ip := netip.AddrFrom16(ipData)

		for _, tc := range testSubnetSetCases {
			want := tc.reference.Contains(ip)

			t.Run(tc.name, func(t *testing.T) {
				checkSubnetSet(t, ip, tc.set, want)
			})
		}
	})
}
