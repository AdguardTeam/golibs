//go:build linux

package arp_test

import (
	"context"
	"net"
	"slices"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/arp"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeos/fakeexec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testARPOutput is sample output from 'arp -a' on Linux.
const (
	testARPOutput = `
invalid.mac (` + testOtherIPv4Str + `) at 12:34:56:78:910 on el0 ifscope [ethernet]
invalid.ip  (` + testInvalidIPv4Str + `) at ab:cd:ef:ab:cd:12 on ek0 ifscope [ethernet]
invalid.fmt 1 at 12:cd:ef:ab:cd:ef on er0 ifscope [ethernet]
? (` + testIPv4Str + `) at ab:cd:ef:ab:cd:ef on en0 ifscope [ethernet]
? (` + testIPv6Str + `) at ef:cd:ab:ef:cd:ab on em0 expires in 100 seconds [ethernet]`

	// testIPNeighOutput is the sample output from 'ip neigh' on Linux.
	testIPNeighOutput = `
` + testInvalidIPv4Str + ` dev enp0s3 lladdr aa:bb:cc:dd:ee:ff DELAY
` + testOtherIPv4Str + ` dev enp0s3 lladdr 12:34:56:78:910 DELAY
` + testIPv4Str + ` dev enp0s3 lladdr ab:cd:ef:ab:cd:ef DELAY
` + testIPv6Str + ` dev enp0s3 lladdr ef:cd:ab:ef:cd:ab router STALE`
)

// testNeighbors contains the neighbors expected from the Linux test data.
var testNeighbors = []*arp.Neighbor{{
	IP:  testIPv4,
	MAC: net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
}, {
	IP:  testIPv6,
	MAC: net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
}}

// Common commands for tests.
var (
	successCmd = &fakeexec.Command{
		OnStart: func(_ context.Context) (err error) { return nil },
		OnWait:  func(_ context.Context) (err error) { return nil },
	}

	failCmd = &fakeexec.Command{
		OnStart: func(_ context.Context) (err error) {
			return testExitError(osutil.ExitCodeFailure)
		},
		OnWait: func(_ context.Context) (err error) { return nil },
	}
)

// newTestFailCmdCons is a helper that returns mock command constructor that
// writes given stdout to the conf stdout if path and args match successPath and
// successArgs.  Constructor returns error otherwise.
func newTestFailCmdCons(
	tb testing.TB,
	stdout string,
	successPath string,
	successArgs []string,
) (cmdCons executil.CommandConstructor) {
	tb.Helper()

	return &fakeexec.CommandConstructor{
		OnNew: func(_ context.Context, conf *executil.CommandConfig) (c executil.Command, err error) {
			var success bool

			if conf.Path == successPath && slices.Equal(conf.Args, successArgs) {
				success = true

				_, err = conf.Stdout.Write([]byte(stdout))
				require.NoError(tb, err)
			}

			if success {
				return successCmd, nil
			}

			return failCmd, nil
		},
	}
}

func TestSystem_Table_linux(t *testing.T) {
	t.Parallel()

	const errMsg = `each table failed: opening "foo.bar": open foo.bar: no such file or directory
running command: unexpected exit code 1
running command: unexpected exit code 1`

	testCases := []struct {
		name        string
		arpFilePath string
		cmdPath     string
		stdout      string
		wantErrMsg  string
		cmdArgs     []string
	}{{
		name:        "fsys",
		arpFilePath: "proc_net_arp",
		cmdPath:     "",
		stdout:      "",
		wantErrMsg:  "",
		cmdArgs:     []string{},
	}, {
		name:        "arp",
		arpFilePath: "",
		cmdPath:     "arp",
		stdout:      testARPOutput,
		wantErrMsg:  "",
		cmdArgs:     []string{"-a", "-n"},
	}, {
		name:        "ip_neigh",
		arpFilePath: "",
		cmdPath:     "ip",
		stdout:      testIPNeighOutput,
		wantErrMsg:  "",
		cmdArgs:     []string{"neigh"},
	}, {
		name:        "fail",
		arpFilePath: "foo.bar",
		cmdPath:     "",
		stdout:      testIPNeighOutput,
		wantErrMsg:  errMsg,
		cmdArgs:     []string{""},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			conf := &arp.Config{
				CmdCons:     newTestFailCmdCons(t, tc.stdout, tc.cmdPath, tc.cmdArgs),
				ARPFilePath: tc.arpFilePath,
			}

			s := newTestSystem(t, conf)
			table, err := s.Table(testutil.ContextWithTimeout(t, testTimeout))
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrMsg == "" {
				wantTable := &arp.Table{Neighbors: testNeighbors}
				assert.Equal(t, wantTable, table)
			}
		})
	}
}
