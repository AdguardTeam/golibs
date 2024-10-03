package netutil

import (
	"math"
	"net"
	"slices"
	"strings"
)

// IP Address Constants And Utilities

// Bit lengths of IP addresses.
const (
	IPv4BitLen = net.IPv4len * 8
	IPv6BitLen = net.IPv6len * 8
)

// CloneIPs returns a deep clone of ips.
func CloneIPs(ips []net.IP) (clone []net.IP) {
	if ips == nil {
		return nil
	}

	clone = make([]net.IP, len(ips))
	for i, ip := range ips {
		clone[i] = slices.Clone(ip)
	}

	return clone
}

// IPAndPortFromAddr returns the IP address and the port from addr.  If addr is
// neither a [*net.TCPAddr] nor a [*net.UDPAddr], it returns nil and 0.  The
// port of the address should fit in uint16.
func IPAndPortFromAddr(addr net.Addr) (ip net.IP, port uint16) {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		// #nosec G115 -- Assume that ports always fit in uint16.
		return addr.IP, uint16(addr.Port)
	case *net.UDPAddr:
		// #nosec G115 -- Assume that ports always fit in uint16.
		return addr.IP, uint16(addr.Port)
	}

	return nil, 0
}

// IPv4bcast returns a new limited broadcast IPv4 address, 255.255.255.255.  It
// has the same name as the variable in package net, but the result always has
// four bytes.
func IPv4bcast() (ip net.IP) { return net.IP{255, 255, 255, 255} }

// IPv4allsys returns a new all systems (aka all hosts) IPv4 address, 224.0.0.1.
// It has the same name as the variable in package net, but the result always
// has four bytes.
func IPv4allsys() (ip net.IP) { return net.IP{224, 0, 0, 1} }

// IPv4allrouter returns a new all routers IPv4 address, 224.0.0.2.  It has the
// same name as the variable in package net, but the result always has four
// bytes.
func IPv4allrouter() (ip net.IP) { return net.IP{224, 0, 0, 2} }

// IPv4Zero returns a new unspecified (aka empty or null) IPv4 address, 0.0.0.0.
// It has the same name as the variable in package net, but the result always
// has four bytes.
func IPv4Zero() (ip net.IP) { return net.IP{0, 0, 0, 0} }

// IPv6Zero returns a new unspecified (aka empty or null) IPv6 address, [::].
// It has the same name as the variable in package net.
func IPv6Zero() (ip net.IP) {
	return net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

// ParseIP is a wrapper around net.ParseIP that returns a useful error.
//
// Any error returned will have the underlying type of [*AddrError].
func ParseIP(s string) (ip net.IP, err error) {
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, &AddrError{
			Kind: AddrKindIP,
			Addr: s,
		}
	}

	return ip, nil
}

// ParseIPv4 is a wrapper around net.ParseIP that makes sure that the parsed IP
// is an IPv4 address and returns a useful error.
//
// Any error returned will have the underlying type of either [*AddrError].
func ParseIPv4(s string) (ip net.IP, err error) {
	ip, err = ParseIP(s)
	if err != nil {
		err.(*AddrError).Kind = AddrKindIPv4

		return nil, err
	}

	if ip = ip.To4(); ip == nil {
		return nil, &AddrError{
			Kind: AddrKindIPv4,
			Addr: s,
		}
	}

	return ip, nil
}

// ValidateIP returns an error if ip is not a valid IPv4 or IPv6 address.
//
// Any error returned will have the underlying type of [*AddrError].
func ValidateIP(ip net.IP) (err error) {
	// TODO(a.garipov):  Get rid of unnecessary allocations in case of valid IP.
	defer makeAddrError(&err, ip.String(), AddrKindIP)

	switch l := len(ip); l {
	case 0:
		return &LengthError{
			Kind:   AddrKindIP,
			Length: 0,
		}
	case net.IPv4len, net.IPv6len:
		return nil
	default:
		return &LengthError{
			Kind:    AddrKindIP,
			Allowed: []int{net.IPv4len, net.IPv6len},
			Length:  l,
		}
	}
}

// IsValidIPString returns true if s is a valid IPv4 or IPv6 address string
// representation as accepted by [netip.ParseAddr].
func IsValidIPString(s string) (ok bool) {
	// maxSignificant is the maximum number of significant characters in a field
	// of an IPv6 address.  There is no need to search for separator longer than
	// this number of bytes.
	//
	// See https://www.rfc-editor.org/rfc/rfc4291#section-2.2.
	const maxSignificant = 4

	strLen := len(s)
	for i, significant := 0, 0; i < strLen && significant <= maxSignificant; i++ {
		switch s[i] {
		case '.':
			return isValidIPv4String(s)
		case ':':
			withoutZone, zone, hasZone := strings.Cut(s, "%")
			if hasZone && zone == "" {
				// Zone cannot be empty.
				return false
			}

			return isValidIPv6String(withoutZone)
		default:
			significant++
		}
	}

	return false
}

