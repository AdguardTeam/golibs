//go:build go1.21

package slogutil_test

import (
	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func ExampleNew_default() {
	l := slogutil.New(&slogutil.Config{
		Verbose: true,
	})

	l.Info("test info")
	l.Debug("test debug")

	// Output:
	// INFO test info
	// DEBUG test debug
}

func ExampleNew_json() {
	l := slogutil.New(&slogutil.Config{
		Format:  slogutil.FormatJSON,
		Verbose: true,
	})

	l.Info("test info")
	l.Debug("test debug")

	l.WithGroup("test_group").Info("group test info", "time", "too late")
	l.WithGroup("test_group").Debug("group test debug", "time", "too late")

	// Output:
	// {"level":"INFO","msg":"test info"}
	// {"level":"DEBUG","msg":"test debug"}
	// {"level":"INFO","msg":"group test info","test_group":{"time":"too late"}}
	// {"level":"DEBUG","msg":"group test debug","test_group":{"time":"too late"}}
}

func ExampleNew_text() {
	l := slogutil.New(&slogutil.Config{
		Format:  slogutil.FormatText,
		Verbose: true,
	})

	l.Info("test info")
	l.Debug("test debug")

	l.WithGroup("test_group").Info("group test info", "time", "too late")
	l.WithGroup("test_group").Debug("group test debug", "time", "too late")

	// Output:
	// level=INFO msg="test info"
	// level=DEBUG msg="test debug"
	// level=INFO msg="group test info" test_group.time="too late"
	// level=DEBUG msg="group test debug" test_group.time="too late"
}
