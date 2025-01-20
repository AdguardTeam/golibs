package service

import (
	"cmp"
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/contextutil"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/timeutil"
)

// RefreshWorker is an [Interface] implementation that updates its [Refresher]
// every tick of the provided ticker.
type RefreshWorker struct {
	done           chan unit
	clock          timeutil.ClockAfter
	contextCons    contextutil.Constructor
	errHdlr        ErrorHandler
	refr           Refresher
	schedule       timeutil.Schedule
	refrOnShutdown bool
}

// RefreshWorkerConfig is the configuration structure for a *RefreshWorker.
type RefreshWorkerConfig struct {
	// Clock is used for time-related operations of the refresher.  If it is
	// nil, [timeutil.SystemClock] is used.
	Clock timeutil.ClockAfter

	// ContextConstructor is used to provide a context for the Refresh method of
	// Refresher.  If it is nil, [contextutil.EmptyConstructor] is used.
	ContextConstructor contextutil.Constructor

	// ErrorHandler is used to handle the errors arising during the run of the
	// worker.  The passed context is the same context that has been passed to
	// Refresh.  If it is nil, [IgnoreErrorHandler] is used.
	ErrorHandler ErrorHandler

	// Refresher is the entity being refreshed.  It must not be nil.
	//
	// The error returned by its Refresh method is returned from
	// [RefreshWorker.Shutdown] only when
	// [RefreshWorkerConfig.RefreshOnShutdown] is true.  In all other cases,
	// [RefreshWorkerConfig.ErrorHandler] is used.
	Refresher Refresher

	// Schedule defines when a refresher is run.  Schedule.UntilNext is called
	// after each refresh to determine how long the worker sleeps until the next
	// refresh.  It must not be nil.
	Schedule timeutil.Schedule

	// RefreshOnShutdown, if true, instructs the worker to refresh before
	// shutting down the worker.  This is useful for refreshers that use Refresh
	// to persist to disk or remote storage before shutting down.
	RefreshOnShutdown bool
}

// NewRefreshWorker returns a new valid *RefreshWorker with the provided
// parameters.  c must not be nil.
func NewRefreshWorker(c *RefreshWorkerConfig) (w *RefreshWorker) {
	return &RefreshWorker{
		done: make(chan unit),
		contextCons: cmp.Or[contextutil.Constructor](
			c.ContextConstructor,
			contextutil.EmptyConstructor{},
		),
		clock:          cmp.Or[timeutil.ClockAfter](c.Clock, timeutil.SystemClock{}),
		errHdlr:        cmp.Or[ErrorHandler](c.ErrorHandler, IgnoreErrorHandler{}),
		refr:           c.Refresher,
		schedule:       c.Schedule,
		refrOnShutdown: c.RefreshOnShutdown,
	}
}

// type check
var _ Interface = (*RefreshWorker)(nil)

// Start implements the [Interface] interface for *RefreshWorker.  err is always
// nil.
//
// If ctx has a logger added with [slogutil.ContextWithLogger], that logger is
// used to report panics; otherwise, [slog.Default] is used to report them.
func (w *RefreshWorker) Start(ctx context.Context) (err error) {
	go w.refreshInALoop(ctx)

	return nil
}

// refreshInALoop refreshes the entity in accordance with w.schedule until
// Shutdown is called.
func (w *RefreshWorker) refreshInALoop(ctx context.Context) {
	defer slogutil.RecoverAndLogDefault(ctx)

	waitDur := w.schedule.UntilNext(w.clock.Now())

	for {
		select {
		case <-w.done:
			return
		case <-w.clock.After(waitDur):
			err := w.refresh(ctx)
			if err != nil {
				w.errHdlr.Handle(ctx, err)
			}

			waitDur = w.schedule.UntilNext(w.clock.Now())
		}
	}
}

// refresh refreshes the entity using w.contextCons to create a context for the
// refresh.
func (w *RefreshWorker) refresh(ctx context.Context) (err error) {
	ctx, cancel := w.contextCons.New(ctx)
	defer cancel()

	return w.refr.Refresh(ctx)
}

// Shutdown implements the [Interface] interface for *RefreshWorker.
func (w *RefreshWorker) Shutdown(ctx context.Context) (err error) {
	close(w.done)

	if w.refrOnShutdown {
		// TODO(a.garipov):  Ensure that a refresh isn't already running.
		err = w.refresh(ctx)
		if err != nil {
			return fmt.Errorf("refresh on shutdown: %w", err)
		}
	}

	return nil
}
