package optslog_test

import (
	"context"

	"github.com/AdguardTeam/golibs/logutil/optslog"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

func Example_trace() {
	l := slogutil.New(&slogutil.Config{
		Level:  slogutil.LevelInfo,
		Format: slogutil.FormatText,
	})

	ctx := context.Background()

	optslog.Trace1(ctx, l, "one arg; not printed", "1", 1)
	optslog.Trace2(ctx, l, "two args; not printed", "1", 1, "2", 2)
	optslog.Trace3(ctx, l, "three args; not printed", "1", 1, "2", 2, "3", 3)
	optslog.Trace4(ctx, l, "four args; not printed", "1", 1, "2", 2, "3", 3, "4", 4)

	l = slogutil.New(&slogutil.Config{
		Level:  slogutil.LevelTrace,
		Format: slogutil.FormatText,
	})

	optslog.Trace1(ctx, l, "one arg; printed", "1", 1)
	optslog.Trace2(ctx, l, "two args; printed", "1", 1, "2", 2)
	optslog.Trace3(ctx, l, "three args; printed", "1", 1, "2", 2, "3", 3)
	optslog.Trace4(ctx, l, "four args; printed", "1", 1, "2", 2, "3", 3, "4", 4)

	// Output:
	// level=TRACE msg="one arg; printed" 1=1
	// level=TRACE msg="two args; printed" 1=1 2=2
	// level=TRACE msg="three args; printed" 1=1 2=2 3=3
	// level=TRACE msg="four args; printed" 1=1 2=2 3=3 4=4
}
