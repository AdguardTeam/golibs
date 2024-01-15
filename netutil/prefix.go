package netutil

import (
	"bytes"
	"encoding"
	"net/netip"
)

// Prefix is a wrapper for [netip.Prefix] providing more functionality in
// encoding.  Unlike [netip.Prefix] it decodes IP addresses with unspecified
// mask bits (i.e. "127.0.0.1") as single-IP CIDR prefixes.
type Prefix struct {
	netip.Prefix
}

// type check
var _ encoding.TextUnmarshaler = (*Prefix)(nil)

// UnmarshalText implements [encoding.TextUnmarshaler] interface for *Prefix.
func (p *Prefix) UnmarshalText(b []byte) (err error) {
	if bytes.Contains(b, []byte("/")) {
		return p.Prefix.UnmarshalText(b)
	}

	var ip netip.Addr
	err = ip.UnmarshalText(b)
	if err != nil {
		// Don't wrap the error, since it's already informative enough as is.
		return err
	}

	p.Prefix = netip.PrefixFrom(ip, ip.BitLen())

	return nil
}

// UnembedPrefixes returns a slice of [netip.Prefix] from embed slice.
func UnembedPrefixes(embed []Prefix) (ps []netip.Prefix) {
	ps = make([]netip.Prefix, 0, len(embed))

	for _, p := range embed {
		ps = append(ps, p.Prefix)
	}

	return ps
}
