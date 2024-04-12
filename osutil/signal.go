package osutil

import (
	"os"
	"os/signal"
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

// DefaultSignalNotifier is a [SignalNotifier] that uses [signal.Notify] and
// [signal.Stop].
type DefaultSignalNotifier struct{}

// type check
var _ SignalNotifier = DefaultSignalNotifier{}

// Notify implements the [SignalNotifier] interface for DefaultSignalNotifier.
func (n DefaultSignalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

// Stop implements the [SignalNotifier] interface for DefaultSignalNotifier.
func (n DefaultSignalNotifier) Stop(c chan<- os.Signal) {
	signal.Stop(c)
}

// IsShutdownSignal returns true if sig is a shutdown signal.
func IsShutdownSignal(sig os.Signal) (ok bool) {
	return isShutdownSignal(sig)
}

// NotifyShutdownSignal notifies c on receiving shutdown signals using n.
func NotifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	notifyShutdownSignal(n, c)
}
