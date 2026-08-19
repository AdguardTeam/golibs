package slogutil

import (
	"fmt"
	"log/slog"
)

// StringerLogValuer is a wrapper over the fmt.Stringer interface.  It can be
// used to optimise allocations.
type StringerLogValuer[T fmt.Stringer] struct {
	value T
}

// NewStringerLogValuer returns a StringerLogValuer for v.
func NewStringerLogValuer[T fmt.Stringer](v T) (s StringerLogValuer[T]) {
	return StringerLogValuer[T]{
		value: v,
	}
}

// type check
var _ slog.LogValuer = (*StringerLogValuer[fmt.Stringer])(nil)

// LogValue implements the [slog.LogValuer] interface for StringerLogValuer.
func (s StringerLogValuer[T]) LogValue() (l slog.Value) {
	return slog.StringValue(s.value.String())
}
