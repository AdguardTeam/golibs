//go:build go1.21

package slogutil_test

import (
	"log/slog"
	"os"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// removeTime is a [slog.HandlerOptions.ReplaceAttr] function that removes the
// "time" attribute.
func removeTime(groups []string, a slog.Attr) (res slog.Attr) {
	if len(groups) == 0 && a.Key == "time" {
		return slog.Attr{}
	}

	return a
}

func ExampleJSONHybridHandler() {
	h := slogutil.NewJSONHybridHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		// Use removeTime to make the example reproducible.
		ReplaceAttr: removeTime,
	})
	l := slog.New(h)

	l.Debug("debug with no attributes")
	l.Debug("debug with attributes", "number", 123)

	l.Info("info with no attributes")
	l.Info("info with attributes", "number", 123)

	l = l.With("attr", "abc")
	l.Info("new info with no attributes")
	l.Info("new info with attributes", "number", 123)

	// Output:
	// {"level":"DEBUG","msg":"debug with no attributes"}
	// {"level":"DEBUG","msg":"debug with attributes; attrs: number=123"}
	// {"level":"INFO","msg":"info with no attributes"}
	// {"level":"INFO","msg":"info with attributes; attrs: number=123"}
	// {"level":"INFO","msg":"new info with no attributes; attrs: attr=abc"}
	// {"level":"INFO","msg":"new info with attributes; attrs: attr=abc number=123"}
}
