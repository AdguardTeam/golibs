package container_test

import (
	"fmt"
	"slices"

	"github.com/AdguardTeam/golibs/container"
)

func ExampleSortedSliceSet() {
	const x = 1
	set := container.NewSortedSliceSet[int]()

	ok := set.Has(x)
	fmt.Printf("%s contains %v is %t\n", set, x, ok)

	set.Add(x)
	ok = set.Has(x)
	fmt.Printf("%s contains %v is %t\n", set, x, ok)

	other := container.NewSortedSliceSet(x)
	ok = set.Equal(other)
	fmt.Printf("%s is equal to %s is %t\n", set, other, ok)

	set.Add(2)
	values := set.Values()
	slices.Sort(values)
	fmt.Printf("values of %s are %v\n", set, values)

	set.Delete(x)
	ok = set.Has(x)
	fmt.Printf("%s contains %v is %t\n", set, x, ok)

	for n := range set.Range {
		fmt.Printf("got value %d\n", n)

		break
	}

	set = container.NewSortedSliceSet(x)
	fmt.Printf("%s has length %d\n", set, set.Len())

	set.Clear()
	fmt.Printf("%s has length %d\n", set, set.Len())

	// Output:
	//
	// [] contains 1 is false
	// [1] contains 1 is true
	// [1] is equal to [1] is true
	// values of [1 2] are [1 2]
	// [2] contains 1 is false
	// got value 2
	// [1] has length 1
	// [] has length 0
}

func ExampleSortedSliceSet_Clone() {
	var set *container.SortedSliceSet[int]
	fmt.Printf("nil:   %#v\n", set.Clone())

	const x, y = 1, 2
	set = container.NewSortedSliceSet(x)
	clone := set.Clone()
	clone.Add(y)

	fmt.Printf("orig:  %t %t\n", set.Has(x), set.Has(y))
	fmt.Printf("clone: %t %t\n", clone.Has(x), clone.Has(y))

	// Output:
	// nil:   (*container.SortedSliceSet[int])(nil)
	// orig:  true false
	// clone: true true
}

func ExampleSortedSliceSet_Equal() {
	const x, y = 1, 2
	set := container.NewSortedSliceSet(x)

	fmt.Printf("same:       %t\n", set.Equal(container.NewSortedSliceSet(x)))
	fmt.Printf("other elem: %t\n", set.Equal(container.NewSortedSliceSet(y)))
	fmt.Printf("other len:  %t\n", set.Equal(container.NewSortedSliceSet(x, y)))
	fmt.Printf("nil:        %t\n", set.Equal(nil))
	fmt.Printf("nil eq nil: %t\n", (*container.SortedSliceSet[int])(nil).Equal(nil))

	// Output:
	// same:       true
	// other elem: false
	// other len:  false
	// nil:        false
	// nil eq nil: true
}

func ExampleSortedSliceSet_Intersection() {
	a := container.NewSortedSliceSet(1, 6, 10)
	b := container.NewSortedSliceSet(3, 6, 12)
	set := container.NewSortedSliceSet[int]()

	fmt.Printf("a = [1 6 10], b = [3 6 12]\n")

	fmt.Printf("set = a and b:     %s\n", set.Intersection(a, b).String())

	fmt.Printf("set = nil and nil: %s\n", set.Intersection(nil, nil).String())

	fmt.Printf("set = nil and b:   %s\n", set.Intersection(nil, b).String())

	fmt.Printf("set = a and nil:   %s\n", set.Intersection(a, nil).String())

	fmt.Printf("a = a and b:       %s\n", a.Intersection(a, b))

	a = container.NewSortedSliceSet(1, 6, 10)
	fmt.Printf("b = a and b:       %s\n", b.Intersection(a, b))

	// Output:
	// a = [1 6 10], b = [3 6 12]
	// set = a and b:     [6]
	// set = nil and nil: []
	// set = nil and b:   []
	// set = a and nil:   []
	// a = a and b:       [6]
	// b = a and b:       [6]
}

func ExampleSortedSliceSet_Union() {
	a := container.NewSortedSliceSet(1, 6, 10)
	b := container.NewSortedSliceSet(3, 6, 12)
	set := container.NewSortedSliceSet[int]()

	fmt.Printf("a = [1 6 10], b = [3 6 12]\n")

	fmt.Printf("set = a and b:     %s\n", set.Union(a, b).String())

	fmt.Printf("set = nil and nil: %s\n", set.Union(nil, nil).String())

	fmt.Printf("set = nil and b:   %s\n", set.Union(nil, b).String())

	fmt.Printf("set = a and nil:   %s\n", set.Union(a, nil).String())

	fmt.Printf("a = a and b:       %s\n", a.Union(a, b))

	a = container.NewSortedSliceSet(1, 6, 10)
	fmt.Printf("b = a and b:       %s\n", b.Union(a, b))

	// Output:
	// a = [1 6 10], b = [3 6 12]
	// set = a and b:     [1 3 6 10 12]
	// set = nil and nil: []
	// set = nil and b:   [3 6 12]
	// set = a and nil:   [1 6 10]
	// a = a and b:       [1 3 6 10 12]
	// b = a and b:       [1 3 6 10 12]
}

func ExampleSortedSliceSet_nil() {
	const x = 1

	var set *container.SortedSliceSet[int]

	panicked := false
	setPanicked := func() {
		panicked = recover() != nil
	}

	func() {
		defer setPanicked()

		set.Clear()
	}()
	fmt.Printf("panic after clear: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Delete(x)
	}()
	fmt.Printf("panic after delete: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Has(x)
	}()
	fmt.Printf("panic after has: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Len()
	}()
	fmt.Printf("panic after len: %t\n", panicked)

	func() {
		defer setPanicked()

		for n := range set.Range {
			fmt.Printf("got value %d\n", n)
		}
	}()
	fmt.Printf("panic after range: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Values()
	}()
	fmt.Printf("panic after values: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Add(x)
	}()
	fmt.Printf("panic after add: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Union(set, set)
	}()
	fmt.Printf("panic after union: %t\n", panicked)

	func() {
		defer setPanicked()

		set.Intersection(set, set)
	}()
	fmt.Printf("panic after intersection: %t\n", panicked)

	// Output:
	//
	// panic after clear: false
	// panic after delete: true
	// panic after has: false
	// panic after len: false
	// panic after range: false
	// panic after values: false
	// panic after add: true
	// panic after union: true
	// panic after intersection: true
}
