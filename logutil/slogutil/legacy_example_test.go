package slogutil_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	aglog "github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func ExampleAdGuardLegacyHandler() {
	// Use a buffer to remove the PID#GID prefix.
	output := &bytes.Buffer{}

	aglog.SetOutput(output)
	aglog.SetFlags(0)
	aglog.SetLevel(aglog.DEBUG)

	h := slogutil.NewAdGuardLegacyHandler(slog.LevelDebug)
	l := slog.New(h)

	l.Debug("debug with no attributes")
	l.Debug("debug with attributes", "number", 123)

	l.Info("info with no attributes")
	l.Info("info with attributes", "number", 123)

	l = l.With(slogutil.KeyPrefix, "hdlr")
	l.Info("info with prefix")

	l.Warn("warning with two prefixes (bad!)", slogutil.KeyPrefix, "bad")

	// Remove the PID#GID prefix for a reproducible example.
	for _, line := range strings.Split(output.String(), "\n") {
		_, line, _ = strings.Cut(line, " ")
		fmt.Println(line)
	}

	// Output:
	// [debug] debug with no attributes
	// [debug] debug with attributes number=123
	// [info] info with no attributes
	// [info] info with attributes number=123
	// [info] hdlr: info with prefix
	// [debug] legacy logger: got prefix "bad" in record for logger with prefix "hdlr"
	// [info] hdlr: warning: warning with two prefixes (bad!)
}
