package httputil

import (
	"context"
	"net/http"
	"time"
)

// RequestMetrics is an interface for collection of HTTP request statistics.
type RequestMetrics interface {
	// ObserveRequest observes a single HTTP request.  It should be called after
	// the request has finished.  w and r must not be nil, dur must be positive.
	ObserveRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, dur time.Duration)
}

// EmptyRequestMetrics is an implementation of the [RequestMetrics] interface
// that does nothing.
type EmptyRequestMetrics struct{}

// type check
var _ RequestMetrics = EmptyRequestMetrics{}

// ObserveRequest implements the [RequestMetrics] interface for
// EmptyRequestMetrics.
func (EmptyRequestMetrics) ObserveRequest(
	_ context.Context,
	_ http.ResponseWriter,
	_ *http.Request,
	_ time.Duration,
) {
}
