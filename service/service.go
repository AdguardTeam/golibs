// Package service defines types and interfaces for long-running services that
// can be started and shut down.
package service

import (
	"context"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
)

// unit is a convenient alias for struct{}.
type unit = struct{}

// Interface is the interface for long-running services.
type Interface interface {
	// Start starts the service.  ctx is used for cancellation.
	//
	// It is recommended that Start returns only after the service has
	// completely finished its initialization.  If that cannot be done, the
	// implementation of Start must document that.
	Start(ctx context.Context) (err error)

	Shutdowner
}

// Shutdowner is the interface for types that have a Shutdown method but not
// necessarily a Start method.
type Shutdowner interface {
	// Shutdown gracefully stops the service.  ctx is used to determine a
	// timeout before trying to stop the service less gracefully.
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

// ShutdownService is an adapter for types that implement the [Shutdowner]
// interface but not the [Interface] one.
type ShutdownService struct {
	shutdowner Shutdowner
}

// NewShutdownService returns a properly initialized *ShutdownService.  s must
// not be nil.
func NewShutdownService(s Shutdowner) (svc Interface) {
	errors.Check(validate.NotNilInterface("s", s))

	return &ShutdownService{
		shutdowner: s,
	}
}

// type check
var _ Interface = (*ShutdownService)(nil)

// Start implements the [Interface] interface for *ShutdownService.  It always
// returns nil.
func (s *ShutdownService) Start(_ context.Context) (err error) { return nil }

// Shutdown implements the [Interface] interface for *ShutdownService.
func (s *ShutdownService) Shutdown(ctx context.Context) (err error) {
	return s.shutdowner.Shutdown(ctx)
}
