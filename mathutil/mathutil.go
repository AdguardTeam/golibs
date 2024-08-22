// Package mathutil contains generic helpers for common algorithms and
// mathematic operations.
package mathutil

import "golang.org/x/exp/constraints"

// Number is a type constraint for all numbers.
type Number interface {
	constraints.Integer | constraints.Float
}

// BoolToNumber returns 1 if cond is true and 0 otherwise.
func BoolToNumber[T Number](cond bool) (res T) {
	if cond {
		return 1
	}

	return 0
}
