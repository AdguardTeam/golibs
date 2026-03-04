package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// PanicT can be used with the helpers from package require in cases when
// testing.T and similar standard test helpers aren't safe for use, e.g. stub
// HTTP handlers and goroutines.
//
// While this type also implements [assert.TestingT], prefer to use require
// helper functions, since this helper panics, which immediately fails the test.
type PanicT struct {
	tb testing.TB
}

// NewPanicT creates a new instance of [PanicT].  It is more preferred to use
// this constructor instead of creating [PanicT] as a struct literal.
func NewPanicT(tb testing.TB) (pt PanicT) {
	return PanicT{
		tb: tb,
	}
}

// RequireT is an extension of the require.TestingT interface that contains
// additional methods used to improve error reporting.
type RequireT interface {
	require.TestingT
	Name() (name string)
	Helper()
}

// type check
var _ RequireT = PanicT{}

// Errorf implements the [RequireT] interface for PanicT.  It panics with an
// error with the given format.
func (PanicT) Errorf(format string, args ...any) {
	panic(fmt.Errorf(format, args...))
}

// FailNow implements the [RequireT] interface for PanicT.  It is assumed that
// it will never actually be called, since [PanicT.Errorf] panics.
func (PanicT) FailNow() {
	panic("test failed")
}

// Used by [PanicT.Name] method of PanicT as a no-op return.
const unknownName = "Unknown"

// Name implements the [RequireT] interface for PanicT.  If p was created with a
// non-nil tb, it returns tb.Name.  Otherwise, it returns a stub name.
func (p PanicT) Name() (name string) {
	if p.tb != nil {
		return p.tb.Name()
	}

	return unknownName
}

// Helper implements the [RequireT] interface for PanicT.  If p was created with
// a non-nil tb and tb implements the Helper method, the Helper method is
// called.  Otherwise, it does nothing.
func (p PanicT) Helper() {
	if h, ok := p.tb.(RequireT); ok {
		h.Helper()
	}
}
