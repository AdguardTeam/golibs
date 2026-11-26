//go:build darwin || freebsd

package arp

import (
	"context"
)

// table queries ARP and returns resulting table.  On darwin and FreeBSD it runs
// "arp -a -n" command.  If err is nil, t is not nil.
func (s *System) table(ctx context.Context) (t *Table, err error) {
	// Use -n flag to avoid resolving the hostnames of the neighbors.  By
	// default ARP attempts to resolve the hostnames via DNS.  See man 8
	// arp.
	//
	// See also https://github.com/AdguardTeam/AdGuardHome/issues/3157.
	sc, err := s.execCmd(ctx, "arp", []string{"-a", "-n"})
	if err != nil {
		// Don't wrap the error because it's informative enough as is.
		return nil, err
	}

	return parseArp(ctx, s.logger, sc)
}
