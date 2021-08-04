// Package netutil contains common utilities for IP, MAC, and other kinds of
// network addresses.
//
// TODO(a.garipov): Add more examples.
//
// TODO(a.garipov): Add HostPort and IPPort structs with decoding and encoding,
// fmt.Srtinger implementations, etc.
package netutil

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/idna"
)

// CloneIP returns a clone of an IP address that doesn't share the same
// underlying array with it.
func CloneIP(ip net.IP) (clone net.IP) {
	if ip != nil && len(ip) == 0 {
		return net.IP{}
	}

	return append(clone, ip...)
}

// CloneIPs returns a deep clone of ips.
func CloneIPs(ips []net.IP) (clone []net.IP) {
	if ips == nil {
		return nil
	}

	clone = make([]net.IP, len(ips))
	for i, ip := range ips {
		clone[i] = CloneIP(ip)
	}

	return clone
}

// CloneMAC returns a clone of a MAC address.
func CloneMAC(mac net.HardwareAddr) (clone net.HardwareAddr) {
	if mac != nil && len(mac) == 0 {
		return net.HardwareAddr{}
	}

	return append(clone, mac...)
}

// CloneURL returns a deep clone of u.  The User pointer of clone is the same,
// since a *url.Userinfo is effectively an immutable value.
func CloneURL(u *url.URL) (clone *url.URL) {
	if u == nil {
		return nil
	}

	cloneVal := *u

	return &cloneVal
}

// IPPortFromAddr returns the IP address and the port from addr.  If addr is
// neither a *net.TCPAddr nor a *net.UDPAddr, it returns nil and 0.
func IPPortFromAddr(addr net.Addr) (ip net.IP, port int) {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return addr.IP, addr.Port
	case *net.UDPAddr:
		return addr.IP, addr.Port
	}

	return nil, 0
}

// IsValidHostInnerRune returns true if r is a valid inner—that is, neither
// initial nor final—rune for a hostname label.
func IsValidHostInnerRune(r rune) (ok bool) {
	return r == '-' || IsValidHostOuterRune(r)
}

// IsValidHostOuterRune returns true if r is a valid initial or final rune for
// a hostname label.
func IsValidHostOuterRune(r rune) (ok bool) {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

// JoinHostPort is a convinient wrapper for net.JoinHostPort with port of type
// int.
func JoinHostPort(host string, port int) (hostport string) {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

// ParseIP is a wrapper around net.ParseIP that returns a useful error.
//
// Any error returned will have the underlying type of *BadIPError.
func ParseIP(s string) (ip net.IP, err error) {
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, &BadIPError{IP: s}
	}

	return ip, nil
}

// ParseIPv4 is a wrapper around net.ParseIP that makes sure that the parsed IP
// is an IPv4 address and returns a useful error.
//
// Any error returned will have the underlying type of either *BadIPError or
// *BadIPv4Error,
func ParseIPv4(s string) (ip net.IP, err error) {
	ip, err = ParseIP(s)
	if err != nil {
		return nil, err
	}

	if ip = ip.To4(); ip == nil {
		return nil, &BadIPv4Error{IP: s}
	}

	return ip, nil
}

// SplitHostPort is a convenient wrapper for net.SplitHostPort with port of type
// int.
func SplitHostPort(hostport string) (host string, port int, err error) {
	var portStr string
	host, portStr, err = net.SplitHostPort(hostport)
	if err != nil {
		return "", 0, err
	}

	port, err = strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("parsing port: %w", err)
	}

	return host, port, nil
}

// SplitHost is a wrapper for net.SplitHostPort for cases when the hostport may
// or may not contain a port.
func SplitHost(hostport string) (host string, err error) {
	host, _, err = net.SplitHostPort(hostport)
	if err != nil {
		// Check for the missing port error.  If it is that error, just
		// use the host as is.
		//
		// See the source code for net.SplitHostPort.
		const missingPort = "missing port in address"

		addrErr := &net.AddrError{}
		if !errors.As(err, &addrErr) || addrErr.Err != missingPort {
			return "", err
		}

		host = hostport
	}

	return host, nil
}

// ValidateMAC returns an error if hwa is not a valid EUI-48, EUI-64, or
// 20-octet InfiniBand link-layer address.
//
// Any error returned will have the underlying type of *BadMACError.
func ValidateMAC(mac net.HardwareAddr) (err error) {
	defer func() {
		if err != nil {
			err = &BadMACError{
				Err: err,
				MAC: mac,
			}
		}
	}()

	const kind = "mac address"

	switch l := len(mac); l {
	case 0:
		return &EmptyError{Kind: kind}
	case 6, 8, 20:
		return nil
	default:
		return &BadLengthError{
			Kind:    kind,
			Allowed: []int{6, 8, 20},
			Length:  l,
		}
	}
}

// MaxDomainLabelLen is the maximum allowed length of a domain name label
// according to RFC 1035.
const MaxDomainLabelLen = 63

// MaxDomainNameLen is the maximum allowed length of a full domain name
// according to RFC 1035.
//
// See also: https://stackoverflow.com/a/32294443/1892060.
const MaxDomainNameLen = 253

