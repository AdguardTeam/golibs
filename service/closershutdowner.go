package service

import (
	"context"
	"io"
)

// CloserShutdowner implements the [Shutdowner] interface.  It wraps an
// [io.Closer] in order to close it on shutdown.
type CloserShutdowner struct {
	closer io.Closer
}

// type check
var _ Shutdowner = (*CloserShutdowner)(nil)

// Shutdown implements the [Shutdowner] interface for *CloserShutdowner.
func (s *CloserShutdowner) Shutdown(_ context.Context) (err error) {
	return s.closer.Close()
}

// NewCloserShutdowner returns a [Shutdowner] that closes the given [io.Closer].
// c must not be nil.
func NewCloserShutdowner(c io.Closer) (s Shutdowner) {
	return &CloserShutdowner{
		closer: c,
	}
}
