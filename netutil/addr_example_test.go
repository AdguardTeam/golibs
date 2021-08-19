package netutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleJoinHostPort() {
	fmt.Println(netutil.JoinHostPort("example.com", 12345))

	// Output:
	//
	// example.com:12345
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

func ExampleSplitHostPort() {
	host, port, err := netutil.SplitHostPort("example.com:12345")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%T(%[1]v)\n", host)
	fmt.Printf("%T(%[1]v)\n", port)

	// Output:
	//
	// string(example.com)
	// int(12345)
}

func ExampleSplitHost() {
	host, err := netutil.SplitHost("example.com:12345")
	if err != nil {
		panic(err)
	}

	fmt.Println(host)

	host, err = netutil.SplitHost("example.org")
	if err != nil {
		panic(err)
	}

	fmt.Println(host)

	_, err = netutil.SplitHost("[BAD:!")
	fmt.Println(err)

	// Output:
	//
	// example.com
	// example.org
	// address [BAD:!: missing ']' in address
}
