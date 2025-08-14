package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/osutil"
	"github.com/AdguardTeam/golibs/service"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSignalNotifier is a fake [osutil.SignalNotifier] implementation for
// tests.
type fakeSignalNotifier struct {
	onNotify func(c chan<- os.Signal, sig ...os.Signal)
	onStop   func(c chan<- os.Signal)
}

// type check
var _ osutil.SignalNotifier = (*fakeSignalNotifier)(nil)

// Notify implements the [osutil.SignalNotifier] interface for
// *fakeSignalNotifier.
func (s *fakeSignalNotifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	s.onNotify(c, sig...)
}

// Stop implements the [osutil.SignalNotifier] interface for
// *fakeSignalNotifier.
func (s *fakeSignalNotifier) Stop(c chan<- os.Signal) {
	s.onStop(c)
}

func TestSignalHandler(t *testing.T) {
	shutdownCh := make(chan struct{})
	svc := &fakeservice.Service{
		OnStart: func(ctx context.Context) (err error) { panic(testutil.UnexpectedCall(ctx)) },
		OnShutdown: func(_ context.Context) (err error) {
			close(shutdownCh)

			return nil
		},
	}

	var controlCh chan<- os.Signal
	sigHdlr := service.NewSignalHandler(&service.SignalHandlerConfig{
		SignalNotifier: &fakeSignalNotifier{
			onNotify: func(c chan<- os.Signal, sig ...os.Signal) {
				controlCh = c
			},
			onStop: func(ch chan<- os.Signal) { panic(testutil.UnexpectedCall(ch)) },
		},
		Logger:          slogutil.NewDiscardLogger(),
		RefreshTimeout:  testTimeout,
		ShutdownTimeout: testTimeout,
	})

	require.NotNil(t, controlCh)

	sigHdlr.AddService(svc)

	testutil.RequireSend(t, controlCh, os.Interrupt, testTimeout)

	status := sigHdlr.Handle(context.Background())
	assert.Equal(t, osutil.ExitCodeSuccess, status)

	testutil.RequireReceive(t, shutdownCh, testTimeout)
}
