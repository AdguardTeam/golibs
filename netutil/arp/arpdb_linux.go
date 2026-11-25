//go:build linux

package arp

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/netip"
	"strings"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/stringutil"
)

// newARPDB returns new properly initialized *composite for linux.  c must not
// be nil.  c must be valid.
//
// TODO(f.setrakov): !! Support Open Wrt.
func newARPDB(c *Config) (arp *composite) {
	// Use the common storage among the implementations.
	ns := &neighbors{
		mu: &sync.RWMutex{},
		ns: make([]*Neighbor, 0),
	}

	// Linux implementation uses the result of the first successful run of the
	// following actions:
	//  - /proc/net/arp read.
	//  - "arp -a -n".
	//  - "ip neigh".
	return newComposite(
		&fsysARPDB{
			ns:       ns,
			fsys:     c.RootFS,
			filename: "proc/net/arp",
		},
		&cmdARPDB{
			logger:  c.Logger,
			cmdCons: c.CmdCons,
			parse:   parseArpA,
			ns:      ns,
			cmd:     "arp",
			// NOTE: Use -n flag to avoid resolving the hostnames of the
			// neighbors. By default ARP attempts to resolve the hostnames via
			// DNS.  See man 8 arp.
			//
			// See also https://github.com/AdguardTeam/AdGuardHome/issues/3157.
			args: []string{"-a", "-n"},
		},
		&cmdARPDB{
			logger:  c.Logger,
			cmdCons: c.CmdCons,
			parse:   parseIPNeigh,
			ns:      ns,
			cmd:     "ip",
			args:    []string{"neigh"},
		},
	)
}

// fsysARPDB is the [Interface] implementation that uses the ARP cache file to
// update the database.
type fsysARPDB struct {
	ns       *neighbors
	fsys     fs.FS
	filename string
}

// type check
var _ Interface = (*fsysARPDB)(nil)

// Refresh implements the [Interface] interface for *fsysARPDB.
func (arp *fsysARPDB) Refresh(_ context.Context) (err error) {
	var f fs.File
	f, err = arp.fsys.Open(arp.filename)
	if err != nil {
		return fmt.Errorf("opening %q: %w", arp.filename, err)
	}
	defer func() {
		err = errors.WithDeferred(err, f.Close())
	}()

	sc := bufio.NewScanner(f)

	// Skip the header.
	if !sc.Scan() {
		return nil
	} else if err = sc.Err(); err != nil {
		return err
	}

	ns := make([]*Neighbor, 0, arp.ns.len())
	for sc.Scan() {
		n := parseNeighbor(sc.Text())
		if n != nil {
			ns = append(ns, n)
		}
	}

	arp.ns.reset(ns)

	return nil
}

// parseNeighbor parses line into *Neighbor.  If parsing failed n is nil.
func parseNeighbor(line string) (n *Neighbor) {
	fields := stringutil.SplitTrimmed(line, " ")
	if len(fields) != 6 {
		return nil
	}

	ip, err := netip.ParseAddr(fields[0])
	if err != nil || ip.IsUnspecified() {
		return nil
	}

	mac, err := net.ParseMAC(fields[3])
	if err != nil {
		return nil
	}

	return &Neighbor{
		IP:  ip,
		MAC: mac,
	}
}

// Table implements the [Interface] interface for *fsysARPDB.
func (arp *fsysARPDB) Table(ctx context.Context) (t *Table, err error) {
	return &Table{
		Neighbors: arp.ns.clone(),
	}, nil
}

// parseArpA parses the output of the "arp -a -n" command on Linux.  l and sc
// must not be nil.  The expected input format:
//
//	hostname (192.168.1.1) at ab:cd:ef:ab:cd:ef [ether] on enp0s3
func parseArpA(
	ctx context.Context,
	l *slog.Logger,
	sc *bufio.Scanner,
	lenHint int,
) (ns []*Neighbor) {
	ns = make([]*Neighbor, 0, lenHint)
	for sc.Scan() {
		ln := sc.Text()

		fields := strings.Fields(ln)
		if len(fields) < 4 {
			continue
		}

		ipStr := fields[1]
		if len(ipStr) < 2 {
			continue
		}

		host := validatedHostname(ctx, l, fields[0])
		n, err := newNeighbor(host, ipStr[1:len(ipStr)-1], fields[3])
		if err != nil {
			l.DebugContext(ctx, "parsing arp output", "line", ln, slogutil.KeyError, err)

			continue
		}

		ns = append(ns, n)
	}

	return ns
}

// parseIPNeigh parses the output of the "ip neigh" command on Linux.  l and sc
// must not be nil.  The expected input format:
//
//	192.168.1.1 dev enp0s3 lladdr ab:cd:ef:ab:cd:ef REACHABLE
func parseIPNeigh(
	ctx context.Context,
	l *slog.Logger,
	sc *bufio.Scanner,
	lenHint int,
) (ns []*Neighbor) {
	ns = make([]*Neighbor, 0, lenHint)
	for sc.Scan() {
		ln := sc.Text()

		fields := strings.Fields(ln)
		if len(fields) < 5 {
			continue
		}

		n, err := newNeighbor("", fields[0], fields[4])
		if err != nil {
			l.DebugContext(ctx, "parsing arp output", "line", ln, slogutil.KeyError, err)

			continue
		}

		ns = append(ns, n)
	}

	return ns
}
