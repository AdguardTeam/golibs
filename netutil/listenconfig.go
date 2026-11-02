package netutil

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"github.com/AdguardTeam/golibs/errors"
)

// SyscallConner is an interface for connections that provide access to their
// underlying network connections through [syscall.RawConn].
type SyscallConner interface {
	// SyscallConn returns a raw network connection for the connection.
	SyscallConn() (c syscall.RawConn, err error)
}

// ListenConfig is the interface that allows controlling socket options for used
// connections.
//
// This interface is modeled after [net.ListenConfig].
type ListenConfig interface {
	// Listen announces on the local network address.  It is a wrapper for
	// [net.ListenConfig.Listen].
	Listen(ctx context.Context, network, address string) (l net.Listener, err error)

	// ListenPacket announces on the local network address.  It is a wrapper for
	// [net.ListenConfig.ListenPacket].
	ListenPacket(ctx context.Context, network, address string) (c net.PacketConn, err error)
}

// ControlFunc is a callback that will be called after creating the network
// connection but before binding it to the operating system.  It is intended to
// be used for setting sock opts.  c must not be nil.
type ControlFunc func(network, address string, c syscall.RawConn) (err error)

// ControlConfig is the configuration of socket options.
type ControlConfig struct {
	// Control is a callback that will be used as a [net.ListenConfig.Control].
	Control ControlFunc
}

// defaultControlFunc is a default [ControlFunc] that sets appropriate default
// socket options.
func defaultControlFunc(_, _ string, c syscall.RawConn) (err error) {
	return defaultListenControlWithSO(c)
}

// NewListenConfig returns properly initialized ListenConfig.  If conf is nil, a
// default configuration is used.
func NewListenConfig(conf *ControlConfig) (lc ListenConfig) {
	if conf == nil {
		conf = &ControlConfig{Control: defaultControlFunc}
	}

	return &net.ListenConfig{Control: conf.Control}
}

// NewListenConfigWithOOB returns properly initialized ListenConfig.  Unlike
// [NewListenConfig], resulting config sets additional socket flags and
// processes the control-messages of connections created with
// [ListenConfig.ListenPacket].  If conf is nil, a default configuration is
// used.
func NewListenConfigWithOOB(conf *ControlConfig) (lc ListenConfig) {
	if conf == nil {
		conf = &ControlConfig{Control: defaultControlFunc}
	}

	return &listenConfigOOB{Control: conf.Control}
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

	return wrapPacketConn(c), nil
}
