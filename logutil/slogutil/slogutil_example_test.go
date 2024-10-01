package slogutil_test

import (
	"context"
	"log/slog"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func ExampleNew_default() {
	l := slogutil.New(&slogutil.Config{
		Level: slog.LevelDebug,
	})

	l.Info("test info")
	l.Debug("test debug")

	// Output:
	// INFO test info
	// DEBUG test debug
}

func ExampleNew_json() {
	l := slogutil.New(&slogutil.Config{
		Format: slogutil.FormatJSON,
		Level:  slog.LevelDebug,
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
		Format: slogutil.FormatText,
		Level:  slog.LevelDebug,
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

func ExamplePrintLines() {
	text := `A Very Long Text

This is a very long text with many lines.`
	l := slogutil.New(nil)

	ctx := context.Background()
	slogutil.PrintLines(ctx, l, slog.LevelInfo, "my text", text)

	// Output:
	// INFO my text line_num=1 line="A Very Long Text"
	// INFO my text line_num=2 line=""
	// INFO my text line_num=3 line="This is a very long text with many lines."
}

func ExampleNew_trace() {
	l := slogutil.New(&slogutil.Config{
		Format: slogutil.FormatText,
		Level:  slogutil.LevelTrace,
	})

	l.Log(context.Background(), slogutil.LevelTrace, "test trace")
	l.Info("test info")
	l.Debug("test debug")

	// Output:
	// level=TRACE msg="test trace"
	// level=INFO msg="test info"
	// level=DEBUG msg="test debug"
}
