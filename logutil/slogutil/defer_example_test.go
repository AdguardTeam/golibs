package slogutil_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
)

func ExampleCloseAndLog() {
	ctx := context.Background()
	l := slogutil.New(&slogutil.Config{
		Level: slog.LevelDebug,
	})

	func() {
		defer slogutil.CloseAndLog(ctx, l, nil, slog.LevelDebug)

		fmt.Println("nil closer:")
	}()

	c := &fakeio.Closer{
		OnClose: func() (err error) {
			return nil
		},
	}

	func() {
		defer slogutil.CloseAndLog(ctx, l, c, slog.LevelDebug)

		fmt.Println("actual closer without error:")
	}()

	c = &fakeio.Closer{
		OnClose: func() (err error) {
			return errors.Error("close failed")
		},
	}

	func() {
		defer slogutil.CloseAndLog(ctx, l, c, slog.LevelDebug)

		fmt.Println("actual closer with error:")
	}()

	// Output:
	//
	// nil closer:
	// actual closer without error:
	// actual closer with error:
	// DEBUG deferred closing err="close failed"
}

func ExamplePrintRecovered() {
	ctx := context.Background()

	output := &bytes.Buffer{}
	l := slogutil.New(&slogutil.Config{
		Output: output,
	})

	func() {
		defer func() { slogutil.PrintRecovered(ctx, l, recover()) }()

		l = l.With("extra", "parameters", "added", "later")

		panic("test value")
	}()

	// Only print the first line, since the stack trace is not reproducible in
	// examples.
	lines := strings.Split(output.String(), "\n")
	fmt.Println(lines[0])

	// Output:
	// ERROR recovered from panic extra=parameters added=later value="test value"
}
