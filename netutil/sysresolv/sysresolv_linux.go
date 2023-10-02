//go:build linux

package sysresolv

import (
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

// dockerEmbeddedDNS is the address of Docker's embedded DNS server.
//
// See https://github.com/moby/moby/blob/v1.12.0/docs/userguide/networking/dockernetworks.md.
const dockerEmbeddedDNS = "127.0.0.11"

// validateAddress ensures the passed address is valid and desired.
func validateAddress(address string) (err error) {
	host, err := netutil.SplitHost(address)
	if err != nil {
		// TODO(e.burkov): Maybe use a structured errBadAddrPassed to allow
		// unwrapping of the real error.
		return fmt.Errorf("%s: %w", err, errBadAddrPassed)
	}

	// Exclude Docker's embedded DNS server, as it may cause recursion if the
	// container is set as the host system's default DNS server.
	//
	// See https://github.com/AdguardTeam/AdGuardHome/issues/3064.
	//
	// TODO(a.garipov): Perhaps only do this when we are in the container?
	// Maybe use an environment variable?
	if host == dockerEmbeddedDNS {
		return errFakeDial
	}

	return nil
}
