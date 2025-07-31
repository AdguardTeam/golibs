// Package contextutil contains types and utilities for working with contexts.
package contextutil

import (
	"context"
	"time"
)

// Constructor is an interface for constructing contexts with deadlines, e.g.
// for request contexts.
type Constructor interface {
	// New returns a new context based on parent as well as a cancel function.
	// parent, ctx, and cancel must not be nil.
	New(parent context.Context) (ctx context.Context, cancel context.CancelFunc)
}

// EmptyConstructor is the implementation of the [Constructor] interface that
// returns the parent context and an empty [context.CancelFunc].
type EmptyConstructor struct{}

// type check
var _ Constructor = EmptyConstructor{}

// New implements the [Constructor] interface for EmptyConstructor.
func (EmptyConstructor) New(
	parent context.Context,
) (ctx context.Context, cancel context.CancelFunc) {
	return parent, func() {}
}

// TimeoutConstructor is an implementation of the [Constructor] interface that
// returns a context with the given timeout.
type TimeoutConstructor struct {
	timeout time.Duration
}

// NewTimeoutConstructor returns a new properly initialized *TimeoutConstructor.
func NewTimeoutConstructor(timeout time.Duration) (c *TimeoutConstructor) {
	return &TimeoutConstructor{
		timeout: timeout,
	}
}

// type check
var _ Constructor = (*TimeoutConstructor)(nil)

// New implements the [Constructor] interface for *TimeoutConstructor.  It
// returns a context with its timeout and the corresponding cancellation
// function.
func (c *TimeoutConstructor) New(
	parent context.Context,
) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(parent, c.timeout)
}
