package netutil

import (
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/errors"
)

const (
	// ErrNotAReversedIP is the underlying error returned from validation
	// functions when a domain name is not a full reversed IP address.
	ErrNotAReversedIP errors.Error = "not a full reversed ip address"
)

// BadDomainError is the underlying type of errors returned from validation
// functions when a domain name is invalid.
type BadDomainError struct {
	// Err is the underlying error.  The type of the underlying error is
	// one of the following:
	//
	//   *BadIPError
	//   *BadIPv4Error
	//   *BadLabelError
	//   *BadLengthError
	//   *BadRuneError
	//   *EmptyError
	//   *TooLongError
	//   any error returned by the Punicode validation.
	//
	// It can also be ErrNotAReversedIP.
	Err error
	// Kind is either "arpa domain name" or "domain name".
	Kind string
	// Name is the text of the invalid domain name.
	Name string
}

// Error implements the error interface for *BadDomainError.
func (err *BadDomainError) Error() (msg string) {
	return fmt.Sprintf("bad %s %q: %s", err.Kind, err.Name, err.Err)
}

// Unwrap implements the errors.Wrapper interface for *BadDomainError.  It
// returns err.Err.
func (err *BadDomainError) Unwrap() (unwrapped error) {
	return err.Err
}

// BadIPError is the underlying type of errors returned from validation
// functions when an IP address is invalid.
type BadIPError struct {
	// IP is the text of the invalid IP address.
	IP string
}

// Error implements the error interface for *BadIPError.
func (err *BadIPError) Error() (msg string) {
	return fmt.Sprintf("bad ip address %q", err.IP)
}

// BadIPv4Error is the underlying type of errors returned from validation
// functions when an IP address is not a valid IPv4 address.
type BadIPv4Error struct {
	// IP is the text of the invalid IP address.
	IP string
}

// Error implements the error interface for *BadIPv4Error.
func (err *BadIPv4Error) Error() (msg string) {
	return fmt.Sprintf("bad ipv4 address %q", err.IP)
}

// BadLabelError is the underlying type of errors returned from validation
// functions when a domain name label is invalid.
type BadLabelError struct {
	// Err is the underlying error.  The type of the underlying error is
	// either *BadRuneError, or *EmptyError, or *TooLongError.
	Err error
	// Label is the text of the label.
	Label string
}

// Error implements the error interface for *BadLabelError.
func (err *BadLabelError) Error() (msg string) {
	return fmt.Sprintf("bad domain name label %q: %s", err.Label, err.Err)
}

// Unwrap implements the errors.Wrapper interface for *BadLabelError.  It
// returns err.Err.
func (err *BadLabelError) Unwrap() (unwrapped error) {
	return err.Err
}

// BadLengthError is the underlying type of errors returned from validation
// functions when an address or a part of an address has a bad length.
type BadLengthError struct {
	// Kind is either "arpa domain name" or "mac address".
	Kind string
	// Allowed are the allowed lengths for this kind of address.
	Allowed []int
	// Length is the length of the provided address.
	Length int
}

// Error implements the error interface for *BadLengthError.
func (err *BadLengthError) Error() (msg string) {
	return fmt.Sprintf("bad %s length %d, allowed: %v", err.Kind, err.Length, err.Allowed)
}

// BadMACError is the underlying type of errors returned from validation
// functions when a MAC address is invalid.
type BadMACError struct {
	// Err is the underlying error.  The type of the underlying error is
	// either *EmptyError, or *BadLengthError.
	Err error
	// MAC is the text of the MAC address.
	MAC net.HardwareAddr
}

// Error implements the error interface for *BadMACError.
func (err *BadMACError) Error() (msg string) {
	return fmt.Sprintf("bad mac address %q: %s", err.MAC, err.Err)
}

// Unwrap implements the errors.Wrapper interface for *BadMACError.  It
// returns err.Err.
func (err *BadMACError) Unwrap() (unwrapped error) {
	return err.Err
}

// BadRuneError is the underlying type of errors returned from validation
// functions when a rune in the address is invalid.
type BadRuneError struct {
	// Kind is either "arpa domain name", or "domain name label", or "mac
	// address".
	Kind string
	// Rune is the invalid rune.
	Rune rune
}

// Error implements the error interface for *BadRuneError.
func (err *BadRuneError) Error() (msg string) {
	return fmt.Sprintf("bad %s rune %q", err.Kind, err.Rune)
}

// EmptyError is the underlying type of errors returned from validation
// functions when an address or a part of an address is missing.
type EmptyError struct {
	// Kind is either "domain name", or "domain name label", or "mac
	// address".
	Kind string
}

// Error implements the error interface for *EmptyError.
func (err *EmptyError) Error() (msg string) {
	return fmt.Sprintf("%s is empty", err.Kind)
}

// TooLongError is the underlying type of errors returned from validation
// functions when an address or a part of an address is too long.
type TooLongError struct {
	// Kind is either "domain name", or "domain name label", or "mac
	// address".
	Kind string
	// Max is the maximum length for this part or address kind.
	Max int
}

// Error implements the error interface for *TooLongError.
func (err *TooLongError) Error() (msg string) {
	return fmt.Sprintf("%s is too long, max: %d", err.Kind, err.Max)
}
