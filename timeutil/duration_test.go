package timeutil_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	}, {
		name: "10d15h",
		val:  10*timeutil.Day + 15*time.Hour,
	}, {
		name: "-10d",
		val:  -10 * timeutil.Day,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := timeutil.Duration(tc.val)
			assert.Equal(t, tc.name, d.String())
		})
	}
}

func TestDuration_encoding(t *testing.T) {
	t.Parallel()

	v := timeutil.Duration(time.Millisecond)

	testutil.AssertMarshalText(t, "1ms", &v)
	testutil.AssertUnmarshalText(t, "1ms", &v)
}

func TestDuration_UnmarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		input      string
		wantErrMsg string
		want       time.Duration
	}{{
		name:       "normal",
		input:      "1d",
		wantErrMsg: "",
		want:       time.Duration(timeutil.Day),
	}, {
		name:       "float",
		input:      "1.5d",
		wantErrMsg: "",
		want:       36 * time.Hour,
	}, {
		name:       "day_sum",
		input:      "1.5d0.5d5m1.5d12h5s",
		wantErrMsg: "",
		want:       96*time.Hour + 5*time.Minute + 5*time.Second,
	}, {
		name:       "hour_sum",
		input:      "10h2d15h",
		wantErrMsg: "",
		want:       3*timeutil.Day + time.Hour,
	}, {
		name:       "minute_sum",
		input:      "1m1m1m1m1m1d1h",
		wantErrMsg: "",
		want:       timeutil.Day + time.Hour + 5*time.Minute,
	}, {
		name:  "microseconds",
		input: "1μs1µs",
		want:  2 * time.Microsecond,
	}, {
		name:       "positive",
		input:      "+1.5d",
		wantErrMsg: "",
		want:       36 * time.Hour,
	}, {
		name:       "negative",
		input:      "-1.5d",
		wantErrMsg: "",
		want:       -36 * time.Hour,
	}, {
		name:       "zero",
		input:      "0",
		wantErrMsg: "",
		want:       0,
	}, {
		name:       "signed_zero",
		input:      "-0",
		wantErrMsg: "",
		want:       0,
	}, {
		name:       "no_days",
		input:      "1ms1h1m1s1ms",
		wantErrMsg: "",
		want:       time.Hour + time.Minute + time.Second + 2*time.Millisecond,
	}, {
		name:       "invalid_string",
		input:      "snsnm",
		wantErrMsg: `parsing common duration units: time: invalid duration "snsnm"`,
		want:       0,
	}, {
		name:       "unknown_unit",
		input:      "1p",
		wantErrMsg: `parsing duration unit: ` + timeutil.ErrInvalidDuration.Error(),
		want:       0,
	}, {
		name:       "empty",
		input:      "",
		wantErrMsg: "parsing duration: " + errors.ErrEmptyValue.Error(),
		want:       0,
	}, {
		name:       "out_of_range",
		input:      "61s",
		wantErrMsg: "",
		want:       time.Minute + time.Second,
	}, {
		name:  "plus",
		input: "+",
		wantErrMsg: "parsing duration: " + timeutil.ErrInvalidDuration.Error() +
			": empty input after sign",
		want: 0,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got timeutil.Duration
			err := got.UnmarshalText([]byte(tc.input))
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
			assert.Equal(t, tc.want, time.Duration(got))
		})
	}
}

func BenchmarkDuration_UnmarshalText(b *testing.B) {
	base := "1d1h1m1s"

	for i := 10; i <= 1000; i *= 10 {
		input := []byte(strings.Repeat(base, i))
		b.Run(fmt.Sprintf("%d_units", i*4), func(b *testing.B) {
			var (
				d   timeutil.Duration
				err error
			)

			b.ReportAllocs()
			for b.Loop() {
				err = d.UnmarshalText(input)
			}

			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/timeutil
	//	cpu: Apple M4 Pro
	//	BenchmarkDuration_UnmarshalText/40_units-14         	 1556010	       754.0 ns/op	     384 B/op	       2 allocs/op
	//	BenchmarkDuration_UnmarshalText/400_units-14        	  148866	      8088 ns/op	    3840 B/op	       2 allocs/op
	//	BenchmarkDuration_UnmarshalText/4000_units-14       	    6024	    196837 ns/op	   38912 B/op	       2 allocs/op
}

func FuzzDuration_UnmarshalText(f *testing.F) {
	for _, seed := range []string{
		"1d",
		"1h",
		"1m",
		"1s",
		"1ms",
		"1ns",
		"1.5d",
		"+1.5d",
		"-0.85d",
		"dhsms",
		"1d1h1.0ms",
		"-1d+1d-1d+10000s",
		"1.0h1.0m1.0s1.0ms1.0ns",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		var got timeutil.Duration
		timeutilErr := got.UnmarshalText([]byte(input))

		d, timeErr := time.ParseDuration(input)
		if timeErr == nil {
			require.NoError(t, timeutilErr)

			assert.Equal(t, d, time.Duration(got))
		}
	})
}

func FuzzDuration_МarshalText_roundTrip(f *testing.F) {
	for _, seed := range []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		time.Minute,
		time.Hour,
		timeutil.Day,
		-time.Nanosecond,
		-time.Microsecond,
		-time.Millisecond,
		-time.Second,
		-time.Minute,
		-time.Hour,
		-timeutil.Day,
	} {
		f.Add(int64(seed))
	}

	f.Fuzz(func(t *testing.T, input int64) {
		want := timeutil.Duration(input)

		str := want.String()

		var got timeutil.Duration
		err := got.UnmarshalText([]byte(str))
		require.NoError(t, err)

		assert.Equal(t, want, got)
	})
}
