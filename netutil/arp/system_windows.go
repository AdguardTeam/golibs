//go:build windows

package arp

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// table queries ARP and returns resulting table.  On Windows it runs "arp /a"
// command.  If err is nil, t is not nil.
func (s *System) table(ctx context.Context) (t *Table, err error) {
	sc, err := s.execCmd(ctx, "arp", []string{"/a"})
	if err != nil {
		// Don't wrap the error because it's informative enough as is.
		return nil, err
	}

	return parseResult(ctx, s.logger, sc)
}

// parseResult parses the output of the "arp /a" command on Windows.  sc must
// not be nil.  The expected input format (the first line is empty):
//
//	Interface: 192.0.2.0 --- 0x7
//	  Internet Address      Physical Address      Type
//	  192.0.2.1.            0a-00-27-00-00-00     dynamic
//	  192.0.2.255.          ff-ff-ff-ff-ff-ff     static
func parseResult(ctx context.Context, l *slog.Logger, sc *bufio.Scanner) (t *Table, err error) {
	ns := make([]*Neighbor, 0)
	for sc.Scan() {
		ln := sc.Text()
		if ln == "" {
			continue
		}

		fields := strings.Fields(ln)
		if len(fields) != 3 {
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