// ValidateDomainNameLabel returns an error if label is not a valid label of
// a domain name.  An empty label is considered invalid.
//
// Any error returned will have the underlying type of *BadLabelError.
func ValidateDomainNameLabel(label string) (err error) {
	defer func() {
		if err != nil {
			err = &BadLabelError{Err: err, Label: label}
		}
	}()

	const kind = "domain name label"

	l := len(label)
	if l == 0 {
		return &EmptyError{Kind: kind}
	} else if l > MaxDomainLabelLen {
		return &TooLongError{Kind: kind, Max: MaxDomainLabelLen}
	}

	if r := rune(label[0]); !IsValidHostOuterRune(r) {
		return &BadRuneError{
			Kind: kind,
			Rune: r,
		}
	} else if l == 1 {
		return nil
	}

	for _, r := range label[1 : l-1] {
		if !IsValidHostInnerRune(r) {
			return &BadRuneError{
				Kind: kind,
				Rune: r,
			}
		}
	}

	if r := rune(label[l-1]); !IsValidHostOuterRune(r) {
		return &BadRuneError{
			Kind: kind,
			Rune: r,
		}
	}

	return nil
}

// ValidateDomainName validates the domain name in accordance to RFC 952, RFC
// 1035, and with RFC-1123's inclusion of digits at the start of the host.  It
// doesn't validate against two or more hyphens to allow punycode and
// internationalized domains.
//
// Any error returned will have the underlying type of *BadDomainError.
func ValidateDomainName(name string) (err error) {
	const kind = "domain name"

	defer func() {
		if err != nil {
			err = &BadDomainError{
				Err:  err,
				Kind: kind,
				Name: name,
			}
		}
	}()

	name, err = idna.ToASCII(name)
	if err != nil {
		return err
	}

	l := len(name)
	if l == 0 {
		return &EmptyError{Kind: kind}
	} else if l > MaxDomainNameLen {
		return &TooLongError{Kind: kind, Max: MaxDomainNameLen}
	}

	labels := strings.Split(name, ".")
	for _, l := range labels {
		err = ValidateDomainNameLabel(l)
		if err != nil {
			return err
		}
	}

	return nil
}

// fromHexByte converts a single hexadecimal ASCII digit character into an
// integer from 0 to 15.  For all other characters it returns 0xff.
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
//   49.91.20.104.in-addr.arpa
//
// An example of IPv6 with a maximum length:
//
//   1.3.b.5.4.1.8.6.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.0.0.7.4.6.0.6.2.ip6.arpa
//
const (
	arpaV4MaxIPLen = len("000.000.000.000")
	arpaV6MaxIPLen = len("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0")

	arpaV4MaxLen = arpaV4MaxIPLen + len(arpaV4Suffix)
	arpaV6MaxLen = arpaV6MaxIPLen + len(arpaV6Suffix)
)

// reverseIP inverts the order of bytes in an IP address in-place.
func reverseIP(ip net.IP) {
	l := len(ip)
	for i := range ip[:l/2] {
		ip[i], ip[l-i-1] = ip[l-i-1], ip[i]
	}
}

// ipv6FromReversedAddr parses an IPv6 reverse address.  It assumes that arpa is
// a valid domain name.
func ipv6FromReversedAddr(arpa string) (ip net.IP, err error) {
	const kind = "arpa domain name"

	ip = make(net.IP, net.IPv6len)

	const addrStep = len("0.0.")
	for i := range ip {
		// Get the two half-byte and merge them together.  Validate the
		// dots between them since while arpa is assumed to be a valid
		// domain name, those labels can still be invalid on their own.
		sIdx := i * addrStep

		c := arpa[sIdx]
		lo := fromHexByte(c)
		if lo == 0xff {
			return nil, &BadRuneError{
				Kind: kind,
				Rune: rune(c),
			}
		}

		c = arpa[sIdx+2]
		hi := fromHexByte(c)
		if hi == 0xff {
			return nil, &BadRuneError{
				Kind: kind,
				Rune: rune(c),
			}
		}

		if arpa[sIdx+1] != '.' || arpa[sIdx+3] != '.' {
			return nil, ErrNotAReversedIP
		}

		ip[net.IPv6len-i-1] = hi<<4 | lo
	}

	return ip, nil
}

// IPFromReversedAddr tries to convert a full reversed ARPA address to a normal
// IP address.  arpa can be domain name or an FQDN.
//
// Any error returned will have the underlying type of *BadDomainError.
func IPFromReversedAddr(arpa string) (ip net.IP, err error) {
	const kind = "arpa domain name"

	arpa = strings.TrimSuffix(arpa, ".")
	err = ValidateDomainName(arpa)
	if err != nil {
		bdErr := err.(*BadDomainError)
		bdErr.Kind = kind

		return nil, bdErr
	}

	defer func() {
		if err != nil {
			err = &BadDomainError{
				Err:  err,
				Kind: kind,
				Name: arpa,
			}
		}
	}()

	// TODO(a.garipov): Add stringutil.HasSuffixFold and remove this.
	arpa = strings.ToLower(arpa)

	if strings.HasSuffix(arpa, arpaV4Suffix) {
		ipStr := arpa[:len(arpa)-len(arpaV4Suffix)]
		ip, err = ParseIPv4(ipStr)
		if err != nil {
			return nil, err
		}

		reverseIP(ip)

		return ip, nil
	}

	if strings.HasSuffix(arpa, arpaV6Suffix) {
		if l := len(arpa); l != arpaV6MaxLen {
			return nil, &BadLengthError{
				Kind:    kind,
				Allowed: []int{arpaV6MaxLen},
				Length:  l,
			}
		}

		ip, err = ipv6FromReversedAddr(arpa)
		if err != nil {
			return nil, err
		}

		return ip, nil
	}

	return nil, ErrNotAReversedIP
}
