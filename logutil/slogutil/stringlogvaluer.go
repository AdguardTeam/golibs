package slogutil

import (
	"fmt"
	"log/slog"
)

// StringLogValuer is a wrapper over the [fmt.Stringer] interface.
type StringLogValuer[T fmt.Stringer] struct {
	value T
}

// NewStringLogValuer returns a [StringLogValuer] for v.
func NewStringLogValuer[T fmt.Stringer](v T) (s StringLogValuer[T]) {
	return StringLogValuer[T]{
		value: v,
	}
}

// type check
var _ slog.LogValuer = (*StringLogValuer[fmt.Stringer])(nil)

// LogValue implements the [slog.LogValuer] interface for [StringLogValuer].
func (s StringLogValuer[T]) LogValue() (l slog.Value) {
	return slog.StringValue(s.value.String())
}
