package netutil_test

import (
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleIPFromReversedAddr() {
	ip, err := netutil.IPFromReversedAddr("4.3.2.1.in-addr.arpa")
	if err != nil {
		panic(err)
	}

	fmt.Println(ip)

	// Output:
	//
	// 1.2.3.4
}

func ExampleIPFromReversedAddr_ipv6() {
	a := `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa`
	ip, err := netutil.IPFromReversedAddr(a)
	if err != nil {
		panic(err)
	}

	fmt.Println(ip)

	// Output:
	//
	// ::abcd:1234
}

func ExampleIPToReversedAddr() {
	arpa, err := netutil.IPToReversedAddr(net.IP{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}

	fmt.Println(arpa)

	// Output:
	//
	// 4.3.2.1.in-addr.arpa
}

func ExampleIPToReversedAddr_ipv6() {
	ip := net.ParseIP("::abcd:1234")
	arpa, err := netutil.IPToReversedAddr(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println(arpa)

	// Output:
	//
	// 4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa
}
