package otelutil

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpanf is a helper for formatted span names.
func StartSpanf(
	parent context.Context,
	tracer trace.Tracer,
	spanNameFmt string,
	args ...any,
) (ctx context.Context, span trace.Span) {
	return tracer.Start(parent, fmt.Sprintf(spanNameFmt, args...))
}

// EndSpan is a deferred helper that records the error sets the status of the
// span before ending it.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		// TODO(a.garipov):  Think of a better description here that doesn't
		// duplicate the work done by span.RecordError above.  Perhaps, the type
		// of the deepest error in the error chain?
		span.SetStatus(codes.Error, "error")
	}

	span.End()
}
