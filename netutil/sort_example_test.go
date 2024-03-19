package netutil_test

import (
	"fmt"
	"net/netip"
	"slices"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExamplePreferIPv4() {
	addrs := []netip.Addr{
		testIPv6Addr.Next(),
		testIPv6Addr,
		{},
		testIPv4Addr.Next(),
		testIPv4Addr,
	}
	slices.SortFunc(addrs, netutil.PreferIPv4)

	fmt.Printf("%q\n", addrs)

	// Output:
	//
	// ["1.2.3.4" "1.2.3.5" "1234::cdef" "1234::cdf0" "invalid IP"]
}

func ExamplePreferIPv6() {
	addrs := []netip.Addr{
		testIPv4Addr.Next(),
		testIPv4Addr,
		{},
		testIPv6Addr.Next(),
		testIPv6Addr,
	}
	slices.SortFunc(addrs, netutil.PreferIPv6)

	fmt.Printf("%q\n", addrs)

	// Output:
	//
	// ["1234::cdef" "1234::cdf0" "1.2.3.4" "1.2.3.5" "invalid IP"]
}
