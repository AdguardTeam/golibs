//go:build !linux

package sysresolv

import (
	"fmt"

	"github.com/AdguardTeam/golibs/netutil"
)

// validateAddress ensures the passed address is valid.
func validateAddress(address string) (err error) {
	_, err = netutil.SplitHost(address)
	if err != nil {
		// TODO(e.burkov): Maybe use a structured errBadAddrPassed to allow
		// unwrapping of the real error.
		return fmt.Errorf("%s: %w", err, errBadAddrPassed)
	}

	return nil
}
