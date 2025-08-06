package httputil_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/faketime"
	"github.com/stretchr/testify/assert"
)

// testRequestMetrics is the [httputil.RequestMetrics] implementation for tests.
type testRequestMetrics struct {
	onObserveRequest func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		dur time.Duration,
	)
}

// type check
var _ httputil.RequestMetrics = (*testRequestMetrics)(nil)

// ObserveRequest implements the [RequestMetrics] interface for
// *testRequestMetrics.
func (m *testRequestMetrics) ObserveRequest(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	dur time.Duration,
) {
	m.onObserveRequest(ctx, w, r, dur)
}

func TestMetricsMiddleware(t *testing.T) {
	const wantStatus = http.StatusCreated
	const wantDuration = 1 * time.Second

	offset := &atomic.Uint32{}
	now := time.Now()

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	req := httptest.NewRequest(http.MethodGet, testPath, nil).WithContext(ctx)

	mw := httputil.NewMetricsMiddleware(&httputil.MetricsMiddlewareConfig{
		Clock: &faketime.Clock{
			OnNow: func() (n time.Time) {
				return now.Add(time.Duration(offset.Add(1)) * wantDuration)
			},
		},
		Metrics: &testRequestMetrics{
			onObserveRequest: func(
				_ context.Context,
				w http.ResponseWriter,
				r *http.Request,
				dur time.Duration,
			) {
				assert.Equal(t, req, r)
				assert.Equal(t, wantDuration, dur)

				rec := w.(*httptest.ResponseRecorder)
				assert.Equal(t, wantStatus, rec.Code)
			},
		},
	})

	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(wantStatus)
	}))

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
}
