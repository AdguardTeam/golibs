//go:build windows

package arp

import (
	"bufio"
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// newARPDB returns new properly initialized *cmdARPDB for Windows.  c must not
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
		cmd:  "arp",
		args: []string{"/a"},
	}
}

// parseArpA parses the output of the "arp /a" command on Windows.  l and sc
// must not be nil.  The expected input format (the first line is empty):
//
//	Interface: 192.168.56.16 --- 0x7
//	  Internet Address      Physical Address      Type
//	  192.168.56.1          0a-00-27-00-00-00     dynamic
//	  192.168.56.255        ff-ff-ff-ff-ff-ff     static
func parseArpA(
	ctx context.Context,
	l *slog.Logger,
	sc *bufio.Scanner,
	lenHint int,
) (ns []*Neighbor) {
	ns = make([]*Neighbor, 0, lenHint)
	for sc.Scan() {
		ln := sc.Text()
		if ln == "" {
			continue
		}

		fields := strings.Fields(ln)
		if len(fields) != 3 {
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
