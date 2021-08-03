// Package netutil contains common utilities for IP, MAC, and other kinds of
// network addresses.
//
// TODO(a.garipov): Add examples.
package netutil

import (
	"fmt"
	"net"
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
	defer func() {
		if err != nil {
			err = &BadDomainError{Err: err, Name: name}
		}
	}()

	const kind = "domain name"

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
