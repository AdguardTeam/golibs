package sysresolv

import "github.com/AdguardTeam/golibs/errors"

const (
	// errFakeDial is an error which [dialFunc] is expected to return.
	errFakeDial errors.Error = "this error signals the successful dialFunc work"

	// errBadAddrPassed is returned by validateDialedHost when the host is not
	// an IP address and not an IP address with a port.
	errBadAddrPassed errors.Error = "address is neither ip, nor ip:port"
)