// isValidIPv4String returns true if s is a valid IPv4 address string
// representation in dotted decimal form.
func isValidIPv4String(s string) (ok bool) {
	num := 1
	label, s, ok := strings.Cut(s, ".")
	for ; num < net.IPv4len && ok; label, s, ok = strings.Cut(s, ".") {
		if !isIPv4Label(label) {
			return false
		}

		num++
	}

	return num == net.IPv4len && !ok && isIPv4Label(label)
}

// isIPv4Label reports whether label is a valid label for an IPv4 address, i.e.
// a decimal number in the range [0, 255] with no leading zeros.
func isIPv4Label(label string) (ok bool) {
	switch l := len(label); {
	case l < 1, l > 3:
		return false
	case l == 1:
		return label[0] >= '0' && label[0] <= '9'
	case label[0] == '0':
		return false
	default:
		val := 0
		for _, c := range label {
			if c < '0' || c > '9' {
				return false
			}

			val = val*10 + int(c-'0')
		}

		return val <= math.MaxUint8
	}
}

// maxIPv6FieldsNum is the maximum number of fields in an IPv6 address.
const maxIPv6FieldsNum = net.IPv6len / 2

// isValidIPv6String returns true if s is a valid IPv6 address string
// representation.  Note, that the address is expected to have no zone.
func isValidIPv6String(s string) (ok bool) {
	hasEllipsis := strings.HasPrefix(s, "::")
	if hasEllipsis {
		s = s[2:]
	}

	fieldsNum := 0
	for ; fieldsNum < maxIPv6FieldsNum && s != ""; fieldsNum++ {
		s, ok = trimValidIPv6Field(s, fieldsNum, hasEllipsis)
		if !ok {
			return false
		} else if s == "" {
			return true
		}

		var sepLen int
		sepLen, hasEllipsis = countIPv6SepRunes(s, hasEllipsis)
		if sepLen == 0 {
			return false
		}

		s = s[sepLen:]
	}

	return s == "" && (hasEllipsis == (fieldsNum < maxIPv6FieldsNum))
}

// countIPv6SepRunes returns the number of first runes in s that are a separator
// in an IPv6 address, considering whether an ellipsis has been found before.
// It returns 0 if there is no valid separator in the beginning of s.
func countIPv6SepRunes(s string, hadEllipsis bool) (n int, hasEllipsis bool) {
	switch {
	case
		// Not a valid character.
		s[0] != ':',
		// Colon at the end.
		len(s) == 1:
		return 0, hadEllipsis
	case s[1] == ':':
		// Ellipsis.
		if hadEllipsis {
			// Only one ellipsis is allowed.
			return 0, false
		}

		return 2, true
	default:
		return 1, hadEllipsis
	}
}

// trimValidIPv6Field trims the first field from s and returns the rest of the
// string, considering the number of fields got so far and whether there is an
// ellipsis has been found.  ok is false if the field is invalid, otherwise if
// withoutField is empty, the rest of s is considered a valid IPv6 address tail.
func trimValidIPv6Field(s string, gotFields int, hasEllipsis bool) (withoutField string, ok bool) {
	fieldLen := countIPv6FieldRunes(s)
	switch fieldLen {
	case 0:
		// No digits found, fail.
		return "", false
	case len(s):
		// The whole string is a field.
		gotFields++

		return "", hasEllipsis == (gotFields < maxIPv6FieldsNum)
	default:
		// Go on.
	}

	if s[fieldLen] == '.' {
		// Probably an IPv4 in the end.
		return "", hasEllipsis == (gotFields < maxIPv6FieldsNum-2) && isValidIPv4String(s)
	}

	return s[fieldLen:], true
}

// countIPv6FieldRunes returns the number of runes in the first field of an IPv6
// address.  It returns 0 if the field is invalid, due to an assumption that the
// field is not shorter than 1 rune.
func countIPv6FieldRunes(s string) (n int) {
	for n = range s {
		if fromHexByte(s[n]) == 0xff {
			// Not a hex digit, return.
			return n
		} else if n > 3 {
			// IPv6 label can't contain more than 4 hex digits.
			return 0
		}
	}

	return len(s)
}
