package netutil_test

import (
	"fmt"
	"net/netip"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleSubnetSet_Contains() {
	var s netutil.SubnetSet = netip.MustParsePrefix("1.2.3.4/16")

	fmt.Println("contains 1.2.3.4:", s.Contains(netip.MustParseAddr("1.2.3.4")))
	fmt.Println("contains 4.3.2.1:", s.Contains(netip.MustParseAddr("4.3.2.1")))

	s = netip.MustParsePrefix("0.0.0.0/0")
	fmt.Println("0.0.0.0/0 contains 1234::cdef:", s.Contains(netip.MustParseAddr("1234::cdef")))

	s = netip.MustParsePrefix("::/0")
	fmt.Println("::/0 contains 1.2.3.4:", s.Contains(netip.MustParseAddr("1.2.3.4")))

	// Output:
	//
	// contains 1.2.3.4: true
	// contains 4.3.2.1: false
	// 0.0.0.0/0 contains 1234::cdef: false
	// ::/0 contains 1.2.3.4: false
}

func ExampleSliceSubnetSet_Contains() {
	nets := []netip.Prefix{
		netip.MustParsePrefix("1.2.3.0/24"),
		netip.MustParsePrefix("ffff:12ab::/16"),
	}

	s := netutil.SliceSubnetSet(nets)

	fmt.Println("contains 1.2.3.4:      ", s.Contains(netip.MustParseAddr("1.2.3.4")))
	fmt.Println("contains 4.3.2.1:      ", s.Contains(netip.MustParseAddr("4.3.2.1")))
	fmt.Println("contains ffff:12ab::10:", s.Contains(netip.MustParseAddr("ffff:12ab::10")))
	fmt.Println("contains 12ab:ffff::10:", s.Contains(netip.MustParseAddr("12ab:ffff::10")))

	fmt.Println()

	s = netutil.SliceSubnetSet{}

	fmt.Println("contains 1.2.3.4:      ", s.Contains(netip.MustParseAddr("1.2.3.4")))
	fmt.Println("contains ffff:12ab::10:", s.Contains(netip.MustParseAddr("ffff:12ab::10")))

	fmt.Println()

	s = netutil.SliceSubnetSet{
		netip.MustParsePrefix("0.0.0.0/0"),
		netip.MustParsePrefix("::/0"),
	}

	fmt.Println("contains 1.2.3.4:      ", s.Contains(netip.MustParseAddr("1.2.3.4")))
	fmt.Println("contains ffff:12ab::10:", s.Contains(netip.MustParseAddr("ffff:12ab::10")))
	fmt.Println("contains zero value:   ", s.Contains(netip.Addr{}))

	// Output:
	//
	// contains 1.2.3.4:       true
	// contains 4.3.2.1:       false
	// contains ffff:12ab::10: true
	// contains 12ab:ffff::10: false
	//
	// contains 1.2.3.4:       false
	// contains ffff:12ab::10: false
	//
	// contains 1.2.3.4:       true
	// contains ffff:12ab::10: true
	// contains zero value:    false
}

func ExampleSubnetSetFunc() {
	var (
		ipv4      = netip.MustParseAddr("1.2.3.4")
		ipv6      = netip.MustParseAddr("1234::cdef")
		ipInvalid = netip.Addr{}
	)

	var s netutil.SubnetSet = netutil.SubnetSetFunc(func(ip netip.Addr) (ok bool) {
		return true
	})

	fmt.Printf("all contains %s: %t\n", ipv4, s.Contains(ipv4))
	fmt.Printf("all contains %s: %t\n", ipv6, s.Contains(ipv6))
	fmt.Printf("all contains %s: %t\n", ipInvalid, s.Contains(ipInvalid))

	s = netutil.SubnetSetFunc(func(ip netip.Addr) (ok bool) {
		return false
	})

	fmt.Printf("none contains %s: %t\n", ipv4, s.Contains(ipv4))
	fmt.Printf("none contains %s: %t\n", ipv6, s.Contains(ipv6))
	fmt.Printf("none contains %s: %t\n", ipInvalid, s.Contains(ipInvalid))

	s = netutil.SubnetSetFunc(netip.Addr.IsValid)

	fmt.Printf("valid contains %s: %t\n", ipv4, s.Contains(ipv4))
	fmt.Printf("valid contains %s: %t\n", ipv6, s.Contains(ipv6))
	fmt.Printf("valid contains %s: %t\n", ipInvalid, s.Contains(ipInvalid))

	// Output:
	//
	// all contains 1.2.3.4: true
	// all contains 1234::cdef: true
	// all contains invalid IP: true
	// none contains 1.2.3.4: false
	// none contains 1234::cdef: false
	// none contains invalid IP: false
	// valid contains 1.2.3.4: true
	// valid contains 1234::cdef: true
	// valid contains invalid IP: false
}

func ExampleSubnetSetFunc_Contains() {
	s := netutil.SubnetSetFunc(func(ip netip.Addr) (ok bool) {
		slice := ip.AsSlice()

		return len(slice) > 0 && slice[0] == 0xFF
	})

	fmt.Println("contains 255.0.0.1:", s.Contains(netip.MustParseAddr("255.0.0.1")))
	fmt.Println("contains 254.0.0.1:", s.Contains(netip.MustParseAddr("254.0.0.1")))
	fmt.Println("contains ff00:::1: ", s.Contains(netip.MustParseAddr("ff00::1")))
	fmt.Println("contains ff:::1:   ", s.Contains(netip.MustParseAddr("ff::1")))

	// Output:
	//
	// contains 255.0.0.1: true
	// contains 254.0.0.1: false
	// contains ff00:::1:  true
	// contains ff:::1:    false
}
