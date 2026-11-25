//go:build openbsd

package arp

import (
	"bufio"
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// newARPDB returns new properly initialized *cmdARPDB for OpenBSD.  c must not
// be nil.  c must be valid.
func newARPDB(c *Config) (arp *cmdARPDB) {
	return &cmdARPDB{
		logger:  c.Logger,
		cmdCons: c.CmdCons,
		parse:   parseArpA,
		ns: &neighbors{
			mu: &sync.RWMutex{},
			ns: make([]*Neighbor, 0),
		},
		cmd: "arp",
		// Use -n flag to avoid resolving the hostnames of the neighbors.  By
		// default ARP attempts to resolve the hostnames via DNS.  See man 8
		// arp.
		//
		// See also https://github.com/AdguardTeam/AdGuardHome/issues/3157.
		args: []string{"-a", "-n"},
	}
}

// parseArpA parses the output of the "arp -a -n" command on OpenBSD.  l and sc
// must not be nil.  The expected input format:
//
//	Host        Ethernet Address  Netif Expire    Flags
//	192.168.1.1 ab:cd:ef:ab:cd:ef   em0 19m59s
func parseArpA(
	ctx context.Context,
	l *slog.Logger,
	sc *bufio.Scanner,
	lenHint int,
) (ns []*Neighbor) {
	// Skip the header.
	if !sc.Scan() {
		return nil
	}

	ns = make([]*Neighbor, 0, lenHint)
	for sc.Scan() {
		ln := sc.Text()

		fields := strings.Fields(ln)
		if len(fields) < 2 {
			continue
		}

		n, err := newNeighbor("", fields[0], fields[1])
		if err != nil {
			l.DebugContext(ctx, "parsing arp output", "line", ln, slogutil.KeyError, err)

			continue
		}

		ns = append(ns, n)
	}

	return ns
}
