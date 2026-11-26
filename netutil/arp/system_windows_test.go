//go:build windows

package arp_test

import (
	"net"

	"github.com/AdguardTeam/golibs/netutil/arp"
)

// testARPOutput is sample output from arp /a on Windows.
const testARPOutput = `

Interface: 192.0.2.1 --- 0x7
  Internet Address      Physical Address      Type
  ` + testIPv4Str + `             ab-cd-ef-ab-cd-ef     dynamic
  ` + testIPv6Str + `           ef-cd-ab-ef-cd-ab     static`

// testNeighbors contains the neighbors expected from [testARPOutput].
var testNeighbors = []*arp.Neighbor{{
	IP:  testIPv4,
	MAC: net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	IP:  testIPv6,
	MAC: net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}}
