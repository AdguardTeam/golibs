package slogutil

// Format represents an acceptable format of logs.
type Format string

// Valid formats.
const (
	FormatAdGuardLegacy Format = "adguard_legacy"
	FormatDefault       Format = "default"
	FormatJSON          Format = "json"
	FormatJSONHybrid    Format = "jsonhybrid"
	FormatText          Format = "text"
)

// NewFormat returns a new valid format.
func NewFormat(s string) (f Format, err error) {
	switch f = Format(s); f {
	case
		FormatAdGuardLegacy,
		FormatDefault,
		FormatJSON,
		FormatJSONHybrid,
		FormatText:
		return f, nil
	default:
		return "", &BadFormatError{
			Format: s,
		}
	}
}
