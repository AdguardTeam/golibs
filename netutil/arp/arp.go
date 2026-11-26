// Package arp implements the platform-specific ARP querying.
package arp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"slices"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
)

// Interface describes the interface for querying OS-specific ARP and obtaining
// network neighbors table.
type Interface interface {
	// Table queries platform-specific ARP and returns the table of network
	// neighbors.
	Table(ctx context.Context) (t *Table, err error)
}

// Table contains snapshot of the network neighbors table.
type Table struct {
	// Neighbors is a slice of network neighbors.  All elements must not be nil.
	Neighbors []*Neighbor
}

// Neighbor represents single network neighbor reported by ARP.
type Neighbor struct {
	// Name is the hostname of the neighbor.  Empty name is valid since not each
	// implementation of ARP is able to retrieve that.
	Name string

	// IP contains the neighbor IP address.  It can be both IPv4 or IPv6.
	IP netip.Addr

	// MAC contains the neighbor MAC address.
	MAC net.HardwareAddr
}

// Clone returns a deep copy of the neighbor.
func (n *Neighbor) Clone() (clone *Neighbor) {
	return &Neighbor{
		Name: n.Name,
		IP:   n.IP,
		MAC:  slices.Clone(n.MAC),
	}
}

// newNeighbor returns the new initialized [Neighbor] by parsing string
// representations of IP and MAC addresses.
func newNeighbor(host, ipStr, macStr string) (n *Neighbor, err error) {
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, fmt.Errorf("parsing neighbor ip: %w", err)
	}

	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, fmt.Errorf("parsing neighbor mac: %w", err)
	}

	return &Neighbor{
		Name: host,
		IP:   ip,
		MAC:  mac,
	}, nil
}

// Empty is the [Interface] implementation that does nothing.
type Empty struct{}

// type check
var _ Interface = Empty{}

// Table implements the [Interface] for Empty.  t and err are always nil.
func (Empty) Table(_ context.Context) (t *Table, err error) { return nil, nil }

// validatedHostname returns h if it's a valid hostname, or an empty string
// otherwise, logging the validation error.
func validatedHostname(ctx context.Context, logger *slog.Logger, h string) (host string) {
	err := netutil.ValidateHostname(h)
	if err != nil {
		logger.DebugContext(ctx, "parsing host of arp output", slogutil.KeyError, err)

		return ""
	}

	return h
}
