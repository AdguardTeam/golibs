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

func ExamplePrefixFromReversedAddr() {
	fmt.Println(netutil.PrefixFromReversedAddr(ipv4RevGood))
	fmt.Println(netutil.PrefixFromReversedAddr(ipv4NetRevGood))
	fmt.Println(netutil.PrefixFromReversedAddr(ipv6RevGood))
	fmt.Println(netutil.PrefixFromReversedAddr(ipv6NetRevGood))

	fmt.Println(netutil.PrefixFromReversedAddr(ipv4Missing))
	fmt.Println(netutil.PrefixFromReversedAddr(ipv4Char))

	fmt.Println(netutil.PrefixFromReversedAddr(ipv6RevMany))
	fmt.Println(netutil.PrefixFromReversedAddr(ipv6RevCharHi))

	// Output:
	//
	// 1.2.3.4/32 <nil>
	// 10.0.0.0/8 <nil>
	// 1234::cdef/128 <nil>
	// 1234::1000:0:0:0/68 <nil>
	// invalid Prefix bad arpa domain name ".0.0.127.in-addr.arpa": bad domain name label "": domain name label is empty
	// invalid Prefix bad arpa domain name "1.0.z.127.in-addr.arpa": ParseAddr("1.0.z.127"): unexpected character (at "z.127")
	// invalid Prefix bad arpa domain name "4.3.2.1.dbc.b.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.3.2.1.ip6.arpa": not a full reversed ip address
	// invalid Prefix bad arpa domain name "4.3.2.1.d.c.b.a.0.z.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.3.2.1.ip6.arpa": bad arpa domain name rune 'z'
}

func ExamplePrefixFromReversedAddr_domainOnly() {
	a := `in-addr.arpa`
	_, err := netutil.PrefixFromReversedAddr(a)

	fmt.Println(err)

	// Output:
	//
	// bad arpa domain name "in-addr.arpa": not a reversed ip network
}
