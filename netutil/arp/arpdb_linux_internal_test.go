//go:build linux

package arp

import (
	"net"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testARPOutput is sample output from 'arp -a' on Linux.
const testARPOutput = `
invalid.mac (` + testOtherIPv4Str + `) at 12:34:56:78:910 on el0 ifscope [ethernet]
invalid.ip  (` + testInvalidIPv4Str + `) at ab:cd:ef:ab:cd:12 on ek0 ifscope [ethernet]
invalid.fmt 1 at 12:cd:ef:ab:cd:ef on er0 ifscope [ethernet]
? (` + testIPv4Str + `) at ab:cd:ef:ab:cd:ef on en0 ifscope [ethernet]
? (` + testIPv6Str + `) at ef:cd:ab:ef:cd:ab on em0 expires in 100 seconds [ethernet]`

// testARPOutput is sample output from 'ip neigh' on Linux.
const ipNeighOutput = `
` + testInvalidIPv4Str + ` dev enp0s3 lladdr aa:bb:cc:dd:ee:ff DELAY
` + testOtherIPv4Str + ` dev enp0s3 lladdr 12:34:56:78:910 DELAY
` + testIPv4Str + ` dev enp0s3 lladdr ab:cd:ef:ab:cd:ef DELAY
` + testIPv6Str + ` dev enp0s3 lladdr ef:cd:ab:ef:cd:ab router STALE`

// testNeighbors contains the neighbors expected from the Linux test fixtures.
var testNeighbors = []*Neighbor{{
	IP:  testIPv4,
	MAC: net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	IP:  testIPv6,
	MAC: net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}}

func TestFSysARPDB(t *testing.T) {
	t.Parallel()

	require.NoError(t, fstest.TestFS(testdata, "proc_net_arp"))

	a := &fsysARPDB{
		ns: &neighbors{
			mu: &sync.RWMutex{},
			ns: make([]*Neighbor, 0),
		},
		fsys:     testdata,
		filename: "proc_net_arp",
	}

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	err := a.Refresh(ctx)
	require.NoError(t, err)

	table, err := a.Table(ctx)
	require.NoError(t, err)

	assert.Equal(t, testNeighbors, table.Neighbors)
}

func TestCmdARPDB_linux(t *testing.T) {
	t.Parallel()

	a := &cmdARPDB{
		logger:  slogutil.NewDiscardLogger(),
		cmdCons: newTestCommandConstructor(t, osutil.ExitCodeSuccess, ipNeighOutput, nil),
		parse:   parseIPNeigh,
		cmd:     "ip",
		args:    []string{"neigh"},
		ns: &neighbors{
			mu: &sync.RWMutex{},
			ns: make([]*Neighbor, 0),
		},
	}

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	err := a.Refresh(ctx)
	require.NoError(t, err)

	table, err := a.Table(ctx)
	require.NoError(t, err)
	assert.Equal(t, testNeighbors, table.Neighbors)
}
