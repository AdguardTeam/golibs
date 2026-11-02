//go:build unix

package netutil_test

import (
	"syscall"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// Localhost addrs for tests.
var (
	localhost4 = netutil.JoinHostPort(netutil.IPv4Localhost().String(), 0)
	localhost6 = netutil.JoinHostPort(netutil.IPv6Localhost().String(), 0)
)

func TestNewListenConfigWithOOB(t *testing.T) {
	lc := netutil.NewListenConfigWithOOB(nil)
	require.NotNil(t, lc)

	t.Run("ipv4", func(t *testing.T) {
		c, err := lc.ListenPacket(t.Context(), "udp4", localhost4)
		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)
		scConner := testutil.RequireTypeAssert[netutil.SyscallConner](t, c)

		sc, err := scConner.SyscallConn()
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT)
			require.NoError(t, opErr)

			// TODO(a.garipov): Rewrite this to use actual expected values for
			// each OS.
			assert.NotEqual(t, 0, val)
		})
		require.NoError(t, err)
	})

	t.Run("ipv6", func(t *testing.T) {
		c, err := lc.ListenPacket(t.Context(), "udp6", localhost6)
		if errors.Is(err, syscall.EADDRNOTAVAIL) {
			// Some CI machines have IPv6 disabled.
			t.Skipf("ipv6 seems to not be supported: %s", err)
		}

		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)

		scConner := testutil.RequireTypeAssert[netutil.SyscallConner](t, c)

		sc, err := scConner.SyscallConn()
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT)
			require.NoError(t, opErr)

			assert.NotEqual(t, 0, val)
		})
		require.NoError(t, err)
	})
}

func TestNewListenConfigWithOOB_withSO(t *testing.T) {
	const (
		sndBufSize = 10000
		rcvBufSize = 20000
	)

	cf := func(_, _ string, c syscall.RawConn) (err error) {
		pt := testutil.NewPanicT(t)

		return c.Control(func(fdPtr uintptr) {
			fd := int(fdPtr)

			require.NoError(pt, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1))
			require.NoError(pt, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_RCVBUF, rcvBufSize))
			require.NoError(pt, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_SNDBUF, sndBufSize))
		})
	}

	lc := netutil.NewListenConfigWithOOB(&netutil.ControlConfig{Control: cf})
	require.NotNil(t, lc)

	t.Run("ipv4", func(t *testing.T) {
		c, err := lc.ListenPacket(t.Context(), "udp4", localhost4)
		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)

		scConner := testutil.RequireTypeAssert[netutil.SyscallConner](t, c)

		sc, err := scConner.SyscallConn()
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_SNDBUF)
			require.NoError(t, opErr)

			// TODO(a.garipov): Rewrite this to use actual expected values for
			// each OS.
			assert.LessOrEqual(t, sndBufSize, val)
		})
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF)
			require.NoError(t, opErr)

			assert.LessOrEqual(t, rcvBufSize, val)
		})
		require.NoError(t, err)
	})

	t.Run("ipv6", func(t *testing.T) {
		c, err := lc.ListenPacket(t.Context(), "udp6", localhost6)
		if errors.Is(err, syscall.EADDRNOTAVAIL) {
			// Some CI machines have IPv6 disabled.
			t.Skipf("ipv6 seems to not be supported: %s", err)
		}

		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)

		scConner := testutil.RequireTypeAssert[netutil.SyscallConner](t, c)

		sc, err := scConner.SyscallConn()
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_SNDBUF)
			require.NoError(t, opErr)

			// TODO(a.garipov): Rewrite this to use actual expected values for
			// each OS.
			assert.LessOrEqual(t, sndBufSize, val)
		})
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF)
			require.NoError(t, opErr)

			assert.LessOrEqual(t, rcvBufSize, val)
		})
		require.NoError(t, err)
	})
}
