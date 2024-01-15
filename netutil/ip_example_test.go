package netutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleIPv4Zero() {
	fmt.Println(netutil.IPv4Zero())

	// Output:
	//
	// 0.0.0.0
}

func ExampleIPv6Zero() {
	fmt.Println(netutil.IPv6Zero())

	// Output:
	//
	// ::
}

func ExampleParseIP() {
	ip, err := netutil.ParseIP("1.2.3.4")
	fmt.Println(ip, err)

	ip, err = netutil.ParseIP("1234::cdef")
	fmt.Println(ip, err)

	ip, err = netutil.ParseIP("!!!")
	fmt.Println(ip, err)

	// Output:
	//
	// 1.2.3.4 <nil>
	// 1234::cdef <nil>
	// <nil> bad ip address "!!!"
}

func ExampleParseIPv4() {
	ip, err := netutil.ParseIPv4("1.2.3.4")
	fmt.Println(ip, err)

	ip, err = netutil.ParseIPv4("1234::cdef")
	fmt.Println(ip, err)

	ip, err = netutil.ParseIPv4("!!!")
	fmt.Println(ip, err)

	// Output:
	//
	// 1.2.3.4 <nil>
	// <nil> bad ipv4 address "1234::cdef"
	// <nil> bad ipv4 address "!!!"
}
