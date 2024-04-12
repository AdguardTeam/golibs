//go:build unix

package osutil

import (
	"os"

	"golang.org/x/sys/unix"
)

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

// notifyShutdownSignal notifies c on receiving Unix shutdown signals using n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, unix.SIGINT, unix.SIGQUIT, unix.SIGTERM)
}
