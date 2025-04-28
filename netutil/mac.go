package netutil

import (
	"net"
)

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
	var sep byte
	var fragNum, fragLen int
	l := len(s)
	switch l {
	case len("00:00:5e:00:53:01"):
		fragLen, fragNum = 2, 6
		sep = s[2]
	case len("02:00:5e:10:00:00:00:01"):
		fragLen, fragNum = 2, 8
		sep = s[2]
	case len("00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01"):
		fragLen, fragNum = 2, 20
		sep = s[2]
	case len("0000.5e00.5301"):
		fragLen, fragNum = 4, 3
		sep = '.'
	case len("0200.5e10.0000.0001"):
		fragLen, fragNum = 4, 4
		sep = '.'
	case len("0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001"):
		fragLen, fragNum = 4, 10
		sep = '.'
	default:
		return false
	}

	switch sep {
	case '.', '-', ':':
		// Go on.
	default:
		return false
	}

	return isValidHexSepString(s, fragLen, fragNum, sep)
}

// isValidHexSepString returns true if s is a string containing fragNum
// hexadecimal fragments, each with a length of fragLen, and separated by sep.
// s must be of the appropriate length.
func isValidHexSepString(s string, fragLen, fragNum int, sep byte) (ok bool) {
	idx := 0

	for range fragNum - 1 {
		if !isValidHexString(s[idx : idx+fragLen]) {
			return false
		}

		idx += fragLen

		if s[idx] != sep {
			return false
		}

		idx++
	}

	return isValidHexString(s[idx : idx+fragLen])
}

// isValidHexString returns true if s is a valid hexadecimal string.
func isValidHexString(s string) (ok bool) {
	for i := range len(s) {
		if fromHexByte(s[i]) == 0xff {
			return false
		}
	}

	return true
}
