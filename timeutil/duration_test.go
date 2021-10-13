package timeutil_test

import (
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestDuration_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		val  time.Duration
	}{{
		name: "1s",
		val:  time.Second,
	}, {
		name: "1m",
		val:  time.Minute,
	}, {
		name: "1h",
		val:  time.Hour,
	}, {
		name: "1m1s",
		val:  time.Minute + time.Second,
	}, {
		name: "1h1m",
		val:  time.Hour + time.Minute,
	}, {
		name: "1h0m1s",
		val:  time.Hour + time.Second,
	}, {
		name: "1ms",
		val:  time.Millisecond,
	}, {
		name: "1h0m0.001s",
		val:  time.Hour + time.Millisecond,
	}, {
		name: "1.001s",
		val:  time.Second + time.Millisecond,
	}, {
		name: "1m1.001s",
		val:  time.Minute + time.Second + time.Millisecond,
	}, {
		name: "0s",
		val:  0,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := timeutil.Duration{Duration: tc.val}
			assert.Equal(t, tc.name, d.String())
		})
	}
}

func TestDuration_encoding(t *testing.T) {
	t.Parallel()

	v := &timeutil.Duration{
		Duration: time.Millisecond,
	}

	testutil.AssertMarshalText(t, "1ms", v)
	testutil.AssertUnmarshalText(t, "1ms", v)
}
