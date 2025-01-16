package timeutil

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/robfig/cron/v3"
)

// Schedule is an interface for entities that can decide when a task should be
// performed based on the current time.
type Schedule interface {
	// UntilNext returns the duration left until the next time when a task
	// should be performed.  UntilNext can return different values with the same
	// now value.  d must not be negative.
	UntilNext(now time.Time) (d time.Duration)
}

// ConstSchedule is a [Schedule] for tasks that run with a constant interval
// between the runs.
type ConstSchedule struct {
	interval time.Duration
}

// NewConstSchedule returns a new schedule that runs with a constant interval.
// ivl must be positive.
func NewConstSchedule(ivl time.Duration) (s *ConstSchedule) {
	err := validate.Positive("ivl", ivl)
	if err != nil {
		panic(fmt.Errorf("timeutil.NewConstSchedule: %w", err))
	}

	return &ConstSchedule{
		interval: ivl,
	}
}

// type check
var _ Schedule = (*ConstSchedule)(nil)

// UntilNext implements the [Schedule] interface for *ConstSchedule.  It always
// returns the original interval used during its construction.
func (s *ConstSchedule) UntilNext(_ time.Time) (d time.Duration) {
	return s.interval
}

// CronSchedule is an adapter for the cron module.
type CronSchedule struct {
	cron cron.Schedule
}

// NewCronSchedule returns a new cron schedule adapted to the [Schedule]
// interface.  cs must not be nil.
func NewCronSchedule(cs cron.Schedule) (s *CronSchedule) {
	if cs == nil {
		panic(fmt.Errorf("timeutil.NewCronSchedule: cs: %w", errors.ErrNoValue))
	}

	return &CronSchedule{
		cron: cs,
	}
}

// type check
var _ Schedule = (*CronSchedule)(nil)

// UntilNext implements the [Schedule] interface for *CronSchedule.
func (s *CronSchedule) UntilNext(now time.Time) (d time.Duration) {
	return max(s.cron.Next(now).Sub(now), 0)
}

// RandomizedSchedule adds a random duration to the result of the [Schedule]
// that it wraps.
type RandomizedSchedule struct {
	rand  *rand.Rand
	sched Schedule
	min   time.Duration
	max   time.Duration
}

// NewRandomizedSchedule returns a new schedule that adds a random value between
// minAdd and maxAdd to the result of sched.  sched and r must not be nil.  max
// must be greater than min.
func NewRandomizedSchedule(
	sched Schedule,
	r *rand.Rand,
	minAdd time.Duration,
	maxAdd time.Duration,
) (s *RandomizedSchedule) {
	errs := []error{
		validate.NotNil("r", r),
		validate.GreaterThan("maxAdd", maxAdd, minAdd),
	}

	if sched == nil {
		errs = append(errs, fmt.Errorf("sched: %w", errors.ErrNoValue))
	}

	err := errors.Join(errs...)
	if err != nil {
		panic(fmt.Errorf("timeutil.NewRandomizedSchedule: %w", err))
	}

	return &RandomizedSchedule{
		sched: sched,
		rand:  r,
		min:   minAdd,
		max:   maxAdd,
	}
}

// type check
var _ Schedule = (*RandomizedSchedule)(nil)

// UntilNext implements the [Schedule] interface for *RandomizedSchedule.
func (s *RandomizedSchedule) UntilNext(now time.Time) (d time.Duration) {
	ivl := s.max - s.min
	if ivl < 0 {
		ivl = -ivl
	}

	added := s.rand.Int64N(int64(ivl)) + int64(s.min)

	return max(s.sched.UntilNext(now)+time.Duration(added), 0)
}
