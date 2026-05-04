package timeutil_test

import (
	"fmt"
	"time"

	"github.com/AdguardTeam/golibs/timeutil"
)

func ExampleDuration_String() {
	for _, value := range []time.Duration{
		time.Minute,
		time.Hour,
		timeutil.Day,
		36 * time.Hour,
		timeutil.Day + time.Hour + 30*time.Minute,
	} {
		d := timeutil.Duration(value)
		fmt.Println(d)
	}

	// Output:
	// 1m
	// 1h
	// 1d
	// 1d12h
	// 1d1h30m
}

func ExampleDuration_UnmarshalText() {
	var got timeutil.Duration
	err := got.UnmarshalText([]byte(`1.5d28m150ms`))
	fmt.Println(err)

	want := 36*time.Hour + 28*time.Minute + 150*time.Millisecond
	fmt.Printf("%t\n", timeutil.Duration(want) == got)

	// Output:
	// <nil>
	// true
}
