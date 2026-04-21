package fakeos

import (
	"context"
	"os"

	"github.com/AdguardTeam/golibs/osutil"
)

// Notifier is the mock implementation of the [osutil.ContextSignalNotifier]
// interface for testing.
type Notifier struct {
	OnNotify        func(c chan<- os.Signal, sig ...os.Signal)
	OnStop          func(c chan<- os.Signal)
	OnNotifyContext func(
		parent context.Context,
		sig ...os.Signal,
	) (ctx context.Context, stop context.CancelFunc)
}

// type check
var _ osutil.ContextSignalNotifier = (*Notifier)(nil)

// Notify implements the [osutil.ContextSignalNotifier] interface for *Notifier.
func (n *Notifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	n.OnNotify(c, sig...)
}

// Stop implements the [osutil.ContextSignalNotifier] interface for *Notifier.
func (n *Notifier) Stop(c chan<- os.Signal) {
	n.OnStop(c)
}

// NotifyContext implements the [osutil.ContextSignalNotifier] interface for
// *Notifier.
func (n *Notifier) NotifyContext(
	parent context.Context,
	sig ...os.Signal,
) (ctx context.Context, stop context.CancelFunc) {
	return n.OnNotifyContext(parent, sig...)
}
