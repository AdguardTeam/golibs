package netutil

import "net/netip"

// PreferIPv4 compares two addresses, preferring IPv4 addresses over IPv6 ones.
// It's intended to be used as a compare function in [slices.SortFunc].  Invalid
// addresses are sorted near the end.
func PreferIPv4(a, b netip.Addr) (res int) {
	return prefer(a, b, netip.Addr.Is4)
}

// PreferIPv6 compares two addresses, preferring IPv6 addresses over IPv4 ones.
// It's intended to be used as a compare function in [slices.SortFunc].  Invalid
// addresses are sorted near the end.
func PreferIPv6(a, b netip.Addr) (res int) {
	return prefer(a, b, netip.Addr.Is6)
}

// prefer compares two addresses, preferring the one that satisfies famFunc over
// the other.  If both addresses satisfy famFunc, it compares them using the
// [netip.Addr.Compare] method.
func prefer(a, b netip.Addr, famFunc func(netip.Addr) (ok bool)) (res int) {
	if !a.IsValid() {
		return 1
	} else if !b.IsValid() {
		return -1
	}

	if aFam := famFunc(a); aFam == famFunc(b) {
		return a.Compare(b)
	} else if aFam {
		return -1
	}

	return 1
}
