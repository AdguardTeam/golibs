package mapsutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/mapsutil"
)

func ExampleSortedRange() {
	m := map[string]int{
		"b": 200,
		"a": 100,
		"c": 300,
		"d": 400,
	}

	mapsutil.SortedRange(m, func(k string, v int) (cont bool) {
		fmt.Printf("value for %q is %d\n", k, v)

		// Do not print any values after "c".
		return k != "c"
	})

	// Output:
	// value for "a" is 100
	// value for "b" is 200
	// value for "c" is 300
}

func ExampleSortedRangeError() {
	checkKey := func(k string, v int) (err error) {
		if k == "x" {
			return fmt.Errorf("bad key: %q", k)
		}

		fmt.Printf("value for %q is %d\n", k, v)

		return nil
	}

	err := mapsutil.SortedRangeError(
		map[string]int{
			"b": 200,
			"a": 100,
			"c": 300,
		},
		checkKey,
	)

	fmt.Println(err)

	err = mapsutil.SortedRangeError(
		map[string]int{
			"x": 0,
		},
		checkKey,
	)

	fmt.Println(err)

	// Output:
	// value for "a" is 100
	// value for "b" is 200
	// value for "c" is 300
	// <nil>
	// bad key: "x"
}
