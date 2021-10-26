package netutil_test

import (
	"fmt"
	"net"

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

func ExampleParseSubnet() {
	ip := net.IP{1, 2, 3, 4}
	otherIP := net.IP{1, 2, 3, 5}

	n, err := netutil.ParseSubnet("1.2.3.4")
	fmt.Println(n, err)
	fmt.Printf("%s is in %s: %t\n", ip, n, n.Contains(ip))
	fmt.Printf("%s is in %s: %t\n", otherIP, n, n.Contains(otherIP))

	n, err = netutil.ParseSubnet("1.2.3.4/16")
	fmt.Println(n, err)
	fmt.Printf("%s is in %s: %t\n", ip, n, n.Contains(ip))
	fmt.Printf("%s is in %s: %t\n", otherIP, n, n.Contains(otherIP))

	// Output:
	//
	// 1.2.3.4/32 <nil>
	// 1.2.3.4 is in 1.2.3.4/32: true
	// 1.2.3.5 is in 1.2.3.4/32: false
	// 1.2.3.4/16 <nil>
	// 1.2.3.4 is in 1.2.3.4/16: true
	// 1.2.3.5 is in 1.2.3.4/16: true
}

func ExampleSingleIPSubnet() {
	ip4 := net.IP{1, 2, 3, 4}
	otherIP4 := net.IP{1, 2, 3, 5}

	n := netutil.SingleIPSubnet(ip4)
	fmt.Printf("%s is in %s: %t\n", ip4, n, n.Contains(ip4))
	fmt.Printf("%s is in %s: %t\n", otherIP4, n, n.Contains(otherIP4))

	ip6 := net.ParseIP("1234::cdef")
	otherIP6 := net.ParseIP("1234::cdff")

	n = netutil.SingleIPSubnet(ip6)
	fmt.Printf("%s is in %s: %t\n", ip6, n, n.Contains(ip6))
	fmt.Printf("%s is in %s: %t\n", otherIP6, n, n.Contains(otherIP6))

	// Output:
	//
	// 1.2.3.4 is in 1.2.3.4/32: true
	// 1.2.3.5 is in 1.2.3.4/32: false
	// 1234::cdef is in 1234::cdef/128: true
	// 1234::cdff is in 1234::cdef/128: false
}
