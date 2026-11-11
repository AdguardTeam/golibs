//go:build linux

package netutil_test

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionPacketConn(t *testing.T) {
	t.Parallel()

	require.True(t, t.Run("v4", func(t *testing.T) {
		ipLocal := netutil.IPv4Localhost().AsSlice()
		ipZero := netutil.IPv4Zero().String() + ":0"

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			testSessionPacketConn(ctx, c, "udp4", ipZero, ipLocal)
		}, testTimeout, testTimeout/10)
	}))

	require.True(t, t.Run("v6", func(t *testing.T) {
		ipLocal := netutil.IPv6Localhost().AsSlice()
		skipIfIPv6Unsupported(t, ipLocal)

		ipZero := fmt.Sprintf("[%s]:0", netutil.IPv6Zero().String())
		ctx := testutil.ContextWithTimeout(t, testTimeout)
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			testSessionPacketConn(ctx, c, "udp6", ipZero, ipLocal)
		}, testTimeout, testTimeout/10)
	}))
}

// skipIfIPv6Unsupported skips the test if IPv6 is not supported in the current
// environment.  IPv6 is considered unsupported if binding a UDP6 socket to
// dstIP returns [syscall.EADDRNOTAVAIL].
func skipIfIPv6Unsupported(tb testing.TB, dstIP net.IP) {
	tb.Helper()

	conn, err := net.DialUDP("udp6", nil, &net.UDPAddr{IP: dstIP})
	if errors.Is(err, syscall.EADDRNOTAVAIL) {
		tb.Skipf("ipv6 seems to not be supported: %s", err)
	}

	require.NoError(tb, err)
	require.NoError(tb, conn.Close())
}

// testSessionPacketConn tests session-based packet exchange for given proto.
func testSessionPacketConn(
	ctx context.Context,
	t require.TestingT,
	proto string,
	addr string,
	dstIP net.IP,
) {
	lc, err := netutil.NewListenConfigWithOOB(nil)
	require.NoError(t, err)
	require.NotNil(t, lc)

	c, err := lc.ListenPacket(ctx, proto, addr)
	require.NoError(t, err)
	require.NotNil(t, c)
	defer func() {
		require.NoError(t, c.Close())
	}()

	deadline := time.Now().Add(1 * time.Second)
	err = c.SetDeadline(deadline)
	require.NoError(t, err)

	laddr := testutil.RequireTypeAssert[*net.UDPAddr](t, c.LocalAddr())
	require.NotNil(t, laddr)

	dstAddr := &net.UDPAddr{
		IP:   dstIP,
		Port: laddr.Port,
	}

	remoteConn, err := net.DialUDP(proto, nil, dstAddr)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, remoteConn.Close())
	}()

	err = remoteConn.SetDeadline(deadline)
	require.NoError(t, err)

	msg := []byte("hello")
	msgLen := len(msg)
	_, err = remoteConn.Write(msg)
	require.NoError(t, err)

	require.Implements(t, (*netutil.SessionPacketConn)(nil), c)

	buf := make([]byte, msgLen)
	n, sess, err := netutil.ReadFromSession(c, buf)
	require.NoError(t, err)

	assert.Equal(t, msgLen, n)
	assert.Equal(t, net.Addr(dstAddr), sess.LocalAddr())
	assert.Equal(t, remoteConn.LocalAddr(), sess.RemoteAddr())
	assert.Equal(t, msg, buf)

	respMsg := []byte("world")
	respMsgLen := len(respMsg)
	n, err = netutil.WriteToSession(c, respMsg, sess)
	require.NoError(t, err)

	assert.Equal(t, respMsgLen, n)

	buf = make([]byte, respMsgLen)
	n, err = remoteConn.Read(buf)
	require.NoError(t, err)

	assert.Equal(t, respMsgLen, n)
	assert.Equal(t, respMsg, buf)
}
