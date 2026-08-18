package testutil

import (
	"fmt"
)

// Stringer is a [fmt.Stringer] implementation for tests.
type Stringer struct {
	OnString func() (s string)
}

// type check
var _ fmt.Stringer = (*Stringer)(nil)

// String implements the [fmt.Stringer] interface for *Stringer.
func (s *Stringer) String() string {
	return s.OnString()
}

// NewStringer returns a [fmt.Stringer] implementation that does nothing and
// panics in [fmt.Stringer.String].
func NewStringer() (s *Stringer) {
	return &Stringer{
		OnString: func() (str string) { panic(UnexpectedCall()) },
	}
}
