//go:build unix

package osutil

import (
	"os"

	"golang.org/x/sys/unix"
)

// shutdownSignals is a list of actual Unix shutdown signals.
var shutdownSignals = []os.Signal{
	unix.SIGINT,
	unix.SIGQUIT,
	unix.SIGTERM,
}

// isReconfigureSignal returns true if sig is a Unix reconfigure signal.
func isReconfigureSignal(sig os.Signal) (ok bool) {
	return sig == unix.SIGHUP
}

// notifyReconfigureSignal notifies c on receiving Unix reconfigure signals
// using n.
func notifyReconfigureSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, unix.SIGHUP)
}

// notifyShutdownSignal notifies c on receiving Unix shutdown signals using n.
func notifyShutdownSignal(n SignalNotifier, c chan<- os.Signal) {
	n.Notify(c, shutdownSignals...)
}
