//go:build openbsd

package arp

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// table queries ARP and returns resulting table.  On OpenBSD it runs the
// "arp -a -n" command.  If err is nil, t is not nil.
func (s *System) table(ctx context.Context) (t *Table, err error) {
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

	return parseResult(ctx, s.logger, sc)
}

// parseResult parses the output of the "arp -a -n" command on OpenBSD.  l and
// sc must not be nil.  The expected input format:
//
//	Host        Ethernet Address  Netif Expire    Flags
//	192.0.2.0   ab:cd:ef:ab:cd:ef   em0 19m59s
func parseResult(ctx context.Context, l *slog.Logger, sc *bufio.Scanner) (t *Table, err error) {
	// Skip the header.
	if !sc.Scan() {
		return &Table{Neighbors: []*Neighbor{}}, nil
	}

	ns := make([]*Neighbor, 0)
	for sc.Scan() {
		ln := sc.Text()

		fields := strings.Fields(ln)
		if len(fields) < 2 {
			continue
		}

		var n *Neighbor
		n, err = newNeighbor("", fields[0], fields[1])
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
