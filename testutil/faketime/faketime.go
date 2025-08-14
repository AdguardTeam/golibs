// Package faketime contains fake implementations of interfaces from package
// timeutil.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic(testutil.UnexpectedCall(arg1, arg2))
package faketime

import (
	"time"

	"github.com/AdguardTeam/golibs/timeutil"
)

// Clock is the [timeutil.Clock] implementation for tests.
type Clock struct {
	OnNow func() (now time.Time)
}

// type check
var _ timeutil.Clock = (*Clock)(nil)

// Now implements the [timeutil.Clock] interface for *Clock.
func (c *Clock) Now() (now time.Time) {
	return c.OnNow()
}

// Schedule is the [timeutil.Schedule] implementation for tests.
type Schedule struct {
	OnUntilNext func(now time.Time) (d time.Duration)
}

// type check
var _ timeutil.Schedule = (*Schedule)(nil)

// UntilNext implements the [timeutil.Schedule] interface for *Schedule.
func (s *Schedule) UntilNext(now time.Time) (d time.Duration) {
	return s.OnUntilNext(now)
}
