package netutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleAddrFamily_String() {
	fmt.Println(netutil.AddrFamilyIPv4, netutil.AddrFamilyIPv6)

	// An empty family.
	var fam netutil.AddrFamily
	fmt.Println(fam)

	// An unsupported family.
	fam = netutil.AddrFamily(1234)
	fmt.Println(fam)

	// Output:
	// ipv4 ipv6
	// none
	// !bad_addr_fam_1234
}

func ExampleAddrFamilyFromRRType() {
	// DNS type A.
	fmt.Println(netutil.AddrFamilyFromRRType(1))

	// DNS type AAAA.
	fmt.Println(netutil.AddrFamilyFromRRType(28))

	// Other DNS type.
	fmt.Println(netutil.AddrFamilyFromRRType(1234))

	// Output:
	// ipv4
	// ipv6
	// none
}
