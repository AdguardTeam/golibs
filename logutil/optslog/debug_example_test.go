package optslog_test

import (
	"context"

	"github.com/AdguardTeam/golibs/logutil/optslog"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func Example_debug() {
	l := slogutil.New(&slogutil.Config{
		Level:  slogutil.LevelInfo,
		Format: slogutil.FormatText,
	})

	ctx := context.Background()

	optslog.Debug1(ctx, l, "one arg; not printed", "1", 1)
	optslog.Debug2(ctx, l, "two args; not printed", "1", 1, "2", 2)
	optslog.Debug3(ctx, l, "three args; not printed", "1", 1, "2", 2, "3", 3)
	optslog.Debug4(ctx, l, "four args; not printed", "1", 1, "2", 2, "3", 3, "4", 4)

	l = slogutil.New(&slogutil.Config{
		Level:  slogutil.LevelDebug,
		Format: slogutil.FormatText,
	})

	optslog.Debug1(ctx, l, "one arg; printed", "1", 1)
	optslog.Debug2(ctx, l, "two args; printed", "1", 1, "2", 2)
	optslog.Debug3(ctx, l, "three args; printed", "1", 1, "2", 2, "3", 3)
	optslog.Debug4(ctx, l, "four args; printed", "1", 1, "2", 2, "3", 3, "4", 4)

	// Output:
	// level=DEBUG msg="one arg; printed" 1=1
	// level=DEBUG msg="two args; printed" 1=1 2=2
	// level=DEBUG msg="three args; printed" 1=1 2=2 3=3
	// level=DEBUG msg="four args; printed" 1=1 2=2 3=3 4=4
}
