package timeutil_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/timeutil"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstSchedule(t *testing.T) {
	const ivl = 5 * time.Minute
	s := timeutil.NewConstSchedule(ivl)

	now := time.Now()
	d := s.UntilNext(now)
	assert.Equal(t, ivl, d)
}

func TestCronSchedule(t *testing.T) {
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cs, err := p.Parse("*/5 * * * *")
	require.NoError(t, err)

	s := timeutil.NewCronSchedule(cs)

	now := time.Now()
	for range 100_000 {
		d := s.UntilNext(now)
		assert.InDelta(t, 5*time.Minute, d, float64(5*time.Minute))
	}
}

func TestRandomizedSchedule(t *testing.T) {
	const (
		ivl    = 5 * time.Minute
		minAdd = -1 * time.Minute
		maxAdd = 1 * time.Minute
	)

	now := time.Now()
	r := rand.New(rand.NewPCG(uint64(now.Unix()), 0))
	s := timeutil.NewRandomizedSchedule(timeutil.NewConstSchedule(ivl), r, minAdd, maxAdd)

	for range 100_000 {
		d := s.UntilNext(now)
		assert.InDelta(t, ivl, d, float64(2*time.Minute))
	}
}
