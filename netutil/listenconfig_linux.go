//go:build linux

package netutil

import (
	"context"
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/errors"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// newListenConfigWithOOB returns new listen config which saves local address
// from out-of-band data.
func newListenConfigWithOOB(conf *ControlConfig) (lc ListenConfig, err error) {
	if conf == nil {
		conf = &ControlConfig{}
	}

	return &listenConfigOOB{Control: conf.Control}, nil
}

// type check
var _ ListenConfig = (*listenConfigOOB)(nil)

// listenConfigOOB is a wrapper around [net.ListenConfig] with modifications
// that set the control-message options on packet conns.
type listenConfigOOB struct {
	net.ListenConfig
}

// ListenPacket implements the [ListenConfig] interface for *listenConfigOOB.
// It sets the control-message flags to receive additional out-of-band data to
// correctly discover the source address when it listens to 0.0.0.0 as well as
// in situations when SO_BINDTODEVICE is used.
//
// network must be "udp", "udp4", or "udp6".
func (lc *listenConfigOOB) ListenPacket(
	ctx context.Context,
	network string,
	address string,
) (c net.PacketConn, err error) {
	c, err = lc.ListenConfig.ListenPacket(ctx, network, address)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return nil, err
	}

	err = setIPOpts(c)
	if err != nil {
		err = errors.WithDeferred(err, c.Close())

		return nil, fmt.Errorf("setting socket options: %w", err)
	}

	return &sessionPacketConn{
		udpConn: c.(udpConn),
	}, nil
}

// setIPOpts sets the IPv4 and IPv6 options on a packet connection.  c must not
// be nil.
func setIPOpts(c net.PacketConn) (err error) {
	// TODO(a.garipov): Returning an error only if both functions return one
	// (which is what module dns does as well) seems rather fragile.  Depending
	// on the OS, the valid errors are ENOPROTOOPT, EINVAL, and maybe others.
	// Investigate and make OS-specific versions to make sure we don't miss any
	// real errors.
	err6 := ipv6.NewPacketConn(c).SetControlMessage(ipv6.FlagDst|ipv6.FlagInterface, true)
	err4 := ipv4.NewPacketConn(c).SetControlMessage(ipv4.FlagDst|ipv4.FlagInterface, true)
	if err4 != nil && err6 != nil {
		return fmt.Errorf("setting ipv4 and ipv6 options: %w", errors.Join(err4, err6))
	}

	return nil
}
