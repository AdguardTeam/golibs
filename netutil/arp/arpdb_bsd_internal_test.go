//go:build darwin || freebsd

package arp

import (
	"net"
	"net/netip"
)

// testOtherIPv6Str is the address of an additional IPv6 neighbor in the BSD fixture.
const testOtherIPv6Str = "2001:db8::3"

// testARPOutput is sample output from arp -a on macOS and FreeBSD.
const testARPOutput = `
invalid.mac (` + testOtherIPv4Str + `) at 12:34:56:78:910 on el0 ifscope [ethernet]
invalid.ip  (` + testInvalidIPv4Str + `) at ab:cd:ef:ab:cd:12 on ek0 ifscope [ethernet]
invalid.fmt 1 at 12:cd:ef:ab:cd:ef on er0 ifscope [ethernet]
hostname.one (` + testIPv4Str + `) at ab:cd:ef:ab:cd:ef on en0 ifscope [ethernet]
hostname.two (` + testIPv6Str + `) at ef:cd:ab:ef:cd:ab on em0 expires in 1198 seconds [ethernet]
? (` + testOtherIPv6Str + `) at aa:bb:cc:dd:ee:ff on ej0 expires in 1918 seconds [ethernet]
`

// testNeighbors contains the neighbors expected from [testARPOutput].
var testNeighbors = []*Neighbor{{
	Name: "hostname.one",
	IP:   testIPv4,
	MAC:  net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	Name: "hostname.two",
	IP:   testIPv6,
	MAC:  net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}, {
	Name: "",
	IP:   netip.MustParseAddr(testOtherIPv6Str),
	MAC:  net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
}}
