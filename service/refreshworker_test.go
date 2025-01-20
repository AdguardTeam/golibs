package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/contextutil"
	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/service"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeservice"
	"github.com/AdguardTeam/golibs/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testIvl                  = 5 * time.Millisecond
	testIvlLong              = 1 * time.Hour
	name                     = "test refresher"
	testError   errors.Error = "test error"
)

// newTestRefresher is a helper that returns refr and linked syncCh channel.
func newTestRefresher(t *testing.T, respErr error) (refr *fakeservice.Refresher, syncCh chan unit) {
	t.Helper()

	pt := testutil.PanicT{}

	syncCh = make(chan unit, 1)
	refr = &fakeservice.Refresher{
		OnRefresh: func(_ context.Context) (err error) {
			testutil.RequireSend(pt, syncCh, unit{}, testTimeout)

			return respErr
		},
	}

	return refr, syncCh
}

// newRefrConfig returns worker configuration.
func newRefrConfig(
	t *testing.T,
	refr service.Refresher,
	ivl time.Duration,
	refrOnShutDown bool,
) (conf *service.RefreshWorkerConfig) {
	t.Helper()

	return &service.RefreshWorkerConfig{
		ContextConstructor: contextutil.NewTimeoutConstructor(testTimeout),
		Refresher:          refr,
		Schedule:           timeutil.NewConstSchedule(ivl),
		RefreshOnShutdown:  refrOnShutDown,
	}
}

func TestRefreshWorker(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		refr, syncCh := newTestRefresher(t, nil)

		w := service.NewRefreshWorker(newRefrConfig(t, refr, testIvl, false))

		err := w.Start(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)

		testutil.RequireReceive(t, syncCh, testTimeout)

		err = w.Shutdown(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)
	})

	t.Run("success_on_shutdown", func(t *testing.T) {
		refr, syncCh := newTestRefresher(t, nil)
		errCh := make(chan unit, 1)

		w := service.NewRefreshWorker(newRefrConfig(t, refr, testIvlLong, true))

		err := w.Start(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)

		err = w.Shutdown(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)

		testutil.RequireReceive(t, syncCh, testTimeout)
		require.Empty(t, errCh)
	})

	t.Run("error", func(t *testing.T) {
		refrWithError, syncCh := newTestRefresher(t, testError)

		w := service.NewRefreshWorker(newRefrConfig(t, refrWithError, testIvl, false))

		err := w.Start(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)

		testutil.RequireReceive(t, syncCh, testTimeout)

		err = w.Shutdown(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)
	})

	t.Run("error_on_shutdown", func(t *testing.T) {
		refrWithError, syncCh := newTestRefresher(t, testError)

		w := service.NewRefreshWorker(newRefrConfig(t, refrWithError, testIvlLong, true))

		err := w.Start(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)

		err = w.Shutdown(testutil.ContextWithTimeout(t, testTimeout))
		assert.ErrorIs(t, err, testError)

		testutil.RequireReceive(t, syncCh, testTimeout)
	})
}
