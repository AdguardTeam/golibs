package netutil

import (
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/stringutil"
)

// Reversed ARPA Addresses

// fromHexByte converts a single hexadecimal ASCII digit character into an
// integer from 0 to 15.  For all other characters it returns 0xff.
//
// TODO(e.burkov):  This should be covered with tests after adding HasSuffixFold
// into stringutil.
func fromHexByte(c byte) (n byte) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0xff
	}
}

// ARPA reverse address domains.
const (
	arpaV4Suffix = ".in-addr.arpa"
	arpaV6Suffix = ".ip6.arpa"
)

// The maximum lengths of the ARPA-formatted reverse addresses.
//
// An example of IPv4 with a maximum length:
//
//	49.91.20.104.in-addr.arpa
//
// An example of IPv6 with a maximum length:
//
//	1.3.b.5.4.1.8.6.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.0.0.7.4.6.0.6.2.ip6.arpa
const (
	arpaV4MaxIPLen = len("000.000.000.000")
	arpaV6MaxIPLen = len("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")

	arpaV4MaxLen = arpaV4MaxIPLen + len(arpaV4Suffix)
	arpaV6MaxLen = arpaV6MaxIPLen + len(arpaV6Suffix)
)

// reverseIPv4 inverts the order of bytes in an IP address.
func reverseIPv4(ip [4]byte) (out [4]byte) {
	l := len(ip)
	for i := range ip[:l/2] {
		ip[i], ip[l-i-1] = ip[l-i-1], ip[i]
	}

	return ip
}

// ipv4FromReversed parses an IPv4 reverse address.  It assumes that arpa is a
// valid domain name.
func ipv4FromReversed(arpa string) (addr netip.Addr, err error) {
	arpaAddr, err := netip.ParseAddr(arpa)
	if err != nil {
		// Don't wrap the error, since it's already informative enough as is.
		return netip.Addr{}, err
	}

	ip := reverseIPv4(arpaAddr.As4())

	return netip.AddrFrom4(ip), nil
}

// ipv6FromReversed parses an IPv6 reverse address.  It assumes that arpa is a
// valid domain name.
func ipv6FromReversed(arpa string) (addr netip.Addr, err error) {
	const addrStep = len("0.0.")

	ip := [net.IPv6len]byte{}
	for i := range ip {
		// Get the two half-byte and merge them together.  Validate the dots
		// between them since while arpa is assumed to be a valid domain name,
		// those labels can still be invalid on their own.
		sIdx := i * addrStep

		c := arpa[sIdx]
		lo := fromHexByte(c)
		if lo == 0xff {
			return netip.Addr{}, &RuneError{
				Kind: AddrKindARPA,
				Rune: rune(c),
			}
		}

		c = arpa[sIdx+2]
		hi := fromHexByte(c)
		if hi == 0xff {
			return netip.Addr{}, &RuneError{
				Kind: AddrKindARPA,
				Rune: rune(c),
			}
		}

		if arpa[sIdx+1] != '.' || arpa[sIdx+3] != '.' {
			return netip.Addr{}, ErrNotAReversedIP
		}

		ip[net.IPv6len-i-1] = hi<<4 | lo
	}

	return netip.AddrFrom16(ip), nil
}

// IPFromReversedAddr tries to convert a full reversed ARPA address to a normal
// IP address.  arpa can be domain name or an FQDN.
//
// Any error returned will have the underlying type of *AddrError.
func IPFromReversedAddr(arpa string) (addr netip.Addr, err error) {
	arpa = strings.TrimSuffix(arpa, ".")
	err = ValidateDomainName(arpa)
	if err != nil {
		replaceKind(err, AddrKindARPA)

		return netip.Addr{}, err
	}

	defer makeAddrError(&err, arpa, AddrKindARPA)

	// TODO(a.garipov): Add stringutil.HasSuffixFold and remove this.
	arpa = strings.ToLower(arpa)

	if strings.HasSuffix(arpa, arpaV4Suffix) {
		ipStr := arpa[:len(arpa)-len(arpaV4Suffix)]

		return ipv4FromReversed(ipStr)
	}

	if strings.HasSuffix(arpa, arpaV6Suffix) {
		if l := len(arpa); l != arpaV6MaxLen {
			return netip.Addr{}, &LengthError{
				Kind:    AddrKindARPA,
				Allowed: []int{arpaV6MaxLen},
				Length:  l,
			}
		}

		return ipv6FromReversed(arpa)
	}

	return netip.Addr{}, ErrNotAReversedIP
}

// IPToReversedAddr returns the reversed ARPA address of ip suitable for reverse
// DNS (PTR) record lookups.  This is a modified version of function ReverseAddr
// from package github.com/miekg/dns package that accepts an IP.
//
// Any error returned will have the underlying type of [*AddrError].
func IPToReversedAddr(ip net.IP) (arpa string, err error) {
	const dot = "."

	var l int
	var suffix string
	var writeByte func(val byte)
	b := &strings.Builder{}
	if ip4 := ip.To4(); ip4 != nil {
		l, suffix = arpaV4MaxLen, arpaV4Suffix[1:]
		ip = ip4
		writeByte = func(val byte) {
			stringutil.WriteToBuilder(b, strconv.Itoa(int(val)), dot)
		}
	} else if ip6 := ip.To16(); ip6 != nil {
		l, suffix = arpaV6MaxLen, arpaV6Suffix[1:]
		ip = ip6
		writeByte = func(val byte) {
			stringutil.WriteToBuilder(
				b,
				strconv.FormatUint(uint64(val&0x0f), 16),
				dot,
				strconv.FormatUint(uint64(val>>4), 16),
				dot,
			)
		}
	} else {
		return "", &AddrError{
			Kind: AddrKindIP,
			Addr: ip.String(),
		}
	}

	b.Grow(l)
	for i := len(ip) - 1; i >= 0; i-- {
		writeByte(ip[i])
	}

	stringutil.WriteToBuilder(b, suffix)

	return b.String(), nil
}

