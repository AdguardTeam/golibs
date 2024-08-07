package container

// KeyValue is a simple structure for pairs of values where one identifies the
// other.  It is useful when multiple values need to be processed in a defined
// order and each value is associated with its own name or ID.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// KeyValues is a slice of [KeyValue]s.
type KeyValues[K comparable, V any] []KeyValue[K, V]
