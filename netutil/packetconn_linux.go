//go:build linux

package netutil

import (
	"net"

	"golang.org/x/net/ipv4"
)

// newIPv4RespOOB returns an encoded control-message assuming origDst is a valid
// IPv4.  origDst is expected to have a [net.IPv4len] length.
func newIPv4RespOOB(origDst net.IP) (b []byte) {
	return (&ipv4.ControlMessage{
		Src: origDst,
	}).Marshal()
}
