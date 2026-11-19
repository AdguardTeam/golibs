//go:build darwin

package netutil

import (
	"net"
)

// newIPv4RespOOB returns an encoded control-message assuming origDst is a valid
// IPv4.  On darwin it returns an empty slice.
func newIPv4RespOOB(_ net.IP) (b []byte) {
	// Do not set the IPv4 source address via OOB, because it can cause the
	// address to become unspecified on darwin.
	//
	// See https://github.com/AdguardTeam/AdGuardHome/issues/2807.
	//
	// TODO(e.burkov): Develop a workaround to make it write OOB only when
	// listening on an unspecified address.
	return []byte{}
}
