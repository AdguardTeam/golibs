package netutil_test

import (
	"bytes"
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func ExampleIPv4Localhost() {
	fmt.Println(netutil.IPv4Localhost())

	// Output:
	// 127.0.0.1
}

func ExampleIPv6Localhost() {
	fmt.Println(netutil.IPv6Localhost())

	// Output:
	// ::1
}

func ExampleZeroPrefix() {
	fmt.Println(netutil.ZeroPrefix(netutil.AddrFamilyIPv4))
	fmt.Println(netutil.ZeroPrefix(netutil.AddrFamilyIPv6))

	// Invalid value for [netutil.AddrFamily].
	func() {
		defer func() { fmt.Println(recover()) }()

		fmt.Println(netutil.ZeroPrefix(1234))
	}()

	// Output:
	// 0.0.0.0/0
	// ::/0
	// netutil.ZeroPrefix: bad address family 1234
}

func ExampleIPToAddr() {
	ip := net.ParseIP("1.2.3.4")
	addr, err := netutil.IPToAddr(ip, netutil.AddrFamilyIPv4)
	fmt.Printf("%q, error: %v\n", addr, err)

	addr, err = netutil.IPToAddr(ip, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", addr, err)

	ip = net.ParseIP("1234::5678")
	addr, err = netutil.IPToAddr(ip, netutil.AddrFamilyIPv4)
	fmt.Printf("%q, error: %v\n", addr, err)

	addr, err = netutil.IPToAddr(ip, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", addr, err)

	ip = nil
	addr, err = netutil.IPToAddr(ip, netutil.AddrFamilyIPv4)
	fmt.Printf("%q, error: %v\n", addr, err)

	ip = net.IP{1, 2, 3, 4, 5}
	addr, err = netutil.IPToAddr(ip, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", addr, err)

	pt := testutil.PanicT{}
	require.Panics(pt, func() {
		_, _ = netutil.IPToAddr(ip, 42)
	})

	// Output:
	// "1.2.3.4", error: <nil>
	// "::ffff:1.2.3.4", error: <nil>
	// "invalid IP", error: bad ipv4 net.IP 1234::5678
	// "1234::5678", error: <nil>
	// "invalid IP", error: nil ip
	// "invalid IP", error: bad net.IP value <nil>
}

func ExampleIPToAddrNoMapped() {
	ip := net.ParseIP("1.2.3.4")
	addr, err := netutil.IPToAddrNoMapped(ip)
	fmt.Printf("%q, error: %v\n", addr, err)

	ip = net.IP{1, 2, 3, 4}
	addr, err = netutil.IPToAddrNoMapped(ip)
	fmt.Printf("%q, error: %v\n", addr, err)

	ip = net.ParseIP("1234::5678")
	addr, err = netutil.IPToAddrNoMapped(ip)
	fmt.Printf("%q, error: %v\n", addr, err)

	// Output:
	// "1.2.3.4", error: <nil>
	// "1.2.3.4", error: <nil>
	// "1234::5678", error: <nil>
}

func ExampleIPNetToPrefix() {
	_, n, _ := net.ParseCIDR("1.2.3.0/24")
	pref, err := netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv4)
	fmt.Printf("%q, error: %v\n", pref, err)

	pref, err = netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", pref, err)

	_, n, _ = net.ParseCIDR("1234::/72")
	pref, err = netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv4)
	fmt.Printf("%q, error: %v\n", pref, err)

	pref, err = netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", pref, err)

	n.Mask = bytes.Repeat([]byte{0xff}, 17)
	pref, err = netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", pref, err)

	n = nil
	pref, err = netutil.IPNetToPrefix(n, netutil.AddrFamilyIPv6)
	fmt.Printf("%q, error: %v\n", pref, err)

	// Output:
	// "1.2.3.0/24", error: <nil>
	// "::ffff:1.2.3.0/24", error: <nil>
	// "invalid Prefix", error: bad ip for subnet 1234::/72: bad ipv4 net.IP 1234::
	// "1234::/72", error: <nil>
	// "invalid Prefix", error: bad subnet <nil>
	// "invalid Prefix", error: nil subnet
}

func ExampleIPNetToPrefixNoMapped() {
	_, n, _ := net.ParseCIDR("1.2.3.0/24")
	pref, err := netutil.IPNetToPrefixNoMapped(n)
	fmt.Printf("%q, error: %v\n", pref, err)

	n = &net.IPNet{
		IP:   net.IP{1, 2, 3, 0},
		Mask: net.CIDRMask(24, 32),
	}
	pref, err = netutil.IPNetToPrefixNoMapped(n)
	fmt.Printf("%q, error: %v\n", pref, err)

	_, n, _ = net.ParseCIDR("1234::/72")
	pref, err = netutil.IPNetToPrefixNoMapped(n)
	fmt.Printf("%q, error: %v\n", pref, err)

	n = nil
	pref, err = netutil.IPNetToPrefixNoMapped(n)
	fmt.Printf("%q, error: %v\n", pref, err)

	// Output:
	// "1.2.3.0/24", error: <nil>
	// "1.2.3.0/24", error: <nil>
	// "1234::/72", error: <nil>
	// "invalid Prefix", error: nil subnet
}

func ExampleNetAddrToAddrPort() {
	tcpAddr := &net.TCPAddr{
		// NOTE: This produces a 16-byte form of the IPv4 address.
		IP:   net.IPv4(1, 2, 3, 4),
		Port: 56789,
	}

	addrPort := netutil.NetAddrToAddrPort(tcpAddr)
	fmt.Println(addrPort)
	fmt.Println(addrPort.Addr().Is4())
	fmt.Println(addrPort.Addr().Is4In6())

	// Output:
	// 1.2.3.4:56789
	// true
	// false
}
