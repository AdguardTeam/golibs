package netutil_test

import (
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/netutil"
)

// check is a helper function for examples that panics on error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExampleIPFromReversedAddr() {
	ip, err := netutil.IPFromReversedAddr("4.3.2.1.in-addr.arpa")
	check(err)

	fmt.Println(ip)

	a := `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa`
	ip, err = netutil.IPFromReversedAddr(a)
	check(err)

	fmt.Println(ip)

	// Output:
	//
	// 1.2.3.4
	// ::abcd:1234
}

func ExampleIPToReversedAddr() {
	arpa, err := netutil.IPToReversedAddr(net.IP{1, 2, 3, 4})
	check(err)

	fmt.Println(arpa)

	ip := net.ParseIP("::abcd:1234")
	arpa, err = netutil.IPToReversedAddr(ip)
	check(err)

	fmt.Println(arpa)

	// Output:
	//
	// 4.3.2.1.in-addr.arpa
	// 4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa
}

func ExamplePrefixFromReversedAddr() {
	pref, err := netutil.PrefixFromReversedAddr("10.in-addr.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.PrefixFromReversedAddr("1.0.0.0.0.0.0.0.0.0.0.0.0.4.3.2.1.ip6.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.PrefixFromReversedAddr("in-addr.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.PrefixFromReversedAddr("ip6.arpa")
	check(err)

	fmt.Println(pref)

	// Output:
	//
	// 10.0.0.0/8
	// 1234::1000:0:0:0/68
	// 0.0.0.0/0
	// ::/0
}

func ExampleExtractReversedAddr() {
	pref, err := netutil.ExtractReversedAddr("_some-srv.10.in-addr.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.ExtractReversedAddr("_some-srv.1.0.0.0.0.0.0.0.0.0.0.0.0.4.3.2.1.ip6.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.ExtractReversedAddr("_some-srv.in-addr.arpa")
	check(err)

	fmt.Println(pref)

	pref, err = netutil.ExtractReversedAddr("_some-srv.ip6.arpa")
	check(err)

	fmt.Println(pref)

	// Output:
	//
	// 10.0.0.0/8
	// 1234::1000:0:0:0/68
	// 0.0.0.0/0
	// ::/0
}
