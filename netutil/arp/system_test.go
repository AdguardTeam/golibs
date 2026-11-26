package arp_test

import (
	"cmp"
	"context"
	"fmt"
	"io/fs"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil/arp"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeos/fakeexec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

var (
	// testdata is the filesystem containing data for testing the package.
	testdata fs.FS = os.DirFS("./testdata")

	// testLogger is a common logger for tests.
	testLogger = slogutil.NewDiscardLogger()
)

// Common IP address strings for tests.
const (
	testIPv4Str        = "192.0.2.0"
	testIPv6Str        = "2001:db8::2"
	testOtherIPv4Str   = "192.0.2.1"
	testOtherIPv6Str   = "2001:db8::3"
	testInvalidIPv4Str = "192.0.2.256"
)

// Common IP addresses for tests.
var (
	testIPv4 = netip.MustParseAddr(testIPv4Str)
	testIPv6 = netip.MustParseAddr(testIPv6Str)
)

// newTestSystem returns a new *arp.System initialized with common test
// defaults.
func newTestSystem(tb testing.TB, c *arp.Config) (s *arp.System) {
	c = cmp.Or(c, &arp.Config{})

	return arp.NewSystem(&arp.Config{
		Logger:      testLogger,
		CmdCons:     cmp.Or[executil.CommandConstructor](c.CmdCons, executil.EmptyCommandConstructor{}),
		RootFS:      testdata,
		ARPFilePath: c.ARPFilePath,
	})
}

// newTestCommandConstructor is a helper function that returns a mock
// [executil.CommandConstructor] for tests.  This mock writes given stdout and
// returns cmdErr or [testExitError] with given code on wait.
func newTestCommandConstructor(
	tb testing.TB,
	code osutil.ExitCode,
	stdout string,
	cmdErr error,
) (cons *fakeexec.CommandConstructor) {
	tb.Helper()

	onNew := func(_ context.Context, conf *executil.CommandConfig) (c executil.Command, err error) {
		cmd := fakeexec.NewCommand()
		cmd.OnStart = func(_ context.Context) (err error) {
			if conf.Stdout != nil {
				_, err = conf.Stdout.Write([]byte(stdout))
				require.NoError(tb, err)
			}

			return nil
		}

		cmd.OnWait = func(_ context.Context) (err error) {
			if cmdErr != nil {
				return cmdErr
			}

			if code != osutil.ExitCodeSuccess {
				return testExitError(code)
			}

			return nil
		}

		return cmd, nil
	}

	return &fakeexec.CommandConstructor{OnNew: onNew}
}

// testExitError is a [executil.ExitCodeError] implementation for tests.
type testExitError int

// type check
var _ executil.ExitCodeError = testExitError(0)

// Error implements the [executil.ExitCodeError] interface for testExitError.
func (e testExitError) Error() (msg string) {
	return fmt.Sprintf("exit code: %d", e)
}

// ExitCode implements the [executil.ExitCodeError] interface for testExitError.
func (e testExitError) ExitCode() (code int) { return int(e) }

func TestSystem_Table(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		conf       *arp.Config
		wantTable  *arp.Table
		wantErrMsg string
	}{{
		name: "success",
		conf: &arp.Config{
			CmdCons: newTestCommandConstructor(t, osutil.ExitCodeSuccess, testARPOutput, nil),
		},
		wantTable:  &arp.Table{Neighbors: testNeighbors},
		wantErrMsg: "",
	}, {
		name: "cmd_error",
		conf: &arp.Config{
			CmdCons: newTestCommandConstructor(t, osutil.ExitCodeSuccess, "", assert.AnError),
		},
		wantTable:  nil,
		wantErrMsg: "running command: running: " + assert.AnError.Error(),
	}, {
		name: "bad_code",
		conf: &arp.Config{
			CmdCons: newTestCommandConstructor(t, osutil.ExitCodeFailure, "", nil),
		},
		wantTable:  nil,
		wantErrMsg: "running command: unexpected exit code 1",
	}, {
		name: "empty",
		conf: &arp.Config{
			CmdCons: newTestCommandConstructor(t, osutil.ExitCodeSuccess, "", nil),
		},
		wantTable:  &arp.Table{Neighbors: []*arp.Neighbor{}},
		wantErrMsg: "",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSystem(t, tc.conf)
			table, err := s.Table(testutil.ContextWithTimeout(t, testTimeout))
			if tc.wantErrMsg == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantErrMsg)
			}

			assert.Equal(t, tc.wantTable, table)
		})
	}
}
