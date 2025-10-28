//go:build plan9

package osutil

import (
	"os"
	"syscall"
)

const sigquit = syscall.Note("quit")

// isReconfigureSignal returns true if sig is a Plan 9 reconfigure signal.
func isReconfigureSignal(sig os.Signal) (ok bool) {
	return sig == syscall.SIGHUP
}

// isShutdownSignal returns true if sig is a Plan 9 shutdown signal.
func isShutdownSignal(sig os.Signal) (ok bool) {
	switch sig {
	case
		syscall.SIGINT,
		sigquit:
		return true
	default:
		return false
	}
}

// notifyReconfigureSignal notifies c on receiving Plan 9 reconfigure signals
// using n.
func notifyReconfigureSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, syscall.SIGHUP)
}

// notifyShutdownSignal notifies c on receiving Plan 9 shutdown signals using n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, syscall.SIGINT, sigquit)
}
