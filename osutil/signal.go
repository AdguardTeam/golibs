package osutil

import (
	"context"
	"os"
	"os/signal"
	"slices"
)

// SignalNotifier is the interface for OS functions that can notify about
// incoming signals using a channel.
type SignalNotifier interface {
	// Notify starts relaying incoming signals to c.  If no signals are
	// provided, all incoming signals are relayed to c.  Otherwise, just the
	// provided signals are.
	//
	// Implementations will not block sending to c: the caller must ensure that
	// c has sufficient buffer space to keep up with the expected signal rate.
	// For a channel used for notification of just one signal value, a buffer of
	// size 1 is sufficient.
	//
	// It is allowed to call Notify multiple times with the same channel: each
	// call expands the set of signals sent to that channel.  The only way to
	// remove signals from the set is to call Stop.
	//
	// It is allowed to call Notify multiple times with different channels and
	// the same signals: each channel receives copies of incoming signals
	// independently.
	//
	// See also [signal.Notify].
	Notify(c chan<- os.Signal, sig ...os.Signal)

	// Stop causes the SignalNotifier to stop relaying incoming signals to c.
	// It undoes the effect of all prior calls to Notify using c.  When Stop
	// returns, it c must not receive any more signals.
	//
	// See also [signal.Stop].
	Stop(c chan<- os.Signal)
}

// ContextSignalNotifier extends the [SignalNotifier] interface allowing to
// handle signals with context cancellation.
type ContextSignalNotifier interface {
	SignalNotifier

	// NotifyContext returns a copy of the parent context that will be canceled
	// when one of the given signals arrives, when the stop function is called,
	// or when the parent context is marked done.  This depends on which event
	// happens first.
	//
	// See also [signal.NotifyContext].
	NotifyContext(
		parent context.Context,
		sig ...os.Signal,
	) (ctx context.Context, stop context.CancelFunc)
}

// EmptySignalNotifier is a [ContextSignalNotifier] that does nothing.
type EmptySignalNotifier struct{}

// type check
var _ ContextSignalNotifier = EmptySignalNotifier{}

// Notify implements the [ContextSignalNotifier] interface for
// EmptySignalNotifier.
func (n EmptySignalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {}

// Stop implements the [ContextSignalNotifier] interface for
// EmptySignalNotifier.
func (n EmptySignalNotifier) Stop(c chan<- os.Signal) {}

// NotifyContext implements the [ContextSignalNotifier] interface for
// EmptySignalNotifier.
func (n EmptySignalNotifier) NotifyContext(
	parent context.Context,
	sig ...os.Signal,
) (ctx context.Context, stop context.CancelFunc) {
	return context.WithCancel(ctx)
}

// DefaultSignalNotifier is a [ContextSignalNotifier] that uses [signal.Notify],
// [signal.Stop] and [signal.NotifyContext].
type DefaultSignalNotifier struct{}

// type check
var _ ContextSignalNotifier = DefaultSignalNotifier{}

// Notify implements the [ContextSignalNotifier] interface for
// DefaultSignalNotifier.
func (n DefaultSignalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

// Stop implements the [ContextSignalNotifier] interface for
// DefaultSignalNotifier.
func (n DefaultSignalNotifier) Stop(c chan<- os.Signal) {
	signal.Stop(c)
}

// NotifyContext implements the [ContextSignalNotifier] interface for
// DefaultSignalHandler.
func (n DefaultSignalNotifier) NotifyContext(
	parent context.Context,
	sig ...os.Signal,
) (ctx context.Context, stop context.CancelFunc) {
	return signal.NotifyContext(parent, sig...)
}

// IsReconfigureSignal returns true if sig is a reconfigure signal.
//
// NOTE:  It always returns false on Windows.
func IsReconfigureSignal(sig os.Signal) (ok bool) {
	return isReconfigureSignal(sig)
}

// IsShutdownSignal returns true if sig is a shutdown signal.
func IsShutdownSignal(sig os.Signal) (ok bool) {
	return slices.Contains(shutdownSignals, sig)
}

// NotifyReconfigureSignal notifies c on receiving reconfigure signals using n.
//
// NOTE:  It does nothing on Windows.
func NotifyReconfigureSignal(n SignalNotifier, c chan<- os.Signal) {
	notifyReconfigureSignal(n, c)
}

// NotifyShutdownSignal notifies c on receiving shutdown signals using n.
func NotifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	notifyShutdownSignal(n, c)
}

// NotifyContextShutdownSignal returns a copy of the parent context, which will
// be canceled when the shutdown signal arrives.
func NotifyContextShutdownSignal(
	n ContextSignalNotifier,
	parent context.Context,
) (ctx context.Context, stop context.CancelFunc) {
	return n.NotifyContext(parent, shutdownSignals...)
}
