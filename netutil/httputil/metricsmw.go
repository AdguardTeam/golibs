package httputil

import (
	"cmp"
	"net/http"

	"github.com/AdguardTeam/golibs/timeutil"
)

// MetricsMiddleware wraps a handler with [RequestMetrics].
type MetricsMiddleware struct {
	clock   timeutil.Clock
	metrics RequestMetrics
}

// MetricsMiddlewareConfig is the configuration structure for a
// *MetricsMiddleware.
type MetricsMiddlewareConfig struct {
	// Clock is used to measure the times for requests' durations.  If Clock is
	// nil, [timeutil.SystemClock] is used.
	Clock timeutil.Clock

	// Metrics is used to observe HTTP requests.  It must not be nil.
	Metrics RequestMetrics
}

// NewMetricsMiddleware returns a new properly initialized *MetricsMiddleware.
// c must not be nil and must be valid.
func NewMetricsMiddleware(c *MetricsMiddlewareConfig) (mw *MetricsMiddleware) {
	return &MetricsMiddleware{
		clock:   cmp.Or[timeutil.Clock](c.Clock, timeutil.SystemClock{}),
		metrics: c.Metrics,
	}
}

// type check
var _ Middleware = (*MetricsMiddleware)(nil)

// Wrap implements the [Middleware] interface for *MetricsMiddleware.
func (mw *MetricsMiddleware) Wrap(h http.Handler) (wrapped http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := mw.clock.Now()

		ctx := r.Context()

		defer func() { mw.metrics.ObserveRequest(ctx, w, r, mw.clock.Now().Sub(start)) }()

		h.ServeHTTP(w, r)
	})
}
