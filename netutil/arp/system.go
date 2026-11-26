package arp

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/osutil/executil"
)

// defaultARPFilePath is the default ARP cache file path.  It is used only on
// Linux.
const defaultARPFilePath = "proc/net/arp"

// Config is the configuration structure for [*System].
type Config struct {
	// Logger is used to log ARP operations.  If it is nil, [slog.Default] is
	// used.
	Logger *slog.Logger

	// CmdCons is used for constructing OS commands.  If it is nil,
	// [executil.SystemCommandConstructor] is used.
	CmdCons executil.CommandConstructor

	// RootFS is a root file system.  It is used on Linux to open ARP cache
	// file.  If it is nil, [osutil.RootDirFS] is used.
	RootFS fs.FS

	// ARPFilePath is a path to ARP cache file.  It is used only on Linux.  If
	// it is empty, "proc/net/arp" is used.
	ARPFilePath string
}

// System is a default implementation of the [Interface] which uses system
// commands to query OS-specific ARP.  It must be created using [NewSystem].
type System struct {
	logger      *slog.Logger
	cmdCons     executil.CommandConstructor
	rootFS      fs.FS
	arpFilePath string
}

// NewSystem returns the properly initialized *System.  If c is nil, default
// values will be used.
func NewSystem(c *Config) (s *System) {
	if c == nil {
		c = &Config{}
	}

	return &System{
		cmdCons: cmp.Or[executil.CommandConstructor](
			c.CmdCons,
			&executil.SystemCommandConstructor{},
		),
		logger:      cmp.Or(c.Logger, slog.Default()),
		rootFS:      cmp.Or[fs.FS](c.RootFS, osutil.RootDirFS()),
		arpFilePath: cmp.Or(c.ARPFilePath, defaultARPFilePath),
	}
}

// type check
var _ Interface = (*System)(nil)

// Table implements the [Interface] for *System.
func (s *System) Table(ctx context.Context) (t *Table, err error) {
	return s.table(ctx)
}

// execCmd executes a command by given path with given args, and returns the
// output scanner.
func (s *System) execCmd(
	ctx context.Context,
	path string,
	args []string,
) (sc *bufio.Scanner, err error) {
	var stdout bytes.Buffer
	err = executil.Run(
		ctx,
		s.cmdCons,
		&executil.CommandConfig{
			Path:   path,
			Args:   args,
			Stdout: &stdout,
		},
	)
	if err != nil {
		if code, ok := executil.ExitCodeFromError(err); ok {
			return nil, fmt.Errorf("running command: unexpected exit code %d", code)
		}

		return nil, fmt.Errorf("running command: %w", err)
	}

	return bufio.NewScanner(&stdout), nil
}

// parseArp parses the output of the "arp -a -n" command.  It is suitable for
// Linux, darwin and FreeBSD.  l and sc must not be nil.  The expected input
// format:
//
//	hostname (192.0.2.0) at ab:cd:ef:ab:cd:ef [ether] on enp0s3
func parseArp(ctx context.Context, l *slog.Logger, sc *bufio.Scanner) (t *Table, err error) {
	ns := make([]*Neighbor, 0)
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

		var n *Neighbor
		n, err = newNeighbor(host, ipStr[1:len(ipStr)-1], fields[3])
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
