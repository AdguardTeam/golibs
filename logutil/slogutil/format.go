package slogutil

// Format represents an acceptable format of logs.
type Format string

// Valid formats.
const (
	FormatAdGuardLegacy = "adguard_legacy"
	FormatDefault       = "default"
	FormatJSON          = "json"
	FormatJSONHybrid    = "jsonhybrid"
	FormatText          = "text"
)

// NewFormat returns a new valid format.
func NewFormat(s string) (f Format, err error) {
	switch s {
	case
		FormatAdGuardLegacy,
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
