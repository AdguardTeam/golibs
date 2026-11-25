package arp

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/netip"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/osutil/executil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeos/fakeexec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

// Common IP address strings for tests.
const (
	testIPv4Str        = "192.0.2.2"
	testIPv6Str        = "2001:db8::2"
	testOtherIPv4Str   = "192.0.2.4"
	testOtherIPv6Str   = "2001:db8::3"
	testInvalidIPv4Str = "192.0.2.256"
)

// Common IP addresses for tests.
var (
	testIPv4 = netip.MustParseAddr(testIPv4Str)
	testIPv6 = netip.MustParseAddr(testIPv6Str)
)

// testdata is the filesystem containing data for testing the package.
var testdata fs.FS = os.DirFS("./testdata")

// newTestCommandConstructor is a helper function that returns a mock
// [executil.CommandConstructor] for tests.
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

// testARP is the mock implementation of [Interface] to use in tests.
type testARP struct {
	onRefresh func(ctx context.Context) (err error)
	onTable   func(ctx context.Context) (t *Table, err error)
}

// type check
var _ Interface = (*testARP)(nil)

// Refresh implements the [Interface] interface for *testARP.
func (arp *testARP) Refresh(ctx context.Context) (err error) {
	return arp.onRefresh(ctx)
}

// Table implements the [Interface] interface for *testARP.
func (arp *testARP) Table(ctx context.Context) (t *Table, err error) {
	return arp.onTable(ctx)
}

func TestComposite(t *testing.T) {
	t.Parallel()

	knownIP := netip.MustParseAddr(testOtherIPv4Str)
	knownMAC := net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF}

	succRefrCount, failRefrCount := 0, 0
	clnp := func() {
		succRefrCount, failRefrCount = 0, 0
	}

	succDB := &testARP{
		onRefresh: func(_ context.Context) (err error) { succRefrCount++; return nil },
		onTable: func(ctx context.Context) (t *Table, err error) {
			return &Table{
				Neighbors: []*Neighbor{{Name: "abc", IP: knownIP, MAC: knownMAC}},
			}, nil
		},
	}
	failDB := &testARP{
		onRefresh: func(_ context.Context) (err error) {
			failRefrCount++

			return errors.Error("refresh failed")
		},
		onTable: func(ctx context.Context) (t *Table, err error) { return nil, nil },
	}

	require.True(t, t.Run("begin_with_success", func(t *testing.T) {
		t.Cleanup(clnp)

		a := newComposite(succDB, failDB)

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err := a.Refresh(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, succRefrCount)
		assert.Zero(t, failRefrCount)

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, table.Neighbors)
	}))

	require.True(t, t.Run("begin_with_fail", func(t *testing.T) {
		t.Cleanup(clnp)

		a := newComposite(failDB, succDB)
		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err := a.Refresh(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, succRefrCount)
		assert.Equal(t, 1, failRefrCount)

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, table.Neighbors)
	}))

	require.True(t, t.Run("fail_only", func(t *testing.T) {
		t.Cleanup(clnp)

		wantMsg := "each arpdb failed: refresh failed\nrefresh failed"

		a := newComposite(failDB, failDB)
		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err := a.Refresh(ctx)
		require.Error(t, err)

		testutil.AssertErrorMsg(t, wantMsg, err)

		assert.Equal(t, 2, failRefrCount)

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.Empty(t, table.Neighbors)
	}))

	require.True(t, t.Run("fail_after_success", func(t *testing.T) {
		t.Cleanup(clnp)

		shouldFail := false
		unstableDB := &testARP{
			onRefresh: func(_ context.Context) (err error) {
				if shouldFail {
					err = errors.Error("unstable failed")
				}
				shouldFail = !shouldFail

				return err
			},
			onTable: func(ctx context.Context) (table *Table, err error) {
				if !shouldFail {
					return failDB.Table(ctx)
				}

				return succDB.Table(ctx)
			},
		}
		a := newComposite(unstableDB, succDB)
		ctx := testutil.ContextWithTimeout(t, testTimeout)

		// Unstable ARPDB should refresh successfully.
		err := a.Refresh(ctx)
		require.NoError(t, err)

		assert.Zero(t, succRefrCount)

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, table.Neighbors)

		// Unstable ARPDB should fail and the succDB should be used.
		err = a.Refresh(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, succRefrCount)

		table, err = a.Table(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, table.Neighbors)

		// Unstable ARPDB should refresh successfully again.
		err = a.Refresh(ctx)
		require.NoError(t, err)

		assert.Equal(t, 1, succRefrCount)

		table, err = a.Table(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, table.Neighbors)
	}))

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		a := newComposite()
		ctx := testutil.ContextWithTimeout(t, testTimeout)
		require.NoError(t, a.Refresh(ctx))

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.Empty(t, table.Neighbors)
	})
}

func TestCmdARPDB_arpa(t *testing.T) {
	t.Parallel()

	a := &cmdARPDB{
		logger: slogutil.NewDiscardLogger(),
		cmd:    "cmd",
		parse:  parseArpA,
		ns: &neighbors{
			mu: &sync.RWMutex{},
			ns: make([]*Neighbor, 0),
		},
	}

	require.True(t, t.Run("arp_a", func(t *testing.T) {
		a.cmdCons = newTestCommandConstructor(t, osutil.ExitCodeSuccess, testARPOutput, nil)

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err := a.Refresh(ctx)
		require.NoError(t, err)

		table, err := a.Table(ctx)
		require.NoError(t, err)

		assert.Equal(t, testNeighbors, table.Neighbors)
	}))

	require.True(t, t.Run("runcmd_error", func(t *testing.T) {
		a.cmdCons = newTestCommandConstructor(
			t,
			osutil.ExitCodeSuccess,
			"",
			errors.Error("can't run"),
		)

		err := a.Refresh(testutil.ContextWithTimeout(t, testTimeout))
		testutil.AssertErrorMsg(t, "cmd arpdb: running command: running: can't run", err)
	}))

	require.True(t, t.Run("bad_code", func(t *testing.T) {
		a.cmdCons = newTestCommandConstructor(t, osutil.ExitCodeFailure, "", nil)

		err := a.Refresh(testutil.ContextWithTimeout(t, testTimeout))
		testutil.AssertErrorMsg(
			t,
			"cmd arpdb: running command: unexpected exit code 1",
			err,
		)
	}))

	require.True(t, t.Run("empty", func(t *testing.T) {
		a.cmdCons = newTestCommandConstructor(t, osutil.ExitCodeSuccess, "", nil)

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		err := a.Refresh(ctx)
		require.NoError(t, err)

		table, err := a.Table(ctx)
		require.NoError(t, err)
		assert.Empty(t, table.Neighbors)
	}))
}
