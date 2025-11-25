package container_test

import (
	"fmt"
	"slices"

	"github.com/AdguardTeam/golibs/container"
)

func ExampleMapSet() {
	const x = 1
	set := container.NewMapSet[int]()

	ok := set.Has(x)
	fmt.Printf("%s contains %v is %t\n", container.MapSetToString(set), x, ok)

	set.Add(x)
	ok = set.Has(x)
	fmt.Printf("%s contains %v is %t\n", container.MapSetToString(set), x, ok)

	other := container.NewMapSet(x)
	ok = set.Equal(other)
	fmt.Printf(
		"%s is equal to %s is %t\n",
		container.MapSetToString(set),
		container.MapSetToString(other),
		ok,
	)

	set.Add(2)
	values := set.Values()
	slices.Sort(values)
	fmt.Printf("values of %s are %v\n", container.MapSetToString(set), values)

	set.Delete(x)
	ok = set.Has(x)
	fmt.Printf("%s contains %v is %t\n", container.MapSetToString(set), x, ok)

	for n := range set.Range {
		fmt.Printf("got value %d\n", n)

		break
	}

	set = container.NewMapSet(x)
	fmt.Printf("%s has length %d\n", container.MapSetToString(set), set.Len())

	set.Clear()
	fmt.Printf("%s has length %d\n", container.MapSetToString(set), set.Len())

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

func ExampleMapSet_Clone() {
	var set *container.MapSet[int]
	fmt.Printf("nil:   %#v\n", set.Clone())

	const x, y = 1, 2
	set = container.NewMapSet(x)
	clone := set.Clone()
	clone.Add(y)

	fmt.Printf("orig:  %t %t\n", set.Has(x), set.Has(y))
	fmt.Printf("clone: %t %t\n", clone.Has(x), clone.Has(y))

	// Output:
	// nil:   (*container.MapSet[int])(nil)
	// orig:  true false
	// clone: true true
}

func ExampleMapSet_Equal() {
	const x, y = 1, 2
	set := container.NewMapSet(x)

	fmt.Printf("same:       %t\n", set.Equal(container.NewMapSet(x)))
	fmt.Printf("other elem: %t\n", set.Equal(container.NewMapSet(y)))
	fmt.Printf("other len:  %t\n", set.Equal(container.NewMapSet(x, y)))
	fmt.Printf("nil:        %t\n", set.Equal(nil))
	fmt.Printf("nil eq nil: %t\n", (*container.MapSet[int])(nil).Equal(nil))

	// Output:
	// same:       true
	// other elem: false
	// other len:  false
	// nil:        false
	// nil eq nil: true
}

func ExampleMapSet_Intersection() {
	a := container.NewMapSet(1, 6, 10)
	b := container.NewMapSet(3, 6, 12)
	set := container.NewMapSet[int]()

	fmt.Printf("a = %s, b = %s\n", container.MapSetToString(a), container.MapSetToString(b))
	fmt.Printf("set = a ∩ b:     %s\n", container.MapSetToString(set.Intersection(a, b)))
	fmt.Printf("set = nil ∩ nil: %s\n", container.MapSetToString(set.Intersection(nil, nil)))
	fmt.Printf("set = nil ∩ b:   %s\n", container.MapSetToString(set.Intersection(nil, b)))
	fmt.Printf("set = a ∩ nil:   %s\n", container.MapSetToString(set.Intersection(a, nil)))
	fmt.Printf("a = a ∩ b:       %s\n", container.MapSetToString(a.Intersection(a, b)))

	a = container.NewMapSet(1, 6, 10)
	fmt.Printf("b = a ∩ b:       %s\n", container.MapSetToString(b.Intersection(a, b)))

	// Output:
	// a = [1 6 10], b = [3 6 12]
	// set = a ∩ b:     [6]
	// set = nil ∩ nil: []
	// set = nil ∩ b:   []
	// set = a ∩ nil:   []
	// a = a ∩ b:       [6]
	// b = a ∩ b:       [6]
}

func ExampleMapSet_Union() {
	a := container.NewMapSet(1, 6, 10)
	b := container.NewMapSet(3, 6, 12)
	set := container.NewMapSet[int]()

	fmt.Printf("a = %s, b = %s\n", container.MapSetToString(a), container.MapSetToString(b))
	fmt.Printf("set = a ∪ b:     %s\n", container.MapSetToString(set.Union(a, b)))
	fmt.Printf("set = nil ∪ nil: %s\n", container.MapSetToString(set.Union(nil, nil)))
	fmt.Printf("set = nil ∪ b:   %s\n", container.MapSetToString(set.Union(nil, b)))
	fmt.Printf("set = a ∪ nil:   %s\n", container.MapSetToString(set.Union(a, nil)))
	fmt.Printf("a = a ∪ b:       %s\n", container.MapSetToString(a.Union(a, b)))

	a = container.NewMapSet(1, 6, 10)
	fmt.Printf("b = a ∪ b:       %s\n", container.MapSetToString(b.Union(a, b)))

	// Output:
	// a = [1 6 10], b = [3 6 12]
	// set = a ∪ b:     [1 3 6 10 12]
	// set = nil ∪ nil: []
	// set = nil ∪ b:   [3 6 12]
	// set = a ∪ nil:   [1 6 10]
	// a = a ∪ b:       [1 3 6 10 12]
	// b = a ∪ b:       [1 3 6 10 12]
}

func ExampleMapSet_Intersects() {
	a := container.NewMapSet(1, 6, 10)
	b := container.NewMapSet(3, 6, 12)
	var nilSet *container.MapSet[int]

	fmt.Printf("a = %s, b = %s\n", container.MapSetToString(a), container.MapSetToString(b))
	fmt.Printf("a ∩ b:     %t\n", a.Intersects(b))
	fmt.Printf("nil ∩ a:   %t\n", nilSet.Intersects(a))
	fmt.Printf("a ∩ nil:   %t\n", a.Intersects(nilSet))
	fmt.Printf("nil ∩ nil: %t\n", nilSet.Intersects(nilSet))

	// Output:
	// a = [1 6 10], b = [3 6 12]
	// a ∩ b:     true
	// nil ∩ a:   false
	// a ∩ nil:   false
	// nil ∩ nil: false
}

func ExampleMapSet_nil() {
	const x = 1

	var set *container.MapSet[int]

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

	func() {
		defer setPanicked()

		set.Intersects(set)
	}()
	fmt.Printf("panic after intersects: %t\n", panicked)

	// Output:
	//
	// panic after clear: false
	// panic after delete: false
	// panic after has: false
	// panic after len: false
	// panic after range: false
	// panic after values: false
	// panic after add: true
	// panic after union: true
	// panic after intersection: true
	// panic after intersects: false
}
