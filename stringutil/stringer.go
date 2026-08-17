package stringutil

import (
	"fmt"

	"github.com/AdguardTeam/golibs/testutil"
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

// NewTestStringer returns a [fmt.Stringer] implementation that does nothing and
// panics in [fmt.Stringer.String].
func NewTestStringer() (s *Stringer) {
	return &Stringer{
		OnString: func() (str string) { panic(testutil.UnexpectedCall()) },
	}
}
