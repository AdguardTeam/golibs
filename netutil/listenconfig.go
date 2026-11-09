package netutil

import (
	"context"
	"net"
	"syscall"
)

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

// type check
var _ ListenConfig = (*net.ListenConfig)(nil)

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
// [ListenConfig.ListenPacket].  It returns [errors.ErrUnsupported] on
// unsupported platforms.  If conf is nil, a default configuration is used.
func NewListenConfigWithOOB(conf *ControlConfig) (lc ListenConfig, err error) {
	return newListenConfigWithOOB(conf)
}
