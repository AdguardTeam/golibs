package arp

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/osutil/executil"
)

// cmdARPDB is the implementation of the [Interface] that uses command line to
// retrieve data.
type cmdARPDB struct {
	logger  *slog.Logger
	cmdCons executil.CommandConstructor
	parse   parseNeighborsFunc
	ns      *neighbors
	cmd     string
	args    []string
}

// type check
var _ Interface = (*cmdARPDB)(nil)

// Refresh implements the [Interface] interface for *cmdARPDB.
func (arp *cmdARPDB) Refresh(ctx context.Context) (err error) {
	defer func() { err = errors.Annotate(err, "cmd arpdb: %w") }()

	var stdout bytes.Buffer
	err = executil.Run(
		ctx,
		arp.cmdCons,
		&executil.CommandConfig{
			Path:   arp.cmd,
			Args:   arp.args,
			Stdout: &stdout,
		},
	)
	if err != nil {
		if code, ok := executil.ExitCodeFromError(err); ok {
			return fmt.Errorf("running command: unexpected exit code %d", code)
		}

		return fmt.Errorf("running command: %w", err)
	}

	sc := bufio.NewScanner(&stdout)
	ns := arp.parse(ctx, arp.logger, sc, arp.ns.len())
	if err = sc.Err(); err != nil {
		// TODO(e.burkov):  This error seems unreachable.  Investigate.
		return fmt.Errorf("scanning the output: %w", err)
	}

	arp.ns.reset(ns)

	return nil
}

// Table implements the [Interface] for *cmdARPDB.
func (arp *cmdARPDB) Table(ctx context.Context) (t *Table, err error) {
	return &Table{
		Neighbors: arp.ns.clone(),
	}, nil
}
