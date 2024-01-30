// Package service defines types and interfaces for long-running services that
// can be started and shut down.
//
// TODO(a.garipov): Add tests.
package service

import (
	"context"
)

// Interface is the interface for long-running services.
type Interface interface {
	// Start starts the service.  ctx is used for cancelation.
	//
	// It is recommended that Start returns only after the service has
	// completely finished its initialization.  If that cannot be done, the
	// implementation of Start must document that.
	Start(ctx context.Context) (err error)

	// Shutdown gracefully stops the service.  ctx is used to determine
	// a timeout before trying to stop the service less gracefully.
	//
	// It is recommended that Shutdown returns only after the service has
	// completely finished its termination.  If that cannot be done, the
	// implementation of Shutdown must document that.
	Shutdown(ctx context.Context) (err error)
}

// type check
var _ Interface = Empty{}

// Empty is an [Interface] implementation that does nothing and returns nil.
type Empty struct{}

// Start implements the [Interface] interface for Empty.
func (Empty) Start(_ context.Context) (err error) { return nil }

// Shutdown implements the [Interface] interface for Empty.
func (Empty) Shutdown(_ context.Context) (err error) { return nil }