// ipv4NetFromReversed parses an IPv4 reverse network.  It assumes that arpa is
// a valid domain name and is not a domain name with a full IPv4 address.
func ipv4NetFromReversed(arpa string) (pref netip.Prefix, err error) {
	var octet64 uint64
	var octetIdx int

	ip := [net.IPv4len]byte{}
	l := 0
	for addr := arpa; addr != ""; addr = addr[:octetIdx-1] {
		octetIdx = strings.LastIndexByte(addr, '.') + 1

		// Don't check for out of range since the domain is validated to have no
		// empty labels.
		octet64, err = strconv.ParseUint(addr[octetIdx:], 10, 8)
		if err != nil {
			// Don't wrap the error, since it's already informative enough as is.
			return netip.Prefix{}, err
		} else if octet64 != 0 && addr[octetIdx] == '0' {
			// Octets of an ARPA domain name shouldn't contain leading zero
			// except an octet itself equals zero.
			//
			// See RFC 1035 Section 3.5.
			return netip.Prefix{}, &AddrError{
				Err:  errors.Error("leading zero is forbidden at this position"),
				Kind: LabelKindDomain,
				Addr: addr[octetIdx:],
			}
		}

		ip[l] = byte(octet64)
		l++

		if octetIdx == 0 {
			// Prevent slicing with negative indices.
			break
		}
	}

	addr := netip.AddrFrom4(ip)

	return netip.PrefixFrom(addr, l*8), nil
}

// ipv6NetFromReversed parses an IPv6 reverse network.  It assumes that arpa is
// a valid domain name and is not a domain name with a full IPv6 address.
func ipv6NetFromReversed(arpa string) (pref netip.Prefix, err error) {
	const nibbleLen = len("0.")

	nibbleIdx := len(arpa) - len(arpaV6Suffix) + len(".") - nibbleLen
	if nibbleIdx%2 != 0 {
		return netip.Prefix{}, ErrNotAReversedSubnet
	}

	var b byte
	ip := [net.IPv6len]byte{}
	l := 0
	for ; nibbleIdx >= 0; nibbleIdx -= nibbleLen {
		if arpa[nibbleIdx+1] != '.' {
			return netip.Prefix{}, ErrNotAReversedSubnet
		}

		c := arpa[nibbleIdx]
		b = fromHexByte(c)
		if b == 0xff {
			return netip.Prefix{}, &RuneError{
				Kind: AddrKindARPA,
				Rune: rune(c),
			}
		}

		if l%2 == 0 {
			// An even digit stands for higher nibble of a byte.
			ip[l/2] |= b << 4
		} else {
			// An odd digit stands for lower nibble of a byte.
			ip[l/2] |= b
		}
		l++
	}

	addr := netip.AddrFrom16(ip)

	return netip.PrefixFrom(addr, l*4), nil
}

// subnetFromReversedV4 tries to convert arpa to a CIDR prefix.  It expects
// arpa to be a valid domain name in a lower case.
func subnetFromReversedV4(arpa string) (subnet netip.Prefix, rErr error) {
	arpa = arpa[:len(arpa)-len(arpaV4Suffix)]

	if dots := strings.Count(arpa, "."); dots > 3 {
		return netip.Prefix{}, ErrNotAReversedSubnet
	} else if dots == 3 {
		addr, err := ipv4FromReversed(arpa)
		if err != nil {
			// Don't wrap the error, since it's already informative enough as is.
			return netip.Prefix{}, err
		}

		return netip.PrefixFrom(addr, addr.BitLen()), nil
	}

	return ipv4NetFromReversed(arpa)
}

// subnetFromReversedV6 tries to convert arpa a CIDR prefix.  It expects
// arpa to be a valid domain name in a lower case.
func subnetFromReversedV6(arpa string) (subnet netip.Prefix, rErr error) {
	if l := len(arpa); l == arpaV6MaxLen {
		addr, err := ipv6FromReversed(arpa)
		if err != nil {
			// Don't wrap the error, since it's already informative enough as is.
			return netip.Prefix{}, err
		}

		return netip.PrefixFrom(addr, addr.BitLen()), nil
	} else if l > arpaV6MaxLen {
		return netip.Prefix{}, &LengthError{
			Kind:   AddrKindARPA,
			Max:    arpaV6MaxLen,
			Length: l,
		}
	}

	return ipv6NetFromReversed(arpa)
}

// PrefixFromReversedAddr tries to convert a reversed ARPA address to an IP
// address prefix.  arpa can be domain name or an FQDN.
//
// Any error returned will have the underlying type of [*AddrError].
func PrefixFromReversedAddr(arpa string) (p netip.Prefix, err error) {
	arpa = strings.TrimSuffix(arpa, ".")
	err = ValidateDomainName(arpa)
	if err != nil {
		replaceKind(err, AddrKindARPA)

		return netip.Prefix{}, err
	}

	defer makeAddrError(&err, arpa, AddrKindARPA)

	// TODO(a.garipov): Add stringutil.HasSuffixFold and remove this.
	arpa = strings.ToLower(arpa)

	if strings.HasSuffix(arpa, arpaV4Suffix) {
		return subnetFromReversedV4(arpa)
	}

	if strings.HasSuffix(arpa, arpaV6Suffix) {
		return subnetFromReversedV6(arpa)
	}

	return netip.Prefix{}, ErrNotAReversedSubnet
}
