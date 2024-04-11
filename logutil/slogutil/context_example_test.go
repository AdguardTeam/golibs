package slogutil_test

import (
	"context"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func ExampleContextWithLogger() {
	handler := func(ctx context.Context) {
		l := slogutil.MustLoggerFromContext(ctx)

		l.Info("handling")
	}

	l := slogutil.New(nil)
	l = l.With("request_id", 123)

	ctx := context.Background()
	ctx = slogutil.ContextWithLogger(ctx, l)

	handler(ctx)

	// Output:
	// INFO handling request_id=123
}
