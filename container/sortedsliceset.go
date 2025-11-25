package container

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/AdguardTeam/golibs/errors"
)

// SortedSliceSet is a simple set implementation that has a sorted set of values
// as its underlying storage.
//
// TODO(a.garipov):  Consider relaxing the type requirement or adding a version
// with a comparison function.
type SortedSliceSet[T cmp.Ordered] struct {
	elems []T
}

// NewSortedSliceSet returns a new *SortedSliceSet.  elems must not be modified
// after calling NewSortedSliceSet.
//
// TODO(f.setrakov): Delegate allocation decisions to the user.
func NewSortedSliceSet[T cmp.Ordered](elems ...T) (set *SortedSliceSet[T]) {
	slices.Sort(elems)
	elems = slices.Compact(elems)

	return &SortedSliceSet[T]{
		elems: elems,
	}
}

// Add adds v to set.
func (set *SortedSliceSet[T]) Add(v T) {
	i, ok := slices.BinarySearch(set.elems, v)
	if !ok {
		set.elems = slices.Insert(set.elems, i, v)
	}
}

// Clear clears set in a way that retains the internal storage for later reuse
// to reduce allocations.  Calling Clear on a nil set has no effect, just like a
// clear on a nil slice doesn't.
func (set *SortedSliceSet[T]) Clear() {
	if set != nil {
		clear(set.elems)
		set.elems = set.elems[:0]
	}
}

// Clone returns a clone of set.  If set is nil, clone is nil.
//
// NOTE:  It calls [slices.Clone] on the underlying storage, so these elements
// are cloned shallowly.
func (set *SortedSliceSet[T]) Clone() (clone *SortedSliceSet[T]) {
	if set == nil {
		return nil
	}

	return NewSortedSliceSet(slices.Clone(set.elems)...)
}

// Delete deletes v from set.
func (set *SortedSliceSet[T]) Delete(v T) {
	i, ok := slices.BinarySearch(set.elems, v)
	if ok {
		set.elems = slices.Delete(set.elems, i, i+1)
	}
}

// Equal returns true if set is equal to other.  set and other may be nil; Equal
// returns true if both are nil, but a nil *SortedSliceSet is not equal to a
// non-nil empty one.
func (set *SortedSliceSet[T]) Equal(other *SortedSliceSet[T]) (ok bool) {
	if set == nil || other == nil {
		return set == other
	}

	return slices.Equal(set.elems, other.elems)
}

// Has returns true if v is in set.  Calling Has on a nil set returns false,
// just like iterating over a nil or empty slice does.
func (set *SortedSliceSet[T]) Has(v T) (ok bool) {
	if set == nil {
		return false
	}

	_, ok = slices.BinarySearch(set.elems, v)

	return ok
}

// Len returns the length of set.  A nil set has a length of zero, just like an
// nil or empty slice.
func (set *SortedSliceSet[T]) Len() (n int) {
	if set == nil {
		return 0
	}

	return len(set.elems)
}

// Range calls f with each value of set in their sorted order.  If cont is
// false, Range stops the iteration.  Calling Range on a nil *SortedSliceSet has
// no effect.
func (set *SortedSliceSet[T]) Range(f func(v T) (cont bool)) {
	if set == nil {
		return
	}

	for _, v := range set.elems {
		if !f(v) {
			break
		}
	}
}

// type check
var _ fmt.Stringer = (*SortedSliceSet[int])(nil)

// String implements the [fmt.Stringer] interface for *SortedSliceSet.  Calling
// String on a nil *SortedSliceSet does not panic.
func (set *SortedSliceSet[T]) String() (s string) {
	return fmt.Sprintf("%v", set.Values())
}

// Values returns the underlying slice of set.  values must not be modified.
// Values returns nil if set is nil.
func (set *SortedSliceSet[T]) Values() (values []T) {
	if set == nil {
		return nil
	}

	return set.elems
}

// Union fills set with values belonging to either a or b.  This function
// guarantees zero-allocation, but may not perform well with large sets.  Union
// returns empty set if both a and b are nil.  set must not be nil.  If neither
// a nor b are equal to set, then the function will rewrite the contents of set.
func (set *SortedSliceSet[T]) Union(a, b *SortedSliceSet[T]) (res *SortedSliceSet[T]) {
	if set == nil {
		panic(fmt.Errorf("set: %v", errors.ErrNoValue))
	}

	if a == nil && b == nil {
		set.Clear()

		return set
	}

	if a == nil {
		set.elems = append(set.elems[:0], b.elems...)

		return set
	}

	if b == nil {
		set.elems = append(set.elems[:0], a.elems...)

		return set
	}

	if set == a {
		set.addMissing(b)

		return set
	}

	if set == b {
		set.addMissing(a)

		return set
	}

	set.elems = set.elems[:0]
	set.union(a, b)

	return set
}

