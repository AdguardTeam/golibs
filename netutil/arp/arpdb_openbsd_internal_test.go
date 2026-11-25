//go:build openbsd

package arp

import "net"

// testARPOutput is sample output from arp -a on OpenBSD.
const testARPOutput = `
Host        Ethernet Address  Netif Expire    Flags
` + testInvalidIPv4Str + ` aa:bb:cc:dd:ee:ff   em0 permanent
` + testOtherIPv4Str + `   12:34:56:78:910     em0 permanent
` + testIPv4Str + `   ab:cd:ef:ab:cd:ef   em0 19m56s
` + testIPv6Str + ` ef:cd:ab:ef:cd:ab   em0 permanent l
`

// testNeighbors contains the neighbors expected from [testARPOutput].
var testNeighbors = []*Neighbor{{
	IP:  testIPv4,
	MAC: net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	IP:  testIPv6,
	MAC: net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}}
