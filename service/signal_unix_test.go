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
			onStop: func(_ chan<- os.Signal) { panic("not implemented") },
		},
		Logger:          slogutil.NewDiscardLogger(),
		RefreshTimeout:  testTimeout,
		ShutdownTimeout: testTimeout,
	})

	require.NotNil(t, controlCh)

	sigHdlr.AddRefresher(refr)

	testutil.RequireSend(t, controlCh, os.Signal(unix.SIGHUP), testTimeout)

	status := sigHdlr.Handle(context.Background())
	assert.Equal(t, osutil.ExitCodeSuccess, status)

	testutil.RequireReceive(t, refrCh, testTimeout)
}
