// Package arp implements the Network Neighborhood Database.
package arp

import (
	"bufio"
	"cmp"
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/netip"
	"slices"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/service"
)

// Interface stores and refreshes the network neighbors table using ARP.
type Interface interface {
	// Refresher updates the stored data.  It must be safe for concurrent use.
	service.Refresher

	// Table returns current snapshot of the network neighbors table.  It must
	// be safe for concurrent use.
	Table(ctx context.Context) (t *Table, err error)
}

// Table contains snapshot of the network neighbors table.
type Table struct {
	// Neighbors is a slice of network neighbors.
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

// Clone returns the deep copy of n.
func (n *Neighbor) Clone() (clone *Neighbor) {
	return &Neighbor{
		Name: n.Name,
		IP:   n.IP,
		MAC:  slices.Clone(n.MAC),
	}
}

// Config is the configuration structure for [Interface] implementations.
type Config struct {
	// Logger is used to log ARP DB operations.  If it is nil, [slog.Default] is
	// used.
	Logger *slog.Logger

	// CmdCons is used for constructing OS commands.  If it is nil,
	// [executil.SystemCommandConstructor] is used.
	CmdCons executil.CommandConstructor

	// RootFS is a root file system.  It is used only for linux.  If it is nil,
	// [osutil.RootDirFS] is used.
	RootFS fs.FS
}

// New returns the [Interface] properly initialized for the OS.  If c is nil,
// default values will be used.
func New(c *Config) (arp Interface) {
	if c == nil {
		c = &Config{}
	}

	return newARPDB(&Config{
		CmdCons: cmp.Or[executil.CommandConstructor](
			c.CmdCons,
			&executil.SystemCommandConstructor{},
		),
		Logger: cmp.Or(c.Logger, slog.Default()),
		RootFS: cmp.Or[fs.FS](c.RootFS, osutil.RootDirFS()),
	})
}

// Empty is the [Interface] implementation that does nothing.
type Empty struct{}

// type check
var _ Interface = Empty{}

// Refresh implements the [Interface] for Empty.  It does nothing and always
// returns nil error.
func (Empty) Refresh(_ context.Context) (err error) { return nil }

// Table implements the [Interface] for Empty.  t and err are always nil.
func (Empty) Table(_ context.Context) (t *Table, err error) { return nil, nil }

// Static is the [Interface] implementation that returns underlying slice as a
// neighbor table.  All the elements must not be nil.
type Static []*Neighbor

// type check
var _ Interface = Static{}

// Refresh implements the [Interface] for Static.  It does nothing and always
// returns nil error.
func (Static) Refresh(_ context.Context) (err error) { return nil }

// Table implements the [Interface] for Static.  It uses underlying slice deep
// copy as a neighbor table.  err is always nil.
func (s Static) Table(_ context.Context) (t *Table, err error) {
	return &Table{Neighbors: cloneNeighbors(s)}, nil
}

// cloneNeighbors returns a deep copy of orig.
func cloneNeighbors(orig []*Neighbor) (n []*Neighbor) {
	n = make([]*Neighbor, 0, len(orig))

	for _, neighbor := range orig {
		n = append(n, neighbor.Clone())
	}

	return n
}

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

// neighbors is the helper type that stores neighbors to avoid copying its
// methods among all the [Interface] implementations.
type neighbors struct {
	mu *sync.RWMutex
	ns []*Neighbor
}

// len returns the length of the neighbors slice.  It's safe for concurrent use.
func (ns *neighbors) len() (l int) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return len(ns.ns)
}

// clone returns a deep copy of the underlying neighbors slice.  It's safe for
// concurrent use.
func (ns *neighbors) clone() (cloned []*Neighbor) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return cloneNeighbors(ns.ns)
}

// reset replaces the underlying slice with the new one.  Every element of with
// must not be nil.  It's safe for concurrent use.
func (ns *neighbors) reset(with []*Neighbor) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.ns = with
}

// parseNeighborsFunc parses the text from sc as if it'd be an output of some
// ARP-related command.  lenHint is a hint for the size of the allocated slice
// of Neighbors.  l and sc must not be nil.
type parseNeighborsFunc func(
	ctx context.Context,
	l *slog.Logger,
	sc *bufio.Scanner,
	lenHint int,
) (ns []*Neighbor)

// composite is the [Interface] that combines several [Interface]
// implementations and consequently switches between those.
type composite struct {
	// arps is the set of [Interface] implementations to range through.
	arps []Interface
	neighbors
}

// newComposite returns a properly initialized *composite.  It begins refreshing
// from the first of arps.
func newComposite(arps ...Interface) (arp *composite) {
	return &composite{
		arps: arps,
		mu:   &sync.RWMutex{},
		ns:   make([]*Neighbor, 0),
	}
}

// type check
var _ Interface = (*composite)(nil)

// Refresh implements the [Interface] for *composite.
func (arp *composite) Refresh(ctx context.Context) (err error) {
	var errs []error

	for _, a := range arp.arps {
		err = a.Refresh(ctx)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		var t *Table
		t, err = a.Table(ctx)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		arp.reset(t.Neighbors)

		return nil
	}

	return errors.Annotate(errors.Join(errs...), "each arpdb failed: %w")
}

// Table implements the [Interface] for *composite.
//
// TODO(e.burkov):  Think of a way to avoid cloning the slice twice.
func (arp *composite) Table(ctx context.Context) (t *Table, err error) {
	return &Table{
		Neighbors: arp.clone(),
	}, nil
}
