//go:build windows

package osutil

import (
	"os"
	"syscall"
)

// shutdownSignals is a list of actual Windows shutdown signals.
//
// NOTE:  Use syscall.SIGTERM as opposed to windows.SIGTERM, because that's the
// type that the Go runtime is sending.
var shutdownSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
}

// isReconfigureSignal returns true if sig is a Windows reconfigure signal.
// Since Windows doesn't seem to have a Unix-compatible mechanism of signaling a
// change in the configuration, it always returns false.
func isReconfigureSignal(_ os.Signal) (ok bool) {
	return false
}

// notifyReconfigureSignal notifies c on receiving Windows reconfigure signals
// using n.  Since Windows doesn't seem to have a Unix-compatible mechanism of
// signaling a change in the configuration, it does nothing.
func notifyReconfigureSignal(_ SignalNotifier, _ chan<- os.Signal) {}

// notifyShutdownSignal notifies c on receiving Windows shutdown signals using
// n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	// [syscall.SIGTERM] is processed automatically.  See go doc os/signal,
	// section Windows.
	n.Notify(c, os.Interrupt)
}
