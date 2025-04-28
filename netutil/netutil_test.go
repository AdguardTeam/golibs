package netutil_test

import (
	"net"
	"net/netip"
)

// Common test IPs.  Do not mutate.
var (
	testIPv4 = net.IP{1, 2, 3, 4}
	testIPv6 = net.IP{
		0x12, 0x34, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xcd, 0xef,
	}

	testIPv4Addr = netip.MustParseAddr("1.2.3.4")
	testIPv6Addr = netip.MustParseAddr("1234::cdef")

	testIPv4Prefix = netip.MustParsePrefix("1.2.3.4/32")
	testIPv6Prefix = netip.MustParsePrefix("1234::cdef/128")
)
