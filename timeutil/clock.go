package timeutil

import "time"

// Clock is an interface for time-related operations.
type Clock interface {
	// Now returns the current time in accordance with the clock.
	Now() (now time.Time)
}

// ClockAfter is an extension of the Clock interface for clocks that can create
// timers.
type ClockAfter interface {
	Clock

	// After returns a channel on which the current time is sent after d has
	// passed.
	After(d time.Duration) (c <-chan time.Time)
}

// SystemClock is a [Clock] that uses the functions from package time.
type SystemClock struct{}

// type check
var _ ClockAfter = SystemClock{}

// Now implements the [ClockAfter] interface for SystemClock.
func (SystemClock) Now() (now time.Time) { return time.Now() }

// After implements the [ClockAfter] interface for SystemClock.
func (SystemClock) After(d time.Duration) (c <-chan time.Time) { return time.After(d) }
