// Package timeutil contains types and utilities for dealing with time and
// duration values.
package timeutil

import (
	"encoding"
	"fmt"
	"time"
)

// Day is the duration of one day.
const Day time.Duration = 24 * time.Hour

// Duration is a helper type for time.Duration providing functionality for
// encoding.
type Duration time.Duration

// type check
var _ fmt.Stringer = Duration(0)

// String implements the [fmt.Stringer] interface for Duration.  It wraps
// time.Duration.String method and additionally cuts off non-leading zero values
// of minutes and seconds.  Some values which are differ between the
// implementations:
//
//	Duration:   "1m", time.Duration:   "1m0s"
//	Duration:   "1h", time.Duration: "1h0m0s"
//	Duration: "1h1m", time.Duration: "1h1m0s"
func (d Duration) String() (str string) {
	timeDur := time.Duration(d)
	str = timeDur.String()

	const (
		tailMin    = len(`0s`)
		tailMinSec = len(`0m0s`)
	)

	const (
		secsInHour = time.Hour / time.Second
		minsInHour = time.Hour / time.Minute
	)

	switch rounded := timeDur / time.Second; {
	case
		rounded == 0,
		rounded*time.Second != timeDur,
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
// *Duration.
//
// TODO(e.burkov):  Fix parsing larger units like days.
func (d *Duration) UnmarshalText(b []byte) (err error) {
	timeDur, err := time.ParseDuration(string(b))
	if err != nil {
		// Don't wrap the error, because it's the only one here.
		return err
	}

	*d = Duration(timeDur)

	return nil
}
