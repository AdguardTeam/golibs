package stringutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/stringutil"
)

func ExampleSet() {
	const s = "a"
	const nl = "\n"

	set := stringutil.NewSet()

	ok := set.Has(s)
	fmt.Printf(`%s contains "a" is %t`+nl, set, ok)

	set.Add(s)
	ok = set.Has(s)
	fmt.Printf(`%s contains "a" is %t`+nl, set, ok)

	fmt.Printf(`values of %s are %q`+nl, set, set.Values())

	set.Del(s)
	ok = set.Has(s)
	fmt.Printf(`%s contains "a" is %t`+nl, set, ok)

	set = stringutil.NewSet(s)
	fmt.Printf(`%s has length %d`+nl, set, set.Len())

	// Output:
	//
	// [] contains "a" is false
	// ["a"] contains "a" is true
	// values of ["a"] are ["a"]
	// [] contains "a" is false
	// ["a"] has length 1
}

func ExampleSet_nil() {
	const s = "a"

	var set *stringutil.Set

	panicked := false
	setPanicked := func() {
		if v := recover(); v != nil {
			panicked = true
		}
	}

	func() {
		defer setPanicked()

		set.Del(s)
	}()
	fmt.Printf("panic after del: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Has(s)
	}()
	fmt.Printf("panic after has: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Len()
	}()
	fmt.Printf("panic after len: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Values()
	}()
	fmt.Printf("panic after values: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Add(s)
	}()
	fmt.Printf("panic after add: %t\n", panicked)

	// Output:
	//
	// panic after del: false
	// panic after has: false
	// panic after len: false
	// panic after values: false
	// panic after add: true
}
