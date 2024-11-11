//go:build unix

package osutil

import (
	"os"

	"golang.org/x/sys/unix"
)

// isReconfigureSignal returns true if sig is a Unix reconfigure signal.
func isReconfigureSignal(sig os.Signal) (ok bool) {
	return sig == unix.SIGHUP
}

// isShutdownSignal returns true if sig is a Unix shutdown signal.
func isShutdownSignal(sig os.Signal) (ok bool) {
	switch sig {
	case
		unix.SIGINT,
		unix.SIGQUIT,
		unix.SIGTERM:
		return true
	default:
		return false
	}
}

// notifyReconfigureSignal notifies c on receiving Unix reconfigure signals
// using n.
func notifyReconfigureSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, unix.SIGHUP)
}

// notifyShutdownSignal notifies c on receiving Unix shutdown signals using n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, unix.SIGINT, unix.SIGQUIT, unix.SIGTERM)
}
