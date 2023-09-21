// Package sysresolv provides cross-platform functionality to discover DNS
// resolvers currently used by the system.
package sysresolv

// HostGenFunc is a function used for generating hostnames to check the system
// DNS.  Implementations must be safe for concurrent use.
type HostGenFunc func() (host string)

// Resolvers helps to work with local resolvers' addresses provided by OS.
type Resolvers interface {
	// Addrs returns the local resolvers' addresses.  It must be safe for
	// concurrent use.
	Addrs() (addrs []string)

	// Refresh refreshes the local resolvers' addresses cache.  It must be safe
	// for concurrent use.
	Refresh() (err error)
}

// NewSystemResolvers returns a Resolvers instance that uses f to generate fake
// hosts for checking.  If f is nil, a default generator is used.
func NewSystemResolvers(f HostGenFunc) (r Resolvers, err error) {
	r = newSystemResolvers(f)

	// Fill cache.
	err = r.Refresh()
	if err != nil {
		return nil, err
	}

	return r, nil
}
