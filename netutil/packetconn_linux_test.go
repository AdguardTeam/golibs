//go:build linux

package netutil_test

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

func TestSessionPacketConn(t *testing.T) {
	t.Parallel()

	require.True(t, t.Run("v4", func(t *testing.T) {
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			ipLocal := netutil.IPv4Localhost().AsSlice()
			ipZero := netutil.IPv4Zero().String() + ":0"

			assert.False(t, testSessionPacketConn(t, "udp4", ipZero, ipLocal))
		}, testTimeout, testTimeout/10)
	}))

	require.True(t, t.Run("v6", func(t *testing.T) {
		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			ipLocal := netutil.IPv6Localhost().AsSlice()
			ipZero := fmt.Sprintf("[%s]:0", netutil.IPv6Zero().String())

			assert.False(t, testSessionPacketConn(t, "udp6", ipZero, ipLocal))
		}, testTimeout, testTimeout/10)
	}))
}

// testSessionPacketConn tests session-based packet exchange for given proto.
// It returns true if a network operation times out.
func testSessionPacketConn(tb testing.TB, proto, addr string, dstIP net.IP) (isTimeout bool) {
	tb.Helper()

	lc, err := netutil.NewListenConfigWithOOB(nil)
	require.NoError(tb, err)
	require.NotNil(tb, lc)

	c, err := lc.ListenPacket(tb.Context(), proto, addr)
	if isTimeoutOrFail(tb, err) {
		return true
	}

	require.NotNil(tb, c)
	testutil.CleanupAndRequireSuccess(tb, c.Close)

	deadline := time.Now().Add(1 * time.Second)
	err = c.SetDeadline(deadline)
	require.NoError(tb, err)

	laddr := testutil.RequireTypeAssert[*net.UDPAddr](tb, c.LocalAddr())
	require.NotNil(tb, laddr)

	dstAddr := &net.UDPAddr{
		IP:   dstIP,
		Port: laddr.Port,
	}

	remoteConn, err := net.DialUDP(proto, nil, dstAddr)
	if proto == "udp6" && errors.Is(err, syscall.EADDRNOTAVAIL) {
		// Some CI machines have IPv6 disabled.
		tb.Skipf("ipv6 seems to not be supported: %s", err)
	} else if isTimeoutOrFail(tb, err) {
		return true
	}

	testutil.CleanupAndRequireSuccess(tb, remoteConn.Close)

	err = remoteConn.SetDeadline(deadline)
	require.NoError(tb, err)

	msg := []byte("hello")
	msgLen := len(msg)
	_, err = remoteConn.Write(msg)
	if isTimeoutOrFail(tb, err) {
		return true
	}

	require.Implements(tb, (*netutil.SessionPacketConn)(nil), c)

	buf := make([]byte, msgLen)
	n, sess, err := netutil.ReadFromSession(c, buf)
	if isTimeoutOrFail(tb, err) {
		return true
	}

	assert.Equal(tb, msgLen, n)
	assert.Equal(tb, net.Addr(dstAddr), sess.LocalAddr())
	assert.Equal(tb, remoteConn.LocalAddr(), sess.RemoteAddr())
	assert.Equal(tb, msg, buf)

	respMsg := []byte("world")
	respMsgLen := len(respMsg)
	n, err = netutil.WriteToSession(c, respMsg, sess)
	if isTimeoutOrFail(tb, err) {
		return true
	}

	assert.Equal(tb, respMsgLen, n)

	buf = make([]byte, respMsgLen)
	n, err = remoteConn.Read(buf)
	if isTimeoutOrFail(tb, err) {
		return true
	}

	assert.Equal(tb, respMsgLen, n)
	assert.Equal(tb, respMsg, buf)

	return false
}

// isTimeoutOrFail is a helper function that returns true if err is a timeout
// error and also calls [require.NoError] on err.
func isTimeoutOrFail(tb testing.TB, err error) (ok bool) {
	tb.Helper()

	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}

	require.NoError(tb, err)

	return false
}
