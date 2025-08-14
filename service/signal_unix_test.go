//go:build unix

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
	"golang.org/x/sys/unix"
)

func TestSignalHandler_unix(t *testing.T) {
	shutdownCh := make(chan struct{})
	svc := &fakeservice.Service{
		OnStart: func(ctx context.Context) (err error) { panic(testutil.UnexpectedCall(ctx)) },
		OnShutdown: func(_ context.Context) (err error) {
			close(shutdownCh)

			return nil
		},
	}

	refrCh := make(chan struct{})
	refr := &fakeservice.Refresher{
		OnRefresh: func(_ context.Context) (err error) {
			close(refrCh)

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
	sigHdlr.AddRefresher(refr)

	go func() {
		pt := &testutil.PanicT{}

		status := sigHdlr.Handle(context.Background())
		assert.Equal(pt, osutil.ExitCodeSuccess, status)
	}()

	testutil.RequireSend(t, controlCh, os.Signal(unix.SIGHUP), testTimeout)
	testutil.RequireReceive(t, refrCh, testTimeout)

	testutil.RequireSend(t, controlCh, os.Interrupt, testTimeout)
	testutil.RequireReceive(t, shutdownCh, testTimeout)
}
