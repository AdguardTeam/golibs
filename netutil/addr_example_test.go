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
	// uint16(12345)
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
	fmt.Printf("%#v\n", netutil.Subdomains("subsub.sub.domain.tld"))

	fmt.Println()

	fmt.Printf("%#v\n", netutil.Subdomains(""))

	// Output:
	//
	// []string{"subsub.sub.domain.tld", "sub.domain.tld", "domain.tld", "tld"}
	//
	// []string(nil)
}

func ExampleIsSubdomain() {
	printResult := func(name string, isImmSub bool) {
		fmt.Printf("%-14s: %5t\n", name, isImmSub)
	}

	printResult("same domain", netutil.IsSubdomain("sub.example.com", "example.com"))
	printResult("not immediate", netutil.IsSubdomain("sub.sub.example.com", "example.com"))
	printResult("empty", netutil.IsSubdomain("", ""))
	printResult("same", netutil.IsSubdomain("example.com", "example.com"))
	printResult("dot only", netutil.IsSubdomain(".example.com", "example.com"))
	printResult("backwards", netutil.IsSubdomain("example.com", "sub.example.com"))
	printResult("other domain", netutil.IsSubdomain("sub.example.com", "example.org"))
	printResult("similar 1", netutil.IsSubdomain("sub.myexample.com", "example.org"))
	printResult("similar 2", netutil.IsSubdomain("sub.example.com", "myexample.org"))

	// Output:
	// same domain   :  true
	// not immediate :  true
	// empty         : false
	// same          : false
	// dot only      : false
	// backwards     : false
	// other domain  : false
	// similar 1     : false
	// similar 2     : false
}

func ExampleIsImmediateSubdomain() {
	printResult := func(name string, isImmSub bool) {
		fmt.Printf("%-14s: %5t\n", name, isImmSub)
	}

	printResult("same domain", netutil.IsImmediateSubdomain("sub.example.com", "example.com"))
	printResult("empty", netutil.IsImmediateSubdomain("", ""))
	printResult("same", netutil.IsImmediateSubdomain("example.com", "example.com"))
	printResult("dot only", netutil.IsImmediateSubdomain(".example.com", "example.com"))
	printResult("backwards", netutil.IsImmediateSubdomain("example.com", "sub.example.com"))
	printResult("other domain", netutil.IsImmediateSubdomain("sub.example.com", "example.org"))
	printResult("not immediate", netutil.IsImmediateSubdomain("sub.sub.example.com", "example.com"))
	printResult("similar 1", netutil.IsImmediateSubdomain("sub.myexample.com", "example.org"))
	printResult("similar 2", netutil.IsImmediateSubdomain("sub.example.com", "myexample.org"))

	// Output:
	// same domain   :  true
	// empty         : false
	// same          : false
	// dot only      : false
	// backwards     : false
	// other domain  : false
	// not immediate : false
	// similar 1     : false
	// similar 2     : false
}
