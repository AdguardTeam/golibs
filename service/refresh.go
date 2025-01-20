package service

import (
	"context"
)

// Refresher is the interface for entities that can update themselves.
type Refresher interface {
	// Refresh is a method that is called to perform some kind of a job on an
	// entity.  Typically, this means updating some data or sending data
	// somewhere else.
	//
	// See [RefreshWorker] for one of the uses.
	Refresh(ctx context.Context) (err error)
}

// RefresherFunc is an adapter to allow the use of ordinary functions as
// [Refresher].
type RefresherFunc func(ctx context.Context) (err error)

// type check
var _ Refresher = RefresherFunc(nil)

// Refresh implements the [Refresher] interface for RefresherFunc.
func (f RefresherFunc) Refresh(ctx context.Context) (err error) { return f(ctx) }

// EmptyRefresher is a [Refresher] that does nothing.
type EmptyRefresher struct{}

// type check
var _ Refresher = EmptyRefresher{}

// Refresh implements the [Refresher] interface for EmptyRefresher.  It returns
// nil immediately.
func (EmptyRefresher) Refresh(_ context.Context) (err error) { return nil }
