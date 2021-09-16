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

	other := stringutil.NewSet("a")
	ok = set.Equal(other)
	fmt.Printf("%s is equal to %s is %t\n", set, other, ok)

	fmt.Printf("values of %s are %q\n", set, set.Values())

	set.Range(func(s string) (cont bool) {
		fmt.Printf("got value %q\n", s)

		return false
	})

	set.Del(s)
	ok = set.Has(s)
	fmt.Printf(`%s contains "a" is %t`+nl, set, ok)

	set = stringutil.NewSet(s)
	fmt.Printf("%s has length %d\n", set, set.Len())

	// Output:
	//
	// [] contains "a" is false
	// ["a"] contains "a" is true
	// ["a"] is equal to ["a"] is true
	// values of ["a"] are ["a"]
	// got value "a"
	// [] contains "a" is false
	// ["a"] has length 1
}

func ExampleSet_Equal() {
	set := stringutil.NewSet("a")

	fmt.Printf("same:       %t\n", set.Equal(stringutil.NewSet("a")))
	fmt.Printf("other elem: %t\n", set.Equal(stringutil.NewSet("b")))
	fmt.Printf("other len:  %t\n", set.Equal(stringutil.NewSet("a", "b")))
	fmt.Printf("nil:        %t\n", set.Equal(nil))
	fmt.Printf("nil eq nil: %t\n", (*stringutil.Set)(nil).Equal(nil))

	// Output:
	// same:       true
	// other elem: false
	// other len:  false
	// nil:        false
	// nil eq nil: true
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

		set.Range(func(s string) (cont bool) {
			fmt.Printf("got value %q\n", s)

			return true
		})
	}()
	fmt.Printf("panic after range: %t\n", panicked)

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
	// panic after range: false
	// panic after values: false
	// panic after add: true
}
