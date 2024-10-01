package slogutil_test

import (
	"context"
	"fmt"
	"log/slog"

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
