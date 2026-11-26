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

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/stringutil"
)

// tableFunc represents the callback for querying Linux ARP.
type tableFunc func(ctx context.Context) (t *Table, err error)

// table queries ARP and returns resulting table.  On Linux, it returns the
// first successful result of the following actions:
//   - Reading ARP cache.
//   - Running "arp -a -n".
//   - Runnig "ip neigh".
//
// If err is nil, t is not nil.
//
// TODO(f.setrakov): !! Support OpenWrt.
func (s *System) table(ctx context.Context) (t *Table, err error) {
	tableFuncs := []tableFunc{
		s.tableFSys,
		s.tableARP,
		s.tableIPNeigh,
	}

	var errs []error
	for _, f := range tableFuncs {
		t, err = f(ctx)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		return t, nil
	}

	return nil, errors.Annotate(errors.Join(errs...), "each table failed: %w")
}

// tableFSys is a [tableFunc] that retrieves the network neighbors table from
// the ARP cache file.
func (s *System) tableFSys(ctx context.Context) (t *Table, err error) {
	var f fs.File
	f, err = s.rootFS.Open(s.arpFilePath)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", s.arpFilePath, err)
	}
	defer func() {
		err = errors.WithDeferred(err, f.Close())
	}()

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		return &Table{}, nil
	} else if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}

	ns := make([]*Neighbor, 0)
	for sc.Scan() {
		n := parseNeighbor(sc.Text())
		if n != nil {
			ns = append(ns, n)
		}
	}

	return &Table{Neighbors: ns}, nil
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

// tableARP is a [tableFunc] that retrieves the network neighbors table from
// "arp -a -n" run output.
func (s *System) tableARP(ctx context.Context) (t *Table, err error) {
	// Use -n flag to avoid resolving the hostnames of the neighbors.  By
	// default ARP attempts to resolve the hostnames via DNS.  See man 8
	// arp.
	//
	// See also https://github.com/AdguardTeam/AdGuardHome/issues/3157.
	sc, err := s.execCmd(ctx, "arp", []string{"-a", "-n"})
	if err != nil {
		// Don't wrap the error because it's informative enough as is.
		return nil, err
	}

	return parseArp(ctx, s.logger, sc)
}

// tableIPNeigh is a [tableFunc] that retrieves the network neighbors table from
// "ip neigh" run output.
func (s *System) tableIPNeigh(ctx context.Context) (t *Table, err error) {
	sc, err := s.execCmd(ctx, "ip", []string{"neigh"})
	if err != nil {
		// Don't wrap the error because it's informative enough as is.
		return nil, err
	}

	return parseIPNeigh(ctx, s.logger, sc)
}

// parseIPNeigh parses the output of the "ip neigh" command on Linux.  l and sc
// must not be nil.  The expected input format:
//
//	192.0.2.0 dev enp0s3 lladdr ab:cd:ef:ab:cd:ef REACHABLE
func parseIPNeigh(ctx context.Context, l *slog.Logger, sc *bufio.Scanner) (t *Table, err error) {
	ns := make([]*Neighbor, 0)
	for sc.Scan() {
		ln := sc.Text()

		fields := strings.Fields(ln)
		if len(fields) < 5 {
			continue
		}

		var n *Neighbor
		n, err = newNeighbor("", fields[0], fields[4])
		if err != nil {
			l.DebugContext(ctx, "parsing arp output", "line", ln, slogutil.KeyError, err)

			continue
		}

		ns = append(ns, n)
	}

	err = sc.Err()
	if err != nil {
		return nil, fmt.Errorf("scanning the output: %w", err)
	}

	return &Table{Neighbors: ns}, nil
}
