package container

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	"github.com/AdguardTeam/golibs/errors"
)

// MapSet is a set that uses a map as its storage.
//
// TODO(a.garipov): Figure out a way to add a reproducible String method.
type MapSet[T comparable] struct {
	m map[T]unit
}

// NewMapSet returns a new map set containing values.
//
// TODO(f.setrakov): Delegate allocation decisions to the user.
func NewMapSet[T comparable](values ...T) (set *MapSet[T]) {
	set = &MapSet[T]{
		m: make(map[T]unit, len(values)),
	}

	for _, v := range values {
		set.Add(v)
	}

	return set
}

// Add adds v to set.
func (set *MapSet[T]) Add(v T) {
	set.m[v] = unit{}
}

// Clear clears set in a way that retains the internal storage for later reuse
// to reduce allocations.  Calling Clear on a nil set has no effect, just like a
// clear on a nil map doesn't.
func (set *MapSet[T]) Clear() {
	if set != nil {
		clear(set.m)
	}
}

// Clone returns a deep clone of set.  If set is nil, clone is nil.
func (set *MapSet[T]) Clone() (clone *MapSet[T]) {
	if set == nil {
		return nil
	}

	return &MapSet[T]{
		m: maps.Clone(set.m),
	}
}

// Delete deletes v from set.  Calling Delete on a nil set has no effect, just
// like delete on a nil map doesn't.
func (set *MapSet[T]) Delete(v T) {
	if set != nil {
		delete(set.m, v)
	}
}

// Equal returns true if set is equal to other.  set and other may be nil; Equal
// returns true if both are nil, but a nil *MapSet is not equal to a non-nil
// empty one.
func (set *MapSet[T]) Equal(other *MapSet[T]) (ok bool) {
	if set == nil || other == nil {
		return set == other
	}

	return maps.Equal(set.m, other.m)
}

// Has returns true if v is in set.  Calling Has on a nil set returns false,
// just like indexing on an empty map does.
func (set *MapSet[T]) Has(v T) (ok bool) {
	if set != nil {
		_, ok = set.m[v]
	}

	return ok
}

// Len returns the length of set.  A nil set has a length of zero, just like an
// empty map.
func (set *MapSet[T]) Len() (n int) {
	if set == nil {
		return 0
	}

	return len(set.m)
}

// Range calls f with each value of set in an undefined order.  If cont is
// false, Range stops the iteration.  Calling Range on a nil *MapSet has no
// effect.
func (set *MapSet[T]) Range(f func(v T) (cont bool)) {
	if set == nil {
		return
	}

	for v := range set.m {
		if !f(v) {
			break
		}
	}
}

// Values returns all values in set.  The order of the values is undefined.
// Values returns nil if set is nil.
func (set *MapSet[T]) Values() (values []T) {
	if set == nil {
		return nil
	}

	values = make([]T, 0, len(set.m))
	for v := range set.m {
		values = append(values, v)
	}

	return values
}

// MapSetToString converts a [MapSet] of values of an ordered type into a
// reproducible string.
func MapSetToString[T cmp.Ordered](set *MapSet[T]) (s string) {
	v := set.Values()
	slices.Sort(v)

	return fmt.Sprintf("%v", v)
}

// MapSetToStringFunc is like [MapSetToString] but uses an explicit comparison
// function.
func MapSetToStringFunc[T comparable](set *MapSet[T], compare func(a, b T) (res int)) (s string) {
	v := set.Values()
	slices.SortStableFunc(v, compare)

	return fmt.Sprintf("%v", v)
}

// Union fills set with values belonging to either a or b.  set must not be nil.
// Union returns empty set if both a and b are nil. If neither a nor b are equal
// to set, then the function will rewrite the contents of set.
func (set *MapSet[T]) Union(a, b *MapSet[T]) (res *MapSet[T]) {
	if set == nil {
		panic(fmt.Errorf("set: %v", errors.ErrNoValue))
	}

	if a == nil && b == nil {
		set.Clear()

		return set
	}

	if set != a && set != b {
		set.Clear()
	}

	if a != nil && set != a {
		maps.Copy(set.m, a.m)
	}

	if b != nil && set != b {
		maps.Copy(set.m, b.m)
	}

	return set
}

// Intersection fills set with values that belong both to a and b.  set must not
// be nil.  Intersection returns empty set if one of the arguments is nil.  If
// neither a nor b are equal to set, then the function will rewrite the contents
// of set.
func (set *MapSet[T]) Intersection(a, b *MapSet[T]) (res *MapSet[T]) {
	if set == nil {
		panic(fmt.Errorf("set: %v", errors.ErrNoValue))
	}

	if a == nil || b == nil {
		set.Clear()

		return set
	}

	if set == a {
		return set.intersection(b)
	}

	if set == b {
		return set.intersection(a)
	}

	set.Clear()

	for v := range a.m {
		if _, ok := b.m[v]; ok {
			set.Add(v)
		}
	}

	return set
}

// intersection removes all elements from set that are not present in other.
// other must not be nil.
func (set *MapSet[T]) intersection(other *MapSet[T]) (res *MapSet[T]) {
	for v := range set.m {
		if _, ok := other.m[v]; !ok {
			delete(set.m, v)
		}
	}

	return set
}

// Intersects returns true if set and other has at least one common element.  If
// set or other is nil, result will be false.
func (set *MapSet[T]) Intersects(other *MapSet[T]) (ok bool) {
	if set == nil || other == nil {
		return false
	}

	for v := range set.m {
		if _, ok = other.m[v]; ok {
			return true
		}
	}

	return false
}
