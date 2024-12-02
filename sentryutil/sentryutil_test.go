package sentryutil_test

import (
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/sentryutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/sentrytest"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportPanics(t *testing.T) {
	t.Parallel()

	tr := sentrytest.NewTransport()
	tr.OnConfigure = func(_ sentry.ClientOptions) {
		// Do nothing.
	}

	eventsCh := make(chan *sentry.Event, 2)
	tr.OnSendEvent = func(e *sentry.Event) {
		eventsCh <- e
	}

	hub := sentry.CurrentHub()
	prevClient := hub.Client()
	t.Cleanup(func() {
		hub.BindClient(prevClient)
	})

	require.NoError(t, sentry.Init(sentry.ClientOptions{
		Dsn:       "https://user:password@does.not.exist/test",
		Transport: tr,
	}))

	require.True(t, t.Run("no panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			defer sentryutil.ReportPanics()

			// Do nothing.
		})

		assert.Empty(t, eventsCh)
	}))

	require.True(t, t.Run("panic", func(t *testing.T) {
		const testError errors.Error = "test error"

		assert.Panics(t, func() {
			defer sentryutil.ReportPanics()

			panic(testError)
		})

		got, _ := testutil.RequireReceive(t, eventsCh, 1*time.Second)
		require.NotNil(t, got)
		require.NotEmpty(t, got.Exception)

		assert.Equal(t, testError.Error(), got.Exception[0].Value)
	}))
}
