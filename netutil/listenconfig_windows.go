//go:build windows

package netutil

import (
	"net"
	"syscall"
)

// defaultListenControlWithSO is nil on Windows, because it doesn't support
// socket options.
var defaultListenControlWithSO func(_ syscall.RawConn) (_ error)

// setIPOpts sets the IPv4 and IPv6 options on a packet connection.
func setIPOpts(c net.PacketConn) (err error) {
	return nil
}
