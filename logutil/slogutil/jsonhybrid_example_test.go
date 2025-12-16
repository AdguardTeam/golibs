package slogutil_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/requestid"
)

func ExampleJSONHybridHandler() {
	ctx := context.Background()
	ctx = requestid.ContextWithRequestID(ctx, testRequestID)

	h := slogutil.NewJSONHybridHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		// Use slogutil.RemoveTime to make the example reproducible.
		ReplaceAttr: slogutil.RemoveTime,
	})
	l := slog.New(h)

	l.DebugContext(ctx, "debug with no attributes")
	l.DebugContext(ctx, "debug with attributes", "number", 123)

	l.InfoContext(ctx, "info with no attributes")
	l.InfoContext(ctx, "info with attributes", "number", 123)

	l = l.With("attr", "abc")
	l.InfoContext(ctx, "new info with no attributes")
	l.InfoContext(ctx, "new info with attributes", "number", 123)

	l.ErrorContext(ctx, "error with no attributes")
	l.ErrorContext(ctx, "error with attributes", "number", 123)

	// Output:
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=DEBUG msg=\"debug with no attributes\""}
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=DEBUG msg=\"debug with attributes\" number=123"}
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=INFO msg=\"info with no attributes\""}
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=INFO msg=\"info with attributes\" number=123"}
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=INFO msg=\"new info with no attributes\" attr=abc"}
	// {"request_id":"abcdefghijklmnop","severity":"NORMAL","message":"level=INFO msg=\"new info with attributes\" number=123 attr=abc"}
	// {"request_id":"abcdefghijklmnop","severity":"ERROR","message":"level=ERROR msg=\"error with no attributes\" attr=abc"}
	// {"request_id":"abcdefghijklmnop","severity":"ERROR","message":"level=ERROR msg=\"error with attributes\" number=123 attr=abc"}
}
