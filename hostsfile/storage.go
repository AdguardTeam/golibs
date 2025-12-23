package hostsfile

import (
	"net/netip"
)

// Storage indexes the hosts file records.
type Storage interface {
	// ByAddr returns the hostnames for the given address.
	ByAddr(addr netip.Addr) (names []string)

	// ByName returns the addresses for the given hostname.
	ByName(name string) (addrs []netip.Addr)
}
