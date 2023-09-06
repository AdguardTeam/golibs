package hostsfile

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

// ErrEmptyLine is returned when the hosts file line is empty or contains only
// comments and spaces.
const ErrEmptyLine errors.Error = "line is empty"

// ErrNoHosts is returned when the record doesn't contain any delimiters, but
// the IP address is valid.
const ErrNoHosts errors.Error = "no hostnames"

// LineError is an error about a specific line in the hosts file.
type LineError struct {
	// err is the original error.
	err error

	// Line is the line number in the hosts file source.
	Line int
}

// type check
var _ error = (*LineError)(nil)

// Error implements the [error] interface for *LineErr.
func (e *LineError) Error() (msg string) {
	return fmt.Sprintf("line %d: %s", e.Line, e.err)
}

// type check
var _ errors.Wrapper = (*LineError)(nil)

// Unwrap implements the [errors.Wrapper] interface for *LineErr.
func (e *LineError) Unwrap() (unwrapped error) { return e.err }
