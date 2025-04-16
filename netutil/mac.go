package netutil

import "net"

// ValidateMAC returns an error if mac is not a valid EUI-48, EUI-64, or
// 20-octet InfiniBand link-layer address.
//
// Any error returned will have the underlying type of [*AddrError].
func ValidateMAC(mac net.HardwareAddr) (err error) {
	defer makeAddrError(&err, mac.String(), AddrKindMAC)

	switch l := len(mac); l {
	case 0:
		return &LengthError{
			Kind:   AddrKindMAC,
			Length: 0,
		}
	case 6, 8, 20:
		return nil
	default:
		return &LengthError{
			Kind:    AddrKindMAC,
			Allowed: []int{6, 8, 20},
			Length:  l,
		}
	}
}

// IsValidMACString is a best-effort check to determine if s is a valid MAC
// address before using [net.ParseMAC], aimed at reducing allocations.
func IsValidMACString(s string) (ok bool) {
	if len(s) < 14 {
		return false
	}

	if s[2] == ':' || s[2] == '-' {
		return hasValidMACFragments(s)
	} else if s[4] == '.' {
		return hasValidMACLongFragments(s)
	}

	return false
}

// hasValidMACLongFragments returns true if s contains valid MAC address
// hex digits.
func hasValidMACFragments(s string) (ok bool) {
	if (len(s)+1)%3 != 0 {
		return false
	}

	n := (len(s) + 1) / 3

	switch n {
	case 6, 8, 20:
		// Go on.
	default:
		return false
	}

	for x, i := 0, 0; i < n; i++ {
		if !startsWith2Hex(s[x:], s[2]) {
			return false
		}

		x += 3
	}

	return true
}

// hasValidMACLongFragments returns true if s contains valid long MAC address
// hex fragments.
func hasValidMACLongFragments(s string) (ok bool) {
	if (len(s)+1)%5 != 0 {
		return false
	}

	n := 2 * (len(s) + 1) / 5
	switch n {
	case 6, 8, 20:
		// Go on.
	default:
		return false
	}

	for x, i := 0, 0; i < n; i += 2 {
		if !startsWith2Hex(s[x:x+2], 0) {
			return false
		}

		if !startsWith2Hex(s[x+2:], s[4]) {
			return false
		}

		x += 5
	}

	return true
}

// isValidHexChar returns true if c is a valid hexadecimal digit character.
func isValidHexChar(c byte) (ok bool) {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// startsWith2Hex checks if the next two hex digits of s could be converted into
// a byte.  If s is longer than 2 bytes then the third byte must be e.  If the
// first two bytes of s are not hex digits or the third byte does not match e,
// false is returned.
func startsWith2Hex(s string, e byte) (ok bool) {
	if s == "" || len(s) > 2 && s[2] != e {
		return false
	}

	for i := range 2 {
		if !isValidHexChar(s[i]) {
			return false
		}
	}

	return true
}
