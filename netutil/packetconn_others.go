//go:build !linux

package netutil

import "net"

// wrapPacketConn wraps c to make it a [SessionPacketConn], if the OS supports
// that.
func wrapPacketConn(c net.PacketConn) (wrapped net.PacketConn) {
	return c
}
