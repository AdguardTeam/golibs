package syncutil

import "sync"

// Map is a generic wrapper around [sync.Map] for better type safety.
//
// TODO(a.garipov):  Remove once https://github.com/golang/go/issues/71076 is
// implemented.
type Map[K comparable, V any] struct {
	m *sync.Map
}

// NewMap returns a new properly initialized *Map.
func NewMap[K comparable, V any]() (m *Map[K, V]) {
	return &Map[K, V]{
		m: &sync.Map{},
	}
}

// Clear deletes all the entries, resulting in an empty Map.
//
// See [sync.Map.Clear].
func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

// CompareAndDelete deletes the entry for key if its value is equal to old.  The
// old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete returns
// false (even if the old value is the nil interface value).
//
// See [sync.Map.CompareAndDelete].
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

// CompareAndSwap swaps the old and new values for key if the value stored in
// the map is equal to old.
//
// See [sync.Map.CompareAndSwap].
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(key, old, new)
}

// Delete deletes the value for a key.
//
// See [sync.Map.Delete].
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Load returns the value stored in the map for a key, or its zero value if no
// value is present.
//
// See [sync.Map.Load].
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, false
	}

	return v.(V), true
}

// LoadAndDelete deletes the value for a key, returning the previous value if
// any.  The loaded result reports whether the key was present.
//
// See [sync.Map.LoadAndDelete].
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)
	if !loaded {
		return value, false
	}

	return v.(V), true
}

// LoadOrStore returns the existing value for the key if present.  Otherwise, it
// stores and returns the given value.
//
// See [sync.Map.LoadOrStore].
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.m.LoadOrStore(key, value)

	return v.(V), loaded
}

// Store sets the value for a key.
//
// See [sync.Map.Store].
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Swap swaps the value for a key and returns the previous value if any.  The
// loaded result reports whether the key was present.
//
// See [sync.Map.Swap].
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := m.m.Swap(key, value)
	if !loaded {
		return previous, false
	}

	return v.(V), true
}

// Range calls f sequentially for each key and value present in the map.  If f
// returns false, range stops the iteration.
//
// See [sync.Map.Range].
func (m *Map[K, V]) Range(f func(key K, value V) (cont bool)) {
	m.m.Range(func(k, v any) (cont bool) {
		return f(k.(K), v.(V))
	})
}
