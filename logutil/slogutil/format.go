//go:build go1.21

package slogutil

// Format represents an acceptable format of logs.
type Format string

// Valid formats.
const (
	FormatDefault    = "default"
	FormatJSON       = "json"
	FormatJSONHybrid = "jsonhybrid"
	FormatText       = "text"
)

// NewFormat returns a new valid format.
func NewFormat(s string) (f Format, err error) {
	switch s {
	case
		FormatDefault,
		FormatJSON,
		FormatJSONHybrid,
		FormatText:
		return Format(s), nil
	default:
		return "", &BadFormatError{
			Format: s,
		}
	}
}
