//go:build darwin || freebsd

package arp

import (
	"bufio"
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// newARPDB returns new properly initialized *cmdARPDB for darwin and FreeBSD.
// c must not be nil.  c must be valid.
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

// parseArpA parses the output of the "arp -a -n" command on darwin and FreeBSD.
// l and sc must not be nil.  The expected input format is:
//
//	host.name (192.168.0.1) at ff:ff:ff:ff:ff:ff on en0 ifscope [ethernet]
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
