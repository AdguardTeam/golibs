// Package mapsutil contains utilities for map handling.
package mapsutil

import (
	"slices"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

// SortedRange is like the usual Go range but sorts the keys before iterating
// ensuring a predictable order.  If cont is false, SortedRange stops the
// iteration.
func SortedRange[K constraints.Ordered, V any, M ~map[K]V](m M, f func(k K, v V) (cont bool)) {
	keys := maps.Keys(m)
	slices.Sort(keys)
	for _, k := range keys {
		if !f(k, m[k]) {
			break
		}
	}
}

// SortedRangeError is like [SortedRange] but uses an error to signal that the
// iteration must be stopped.  err is the same error as the one returned from f,
// or nil if no errors are returned.
func SortedRangeError[K constraints.Ordered, V any, M ~map[K]V](
	m M,
	f func(k K, v V) (err error),
) (err error) {
	keys := maps.Keys(m)
	slices.Sort(keys)
	for _, k := range keys {
		err = f(k, m[k])
		if err != nil {
			return err
		}
	}

	return nil
}
