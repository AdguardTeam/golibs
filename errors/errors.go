// Package errors is a drop-in replacement and extension of the Go standard
// library's package [errors].
package errors

import (
	stderrors "errors"
)

// NOTE:  Please keep only the stdlib-related identifiers in this file.

// ErrUnsupported indicates that a requested operation cannot be performed,
// because it is unsupported.  For example, a call to [os.Link] when using a
// file system that does not support hard links.
var ErrUnsupported = stderrors.ErrUnsupported

// Wrapper is a copy of the hidden wrapper interface from the Go standard
// library.  It is added here for tests, linting, etc.
type Wrapper interface {
	Unwrap() error
}

// WrapperSlice is a copy of the hidden wrapper interface added to the Go
// standard library in Go 1.20.  It is added here for tests, linting, etc.
type WrapperSlice interface {
	Unwrap() []error
}

// Join returns an error that wraps the given errors.  Any nil error values are
// discarded.  Join returns nil if errs contains no non-nil values.  The error
// formats as the concatenation of the strings obtained by calling the Error
// method of each element of errs, with a newline between each string.
//
// It calls [stderrors.Join] from the Go standard library.
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.  Otherwise, it returns false.
//
// It calls [stderrors.As] from the Go standard library.
func As(err error, target any) (ok bool) {
	return stderrors.As(err, target)
}

// Aser is a copy of the hidden aser interface from the Go standard library.  It
// is added here for tests, linting, etc.
type Aser interface {
	As(target any) (ok bool)
}

// Is reports whether any error in err's chain matches target.
//
// It calls [stderrors.Is] from the Go standard library.
func Is(err, target error) (ok bool) {
	return stderrors.Is(err, target)
}

// Iser is a copy of the hidden iser interface from the Go standard library.  It
// is added here for tests, linting, etc.
type Iser interface {
	Is(target error) (ok bool)
}

// New returns an error that formats as the given msg.  Each call to New returns
// a distinct error value even if the text is identical.
//
// It calls [stderrors.New] from the Go standard library.
//
// Deprecated: Use type [Error] and constant errors instead.
func New(msg string) (err error) {
	return stderrors.New(msg)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type
// contains an Unwrap method returning error.  Otherwise, Unwrap returns nil.
//
// It calls [stderrors.Unwrap] from the Go standard library.
func Unwrap(err error) (unwrapped error) {
	return stderrors.Unwrap(err)
}
