package errors

// Error is the constant error type.
//
// See https://dave.cheney.net/2016/04/07/constant-errors.
type Error string

// type check
var _ error = Error("")

// Error implements the error interface for Error.
func (err Error) Error() (msg string) {
	return string(err)
}

// ErrBadEnumValue indicates that the provided value is not a valid value within
// an enumeration of types (a sum type) or values.
//
// For a value that is outside of a range of ordered values, use
// [ErrOutOfRange].
const ErrBadEnumValue Error = "bad enum value"

// ErrEmptyValue indicates that a value is provided but it is empty.  For
// example, a non-null but empty JSON or YAML object.
//
// For an absent value, use [ErrNoValue].
const ErrEmptyValue Error = "empty value"

// ErrNegative indicates that the provided value is negative when it should be
// greater than or equal to zero.
//
// For a value that should be greater than zero, use [ErrNotPositive].
const ErrNegative Error = "negative value"

// ErrNoValue indicates that a required value has not been provided.  For
// example, a null instead of an object in a JSON or YAML document.
//
// For a value that is present but empty, use [ErrEmptyValue].
const ErrNoValue Error = "no value"

// ErrNotPositive indicates that the provided value is negative or zero when it
// should be greater than zero.
//
// For a value that should be greater than or equal to zero, use [ErrNegative].
const ErrNotPositive Error = "not positive"

// ErrOutOfRange indicates that provided value is outside of a valid range of
// ordered values.
//
// For a value that is not a valid enum or sum type value, use [ErrBadEnumValue].
const ErrOutOfRange Error = "out of range"
