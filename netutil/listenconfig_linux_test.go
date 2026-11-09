//go:build linux

package netutil_test

import (
	"net"
	"net/netip"
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
	localhost4 = netip.AddrPortFrom(netutil.IPv4Localhost(), 0).String()
	localhost6 = netip.AddrPortFrom(netutil.IPv6Localhost(), 0).String()
)

// syscallConner is an interface for connections that provide access to their
// underlying network connections through [syscall.RawConn].
type syscallConner interface {
	SyscallConn() (c syscall.RawConn, err error)
}

func TestNewListenConfigWithOOB(t *testing.T) {
	t.Parallel()

	lc, err := netutil.NewListenConfigWithOOB(nil)
	require.NoError(t, err)
	require.NotNil(t, lc)

	require.True(t, t.Run("ipv4", func(t *testing.T) {
		var c net.PacketConn
		c, err = lc.ListenPacket(t.Context(), "udp4", localhost4)
		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)
		scConner := testutil.RequireTypeAssert[syscallConner](t, c)

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
	}))

	require.True(t, t.Run("ipv6", func(t *testing.T) {
		var c net.PacketConn
		c, err = lc.ListenPacket(t.Context(), "udp6", localhost6)
		if errors.Is(err, syscall.EADDRNOTAVAIL) {
			// Some CI machines have IPv6 disabled.
			t.Skipf("ipv6 seems to not be supported: %s", err)
		}

		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)

		scConner := testutil.RequireTypeAssert[syscallConner](t, c)

		sc, err := scConner.SyscallConn()
		require.NoError(t, err)

		err = sc.Control(func(fd uintptr) {
			val, opErr := unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT)
			require.NoError(t, opErr)

			assert.NotEqual(t, 0, val)
		})
		require.NoError(t, err)
	}))
}

func TestNewListenConfigWithOOB_withSO(t *testing.T) {
	t.Parallel()

	const (
		sndBufSize = 10000
		rcvBufSize = 20000
	)

	cf := func(_, _ string, c syscall.RawConn) (err error) {
		return c.Control(func(fdPtr uintptr) {
			fd := int(fdPtr)

			require.NoError(t, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1))
			require.NoError(t, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_RCVBUF, rcvBufSize))
			require.NoError(t, unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_SNDBUF, sndBufSize))
		})
	}

	lc, err := netutil.NewListenConfigWithOOB(&netutil.ControlConfig{Control: cf})
	require.NoError(t, err)
	require.NotNil(t, lc)

	t.Run("ipv4", func(t *testing.T) {
		c, err := lc.ListenPacket(t.Context(), "udp4", localhost4)
		require.NoError(t, err)
		require.NotNil(t, c)
		testutil.CleanupAndRequireSuccess(t, c.Close)

		scConner := testutil.RequireTypeAssert[syscallConner](t, c)

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

		scConner := testutil.RequireTypeAssert[syscallConner](t, c)

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
