package timeutil

import (
	"bytes"
	"encoding"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/AdguardTeam/golibs/errors"
)

// ErrInvalidDuration is returned when an invalid duration string is provided.
const ErrInvalidDuration errors.Error = "invalid duration"

// durationUnitDay represents the [Day] duration unit.
const durationUnitDay = 'd'

// Duration is a helper type for time.Duration providing functionality for
// encoding.
type Duration time.Duration

// type check
var _ fmt.Stringer = Duration(0)

// String implements the [fmt.Stringer] interface for Duration.  It wraps
// [time.Duration.String] method and additionally cuts off non-leading zero
// values of minutes and seconds.  This method also supports day units.  Some
// values which are differ between the implementations:
//
//	Duration:    "1m", time.Duration:    "1m0s"
//	Duration:    "1h", time.Duration:  "1h0m0s"
//	Duration:  "1h1m", time.Duration:  "1h1m0s"
//	Duration: "1d12h", time.Duration: "36h0m0s"
func (d Duration) String() (str string) {
	timeDur := time.Duration(d)
	abs := timeDur
	var sign string
	if timeDur < 0 {
		abs = -timeDur
		sign = "-"
	}

	days := abs / Day
	if days == 0 {
		return stringWithZeroCut(timeDur)
	}

	remainder := abs - days*Day
	if remainder == 0 {
		return fmt.Sprintf("%s%dd", sign, days)
	}

	remainderStr := stringWithZeroCut(remainder)

	return fmt.Sprintf("%s%dd%s", sign, days, remainderStr)
}

// stringWithZeroCut wraps [time.Duration.String] method, and additionally cuts
// off non-leading zero values of minutes and seconds.
func stringWithZeroCut(d time.Duration) (str string) {
	str = d.String()

	const (
		tailMin    = len(`0s`)
		tailMinSec = len(`0m0s`)
	)

	const (
		secsInHour = time.Hour / time.Second
		minsInHour = time.Hour / time.Minute
	)

	switch rounded := d / time.Second; {
	case
		rounded == 0,
		rounded*time.Second != d,
		rounded%60 != 0:
		// Return the uncut value if it's either equal to zero or has fractions
		// of a second or even whole seconds in it.
		return str
	case (rounded%secsInHour)/minsInHour != 0:
		return str[:len(str)-tailMin]
	default:
		return str[:len(str)-tailMinSec]
	}
}

// type check
var _ encoding.TextMarshaler = Duration(0)

// MarshalText implements the [encoding.TextMarshaler] interface for Duration.
func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// type check
var _ encoding.TextUnmarshaler = (*Duration)(nil)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface for
// *Duration.  Unlike default [time.ParseDuration], this implementation supports
// day units, e.g.:
//
//	"1.5d"
//	"10d"
//	"-1d"
func (d *Duration) UnmarshalText(b []byte) (err error) {
	b, sign, err := extractSign(b)
	if err != nil {
		return fmt.Errorf("parsing duration: %w", err)
	}

	input := bytes.Runes(b)

	if len(input) == 1 && input[0] == '0' {
		*d = 0

		return nil
	}

	var days time.Duration
	remainUnits := input
	commonUnitsLen := 0
	for len(remainUnits) > 0 {
		var (
			daysParsed time.Duration
			l          int
		)
		remainUnits, l, daysParsed, err = parseDurationUnit(remainUnits)
		if err != nil {
			return fmt.Errorf("parsing duration unit: %w", err)
		}

		days += daysParsed
		commonUnitsLen += l
	}

	if commonUnitsLen == 0 {
		*d = applySign(Duration(days), sign)

		return nil
	}

	// TODO(f.setrakov): Consider reimplementing [time.ParseDuration] to
	// decrease the amount of allocations.
	timeDur, err := time.ParseDuration(string(input[:commonUnitsLen]))
	if err != nil {
		return fmt.Errorf("parsing common duration units: %w", err)
	}

	result := days + timeDur
	*d = applySign(Duration(result), sign)

	return nil
}

// extractSign makes sure that the input has a valid length and extracts the
// sign from it.
func extractSign(b []byte) (res []byte, sign byte, err error) {
	if len(b) == 0 {
		return nil, 0, errors.ErrEmptyValue
	}

	res = b
	if b[0] == '+' || b[0] == '-' {
		if len(b) == 1 {
			return nil, 0, fmt.Errorf("%w: empty input after sign", ErrInvalidDuration)
		}

		sign = b[0]
		res = b[1:]
	}

	return res, sign, nil
}

// applySign negates d if sign is `-`, otherwise returns d unchanged.
func applySign(d Duration, sign byte) (res Duration) {
	if sign == '-' {
		return -d
	}

	return d
}

// parseDurationUnit parses a single duration unit from the input, removing day
// units.  It returns the resulting slice, the length of the duration unit if it
// is a common one, and the duration in days.
func parseDurationUnit(input []rune) (
	res []rune,
	commonUnitLen int,
	days time.Duration,
	err error,
) {
	pos := slices.IndexFunc(input, isNotFloatChar)
	if pos == -1 || !isDurationLetter(input[pos]) {
		return nil, 0, 0, ErrInvalidDuration
	}

	unit := input[pos]
	pos++

	if unit != durationUnitDay {
		return input[pos:], pos, 0, nil
	}

	num, err := strconv.ParseFloat(string(input[:pos-1]), 64)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("parsing days number: %w", err)
	}

	copy(input, input[pos:])
	input = input[:len(input)-pos]

	return input, 0, time.Duration(num * float64(Day)), nil
}

// isNotFloatChar returns true if r is not a digit or a point.
func isNotFloatChar(r rune) (ok bool) {
	return (r < '0' || r > '9') && r != '.'
}

// isDurationLetter returns true if r can be a part of a duration unit.
func isDurationLetter(r rune) (ok bool) {
	switch r {
	case 'd',
		'h',
		'm',
		'n',
		's',
		'u',
		// NOTE: Check both of those characters, as they have different Unicode
		// values.
		'µ',
		'μ':
		return true
	default:
		return false
	}
}
