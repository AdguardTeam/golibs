package redisutil

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/gomodule/redigo/redis"
)

// Dialer is the interface for dialing Redis connections.
type Dialer interface {
	// DialContext creates and configures a connection with the given context.
	// c must not be in a special state (subscribed to pubsub channel,
	// transaction started, etc.).
	//
	// See [redis.Pool.DialContext].
	DialContext(ctx context.Context) (c redis.Conn, err error)
}

// DefaultDialerConfig is the configuration structure for a [*DefaultDialer].
type DefaultDialerConfig struct {
	// Addr is the address of the Redis server.  It must not be nil and must be
	// valid.
	Addr *netutil.HostPort

	// Resolver is used to resolve the hostname.  If nil, a default pure-Go
	// [*net.Resolver] is used.
	//
	// TODO(a.garipov):  Add the [netutil.Resolver] interface?
	Resolver *net.Resolver

	// Network is the network to dial.  If empty, "ip" is used.  If not empty,
	// must be one of:
	//   - "ip"
	//   - "ip4"
	//   - "ip6"
	Network string

	// DBIndex is the index of Redis database to use.  Zero is the default
	// index.
	DBIndex uint8
}

// DefaultDialer is the default implementation of the [Dialer] interface.
type DefaultDialer struct {
	addr     *netutil.HostPort
	resolver *net.Resolver
	net      string
	dbIdx    int
}

// NewDefaultDialer returns a properly initialized default dialer.  c should not
// be nil and should be valid.
func NewDefaultDialer(c *DefaultDialerConfig) (d *DefaultDialer, err error) {
	err = validate.NotNil("c", c)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return nil, err
	}

	err = validate.NotNil("c.Addr", c.Addr)
	if err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return nil, err
	}

	var errs []error
	err = validate.NotEmpty("c.Addr", *c.Addr)
	if err != nil {
		errs = append(errs, err)
	}

	if c.Resolver == nil {
		c.Resolver = &net.Resolver{
			PreferGo: true,
		}
	}

	switch c.Network {
	case "":
		c.Network = "ip"
	case "ip", "ip4", "ip6":
		// Go on.
	default:
		errs = append(errs, fmt.Errorf("c.Network: %w: %q", errors.ErrBadEnumValue, c.Network))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &DefaultDialer{
		addr:     c.Addr,
		resolver: c.Resolver,
		net:      c.Network,
		dbIdx:    int(c.DBIndex),
	}, nil
}

// type check
var _ Dialer = (*DefaultDialer)(nil)

// DialContext implements the [Dialer] interface for *DefaultDialer.
func (d *DefaultDialer) DialContext(ctx context.Context) (conn redis.Conn, err error) {
	ips, err := d.resolver.LookupNetIP(ctx, d.net, d.addr.Host)
	if err != nil {
		return nil, fmt.Errorf("looking up: %w", err)
	} else if len(ips) == 0 {
		panic(errors.Error(
			"stdlib contract violation: net.Resolver.LookupNetIP: 0 ips with no error",
		))
	}

	var errs []error
	port := d.addr.Port
	for i, ip := range ips {
		addrPort := netip.AddrPortFrom(ip, port)
		conn, err = redis.DialContext(ctx, "tcp", addrPort.String(), redis.DialDatabase(d.dbIdx))
		if err == nil {
			return conn, nil
		}

		err = fmt.Errorf("dialing ip %s at index %d and port %d: %w", ip, i, port, err)
		errs = append(errs, err)
	}

	return nil, errors.Join(errs...)
}
