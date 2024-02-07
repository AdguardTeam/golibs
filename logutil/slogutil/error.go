package slogutil

import "fmt"

// BadFormatError is an error about a bad logging format.
type BadFormatError struct {
	Format string
}

// type check
var _ error = (*BadFormatError)(nil)

// Error implements the [error] interface for *BadFormatError.
func (err *BadFormatError) Error() (msg string) {
	return fmt.Sprintf("bad log format %q", err.Format)
}
