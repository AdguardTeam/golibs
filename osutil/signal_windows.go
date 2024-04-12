//go:build windows

package osutil

import (
	"os"
	"syscall"
)

// isShutdownSignal returns true if sig is a Windows shutdown signal.
func isShutdownSignal(sig os.Signal) (ok bool) {
	switch sig {
	case os.Interrupt, syscall.SIGTERM:
		return true
	default:
		return false
	}
}

// notifyShutdownSignal notifies c on receiving Windows shutdown signals using
// n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	// syscall.SIGTERM is processed automatically.  See go doc os/signal,
	// section Windows.
	n.Notify(c, os.Interrupt)
}
