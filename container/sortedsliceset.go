package container

import (
	"cmp"
	"fmt"
	"slices"
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
func NewSortedSliceSet[T cmp.Ordered](elems ...T) (set *SortedSliceSet[T]) {
	slices.Sort(elems)

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
