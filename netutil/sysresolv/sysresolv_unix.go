//go:build unix

package sysresolv

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/stringutil"
)

// defaultHostGen is the default method of generating host for
// [systemResolvers.Refresh].
func defaultHostGen() (host string) {
	// TODO(e.burkov): Use strings.Builder.
	return fmt.Sprintf("test%d.org", time.Now().UnixNano())
}

// systemResolvers is a default implementation of the Resolvers interface.
type systemResolvers struct {
	// addrs is the set that contains cached local resolvers' addresses.
	addrs *stringutil.Set

	// resolver is used to fetch the resolvers' addresses.
	resolver *net.Resolver

	// hostGenFunc generates hosts to resolve.
	hostGenFunc HostGenFunc

	// addrsLock protects addrs.
	addrsLock sync.RWMutex
}

const (
	// errBadAddrPassed is returned when dialFunc can't parse an IP address.
	errBadAddrPassed errors.Error = "the passed string is not a valid IP address"

	// errFakeDial is an error which dialFunc is expected to return.
	errFakeDial errors.Error = "this error signals the successful dialFunc work"

	// errUnexpectedHostFormat is returned by validateDialedHost when the host has
	// more than one percent sign.
	errUnexpectedHostFormat errors.Error = "unexpected host format"
)

// newSystemResolvers returns a new system resolvers checker for Unix systems.
func newSystemResolvers(hostGenFunc HostGenFunc) (sr Resolvers) {
	if hostGenFunc == nil {
		hostGenFunc = defaultHostGen
	}

	s := &systemResolvers{
		resolver: &net.Resolver{
			PreferGo: true,
		},
		hostGenFunc: hostGenFunc,
		addrs:       stringutil.NewSet(),
	}

	s.resolver.Dial = s.dialFunc

	return s
}

// type check
var _ Resolvers = (*systemResolvers)(nil)

// Refresh implements the [Resolvers] interface for *systemResolvers.
func (sr *systemResolvers) Refresh() (err error) {
	defer func() { err = errors.Annotate(err, "systemResolvers: %w") }()

	_, err = sr.resolver.LookupHost(context.Background(), sr.hostGenFunc())
	dnserr := &net.DNSError{}
	if errors.As(err, &dnserr) && dnserr.Err == errFakeDial.Error() {
		return nil
	}

	return err
}

// validateDialedHost validates the host used by resolvers in dialFunc.
func validateDialedHost(host string) (err error) {
	defer func() { err = errors.Annotate(err, "parsing %q: %w", host) }()

	parts := strings.Split(host, "%")
	switch len(parts) {
	case 1:
		// host
	case 2:
		// Remove the zone and check the IP address part.
		host = parts[0]
	default:
		return errUnexpectedHostFormat
	}

	if _, err = netutil.ParseIP(host); err != nil {
		return errBadAddrPassed
	}

	return nil
}

// dockerEmbeddedDNS is the address of Docker's embedded DNS server.
//
// See https://github.com/moby/moby/blob/v1.12.0/docs/userguide/networking/dockernetworks.md.
const dockerEmbeddedDNS = "127.0.0.11"

// dialFunc gets the resolver's address and puts it into internal cache.
func (sr *systemResolvers) dialFunc(_ context.Context, _, address string) (_ net.Conn, err error) {
	// Just validate the passed address is a valid IP.
	host, err := netutil.SplitHost(address)
	if err != nil {
		// TODO(e.burkov): Maybe use a structured errBadAddrPassed to allow
		// unwrapping of the real error.
		return nil, fmt.Errorf("%s: %w", err, errBadAddrPassed)
	}

	// Exclude Docker's embedded DNS server, as it may cause recursion if the
	// container is set as the host system's default DNS server.
	//
	// See https://github.com/AdguardTeam/AdGuardHome/issues/3064.
	//
	// TODO(a.garipov): Perhaps only do this when we are in the container?
	// Maybe use an environment variable?
	if host == dockerEmbeddedDNS {
		return nil, errFakeDial
	}

	err = validateDialedHost(host)
	if err != nil {
		return nil, fmt.Errorf("validating dialed host: %w", err)
	}

	sr.addrsLock.Lock()
	defer sr.addrsLock.Unlock()

	sr.addrs.Add(host)

	return nil, errFakeDial
}

// Addrs implements the [Resolvers] interface for *systemResolvers.
func (sr *systemResolvers) Addrs() (addrs []string) {
	sr.addrsLock.RLock()
	defer sr.addrsLock.RUnlock()

	return sr.addrs.Values()
}
