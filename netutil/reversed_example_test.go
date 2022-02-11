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

	a := `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa`
	ip, err = netutil.IPFromReversedAddr(a)
	if err != nil {
		panic(err)
	}

	fmt.Println(ip)

	// Output:
	//
	// 1.2.3.4
	// ::abcd:1234
}

func ExampleIPToReversedAddr() {
	arpa, err := netutil.IPToReversedAddr(net.IP{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}

	fmt.Println(arpa)

	ip := net.ParseIP("::abcd:1234")
	arpa, err = netutil.IPToReversedAddr(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println(arpa)

	// Output:
	//
	// 4.3.2.1.in-addr.arpa
	// 4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa
}

func ExampleSubnetFromReversedAddr() {
	subnet, err := netutil.SubnetFromReversedAddr("10.in-addr.arpa")
	if err != nil {
		panic(err)
	}

	fmt.Println(subnet)

	subnet, err = netutil.SubnetFromReversedAddr("0.10.in-addr.arpa")
	if err != nil {
		panic(err)
	}

	fmt.Println(subnet)

	subnet, err = netutil.SubnetFromReversedAddr("3.2.1.d.c.b.a.ip6.arpa")
	if err != nil {
		panic(err)
	}

	fmt.Println(subnet)

	// Output:
	//
	// 10.0.0.0/8
	// 10.0.0.0/16
	// abcd:1230::/28
}

func ExampleSubnetFromReversedAddr_domainOnly() {
	a := `in-addr.arpa`
	_, err := netutil.SubnetFromReversedAddr(a)

	fmt.Println(err)

	// Output:
	//
	// bad arpa domain name "in-addr.arpa": not a reversed ip network
}
