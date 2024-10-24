package slogutil_test

import (
	"log/slog"
	"os"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func ExampleJSONHybridHandler() {
	h := slogutil.NewJSONHybridHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		// Use slogutil.RemoveTime to make the example reproducible.
		ReplaceAttr: slogutil.RemoveTime,
	})
	l := slog.New(h)

	l.Debug("debug with no attributes")
	l.Debug("debug with attributes", "number", 123)

	l.Info("info with no attributes")
	l.Info("info with attributes", "number", 123)

	l = l.With("attr", "abc")
	l.Info("new info with no attributes")
	l.Info("new info with attributes", "number", 123)

	l.Error("error with no attributes")
	l.Error("error with attributes", "number", 123)

	// Output:
	// {"severity":"NORMAL","message":"level=DEBUG msg=\"debug with no attributes\""}
	// {"severity":"NORMAL","message":"level=DEBUG msg=\"debug with attributes\" number=123"}
	// {"severity":"NORMAL","message":"level=INFO msg=\"info with no attributes\""}
	// {"severity":"NORMAL","message":"level=INFO msg=\"info with attributes\" number=123"}
	// {"severity":"NORMAL","message":"level=INFO msg=\"new info with no attributes\" attr=abc"}
	// {"severity":"NORMAL","message":"level=INFO msg=\"new info with attributes\" number=123 attr=abc"}
	// {"severity":"ERROR","message":"level=ERROR msg=\"error with no attributes\" attr=abc"}
	// {"severity":"ERROR","message":"level=ERROR msg=\"error with attributes\" number=123 attr=abc"}
}
