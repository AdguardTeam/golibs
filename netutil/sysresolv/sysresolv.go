// Package sysresolv provides cross-platform functionality to discover DNS
// resolvers currently used by the system.
package sysresolv

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// unit is a convenient type alias for empty struct.
type unit = struct{}

// HostGenFunc is a function used for generating hostnames to check the system
// DNS.  The generated hosts should be unique to avoid resolver's cache.
// Implementations must be safe for concurrent use.
type HostGenFunc func() (hostname string)

// SystemResolvers is a default implementation of the Resolvers interface.
type SystemResolvers struct {
	// lastUpd is the moment when the last update started.
	lastUpd time.Time

	// updMu protects addresses and lastUpd.
	updMu *sync.RWMutex

	// generateHost generates unique hosts to resolve.
	generateHost HostGenFunc

	// addresses is the most recent set of cached local resolvers' addresses.
	addresses []netip.AddrPort

	// defaultPort is the default port to use when parsing an address without
	// port.
	defaultPort uint16
}

// NewSystemResolvers returns a SystemResolvers instance that uses genHost to
// generate fake hosts for dialing, see [HostGenFunc].  The default generator is
// used, if genHost is nil.  The defaultPort is used when resolvers are provided
// without a port number.
func NewSystemResolvers(genHost HostGenFunc, defaultPort uint16) (sr *SystemResolvers, err error) {
	if genHost == nil {
		genHost = defaultHostGenFunc
	}

	sr = &SystemResolvers{
		// TODO(e.burkov):  Probably we should use Unix epoch here.
		lastUpd:      time.Now(),
		updMu:        &sync.RWMutex{},
		addresses:    []netip.AddrPort{},
		generateHost: genHost,
		defaultPort:  defaultPort,
	}

	// Fill the cache initially.
	err = sr.Refresh()
	if err != nil {
		return nil, err
	}

	return sr, nil
}

// Addrs returns all the collected resolvers' addresses.  Caller must clone the
// returned slice before modifying it.  It is safe for concurrent use.
func (sr *SystemResolvers) Addrs() (addrs []netip.AddrPort) {
	sr.updMu.RLock()
	defer sr.updMu.RUnlock()

	return sr.addresses
}

// Refresh updates the internal cache of the resolvers' addresses if no error
// occurred.  It is safe for concurrent use.
func (sr *SystemResolvers) Refresh() (err error) {
	startTime := time.Now()

	defer func() { err = errors.Annotate(err, "system resolvers: %w") }()

	set, err := sr.collectResolvers()
	if err != nil {
		return err
	}

	addrs := maps.Keys(set)
	slices.SortFunc(addrs, compareAddrPorts)

	sr.updMu.Lock()
	defer sr.updMu.Unlock()

	// Don't use [time.Before] here, since on Windows the values returned by
	// [time.Now] has a much lower precision than on Unix, which may omit the
	// initial update.
	if !sr.lastUpd.After(startTime) {
		sr.addresses = addrs
	}

	return nil
}

// compareAddrPorts compares two [netip.AddrPort]s.  It's used for sorting.
func compareAddrPorts(a, b netip.AddrPort) (res int) {
	res = a.Addr().Compare(b.Addr())
	if res != 0 {
		return res
	}

	return int(a.Port()) - int(b.Port())
}

// collectResolvers returns the set of resolvers' addresses used by the system.
func (sr *SystemResolvers) collectResolvers() (set map[netip.AddrPort]unit, err error) {
	setMu := sync.Mutex{}
	set = map[netip.AddrPort]unit{}

	dialFunc := func(_ context.Context, _, address string) (_ net.Conn, err error) {
		addrPort, err := sr.parse(address)
		if err != nil {
			return nil, err
		}

		setMu.Lock()
		defer setMu.Unlock()

		set[addrPort] = unit{}

		return nil, errFakeDial
	}

	resolver := net.Resolver{
		PreferGo: true,
		Dial:     dialFunc,
	}

	_, err = resolver.LookupHost(context.Background(), sr.generateHost())
	dnsErr := &net.DNSError{}
	if !errors.As(err, &dnsErr) || dnsErr.Err != errFakeDial.Error() {
		return nil, err
	}

	return set, nil
}

// parse returns the [netip.AddrPort] parsed from the passed address, using the
// preconfigured default port, if address doesn't contain one.  It also
// immediately returns [errFakeDial] if the address is valid, but useless.
func (sr *SystemResolvers) parse(address string) (addrPort netip.AddrPort, err error) {
	err = validateAddress(address)
	if err != nil {
		return netip.AddrPort{}, err
	}

	addrPort, err = netip.ParseAddrPort(address)
	if err != nil {
		ip, ipErr := netip.ParseAddr(address)
		if ipErr != nil {
			return netip.AddrPort{}, fmt.Errorf("%s: %w", err, errBadAddrPassed)
		}

		addrPort = netip.AddrPortFrom(ip, sr.defaultPort)
	}

	return addrPort, nil
}

// defaultHostGenFunc is the default implementation of the [HostGenFunc].
func defaultHostGenFunc() (hostname string) {
	// This code is written after Unix time value of 1000000000000000000, so
	// it's safe to assume that the length of the generated hostname will be
	// constant for a long time.
	const expectedLen = len("host0000000000000000000.test")

	b := strings.Builder{}
	b.Grow(expectedLen)

	// It's safe to ignore errors from writing to [strings.Builder], since it
	// always returns nil errors and only panics on OOM.
	_, _ = b.WriteString("host")

	// Format the integer manually to avoid allocations.  The order of digits
	// doesn't really matter.
	for n := time.Now().UnixNano(); n > 0; n /= 10 {
		_ = b.WriteByte(byte(n%10) + '0')
	}

	_, _ = b.WriteString(".test")

	return b.String()
}
