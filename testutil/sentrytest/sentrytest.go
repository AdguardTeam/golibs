// Package sentrytest contains fake implementations of interfaces for the Sentry
// module.
//
// TODO(a.garipov):  Add more utilities or rename to fakesentry.
package sentrytest

import (
	"context"
	"time"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/getsentry/sentry-go"
)

// Transport is a [sentry.Transport] implementation for tests.
type Transport struct {
	OnClose            func()
	OnConfigure        func(opts sentry.ClientOptions)
	OnFlush            func(timeout time.Duration) (ok bool)
	OnFlushWithContext func(ctx context.Context) (ok bool)
	OnSendEvent        func(e *sentry.Event)
}

// type check
var _ sentry.Transport = (*Transport)(nil)

// Close implements the [sentry.Transport] interface for the *Transport.
func (t *Transport) Close() {
	t.OnClose()
}

// Configure implements the [sentry.Transport] interface for the *Transport.
func (t *Transport) Configure(opts sentry.ClientOptions) {
	t.OnConfigure(opts)
}

// Flush implements the [sentry.Transport] interface for the *Transport.
func (t *Transport) Flush(timeout time.Duration) (ok bool) {
	return t.OnFlush(timeout)
}

// FlushWithContext implements the [sentry.Transport] interface for the
// *Transport.
func (t *Transport) FlushWithContext(ctx context.Context) (ok bool) {
	return t.OnFlushWithContext(ctx)
}

// SendEvent implements the [sentry.Transport] interface for the *Transport.
func (t *Transport) SendEvent(e *sentry.Event) {
	t.OnSendEvent(e)
}

// NewTransport returns a new *Transport all methods of which panic.
func NewTransport() (tst *Transport) {
	return &Transport{
		OnClose:     func() { panic(testutil.UnexpectedCall()) },
		OnConfigure: func(opts sentry.ClientOptions) { panic(testutil.UnexpectedCall(opts)) },
		OnFlush: func(timeout time.Duration) (ok bool) {
			panic(testutil.UnexpectedCall(timeout))
		},
		OnFlushWithContext: func(ctx context.Context) (ok bool) {
			panic(testutil.UnexpectedCall(ctx))
		},
		OnSendEvent: func(e *sentry.Event) { panic(testutil.UnexpectedCall(e)) },
	}
}