// union merges two SortedSliceSets producing a sorted result.  set, a and b
// must not be nil.
func (set *SortedSliceSet[T]) union(a, b *SortedSliceSet[T]) {
	aIdx, bIdx := 0, 0
	for aIdx < len(a.elems) && bIdx < len(b.elems) {
		if a.elems[aIdx] < b.elems[bIdx] {
			set.elems = append(set.elems, a.elems[aIdx])
			aIdx++
		} else if a.elems[aIdx] > b.elems[bIdx] {
			set.elems = append(set.elems, b.elems[bIdx])
			bIdx++
		} else {
			set.elems = append(set.elems, b.elems[bIdx])
			aIdx++
			bIdx++
		}
	}

	set.elems = append(set.elems, a.elems[aIdx:]...)
	set.elems = append(set.elems, b.elems[bIdx:]...)
}

// addMissing adds all the elements from other that are not present in set.  set
// and other must not be nil.
func (set *SortedSliceSet[T]) addMissing(other *SortedSliceSet[T]) {
	otherIdx, setIdx := 0, 0

	for otherIdx < len(set.elems) && setIdx < len(other.elems) {
		if set.elems[otherIdx] < other.elems[setIdx] {
			otherIdx++
		} else if set.elems[otherIdx] > other.elems[setIdx] {
			set.elems = append(set.elems, other.elems[setIdx])
			setIdx++
		} else {
			otherIdx++
			setIdx++
		}
	}

	set.elems = append(set.elems, other.elems[setIdx:]...)
	slices.Sort(set.elems)
}

// Intersection fills set with values that belong to both a and b.  This
// function guarantees zero-allocation, but may not perform well with large
// sets.  If you need better performance, consider using [MapSet].  Intersection
// returns an empty set if one of the arguments is nil.  set must not be nil.
// If neither a nor b are equal to set, then the function will rewrite the
// contents of set.
func (set *SortedSliceSet[T]) Intersection(a, b *SortedSliceSet[T]) (res *SortedSliceSet[T]) {
	if set == nil {
		panic(fmt.Errorf("set: %v", errors.ErrNoValue))
	}

	if a == nil || b == nil {
		set.Clear()

		return set
	}

	if set == a {
		return set.removeMissing(b)
	}

	if set == b {
		return set.removeMissing(a)
	}

	set.elems = set.elems[:0]

	return set.intersection(a, b)
}

// removeMissing removes all elements from other that are not present in set.
// set and other must not be nil.
func (set *SortedSliceSet[T]) removeMissing(other *SortedSliceSet[T]) (res *SortedSliceSet[T]) {
	setIdx, otherIdx := 0, 0
	lastSavedIdx := 0

	for setIdx < len(set.elems) && otherIdx < len(other.elems) {
		if set.elems[setIdx] < other.elems[otherIdx] {
			setIdx++
		} else if set.elems[setIdx] > other.elems[otherIdx] {
			otherIdx++
		} else {
			set.elems[lastSavedIdx] = set.elems[setIdx]
			lastSavedIdx++
			setIdx++
			otherIdx++
		}
	}

	clear(set.elems[lastSavedIdx:])
	set.elems = set.elems[:lastSavedIdx]

	return set
}

// intersection fills set with values that belong both to b and a.  res will be
// sorted.  set, a and b must not be nil.
func (set *SortedSliceSet[T]) intersection(b, a *SortedSliceSet[T]) (res *SortedSliceSet[T]) {
	aIdx, bIdx := 0, 0
	for aIdx < len(a.elems) && bIdx < len(b.elems) {
		if a.elems[aIdx] < b.elems[bIdx] {
			aIdx++
		} else if a.elems[aIdx] > b.elems[bIdx] {
			bIdx++
		} else {
			set.elems = append(set.elems, a.elems[aIdx])
			aIdx++
			bIdx++
		}
	}

	return set
}

// Intersects returns true if set and other has at least one common element.  If
// set or other is nil, result will be false.
func (set *SortedSliceSet[T]) Intersects(other *SortedSliceSet[T]) (ok bool) {
	if set == nil || other == nil {
		return false
	}

	setIdx, otherIdx := 0, 0
	for setIdx < len(set.elems) && otherIdx < len(other.elems) {
		if set.elems[setIdx] < other.elems[otherIdx] {
			setIdx++
		} else if set.elems[setIdx] > other.elems[otherIdx] {
			otherIdx++
		} else {
			return true
		}
	}

	return false
}
