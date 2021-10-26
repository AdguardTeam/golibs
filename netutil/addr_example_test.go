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

func ExampleSubdomains() {
	subdomains := netutil.Subdomains("subsub.sub.domain.tld")
	for _, sub := range subdomains {
		fmt.Println(sub)
	}

	// Output:
	//
	// subsub.sub.domain.tld
	// sub.domain.tld
	// domain.tld
	// tld
}
