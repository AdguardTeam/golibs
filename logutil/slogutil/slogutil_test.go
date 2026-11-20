package slogutil_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/testutil/fakefmt"
	"github.com/stretchr/testify/assert"
)

func TestPrint_disabledLevel(t *testing.T) {
	const N = 1_000

	l := slogutil.New(&slogutil.Config{
		Output: io.Discard,
		Level:  slog.LevelInfo,
	})

	ctx := context.Background()

	const (
		lvl = slog.LevelDebug
		msg = "test message"
	)

	text := "line 1\nline 2\nline 3"
	textBytes := []byte(text)

	// The stringer panics if String is called.
	stringer := fakefmt.NewStringer()

	testCases := []struct {
		f    func()
		name string
	}{{
		f:    func() { slogutil.PrintLines(ctx, l, lvl, msg, text) },
		name: "PrintLines",
	}, {
		f:    func() { slogutil.PrintByteLines(ctx, l, lvl, msg, textBytes) },
		name: "PrintByteLines",
	}, {
		f:    func() { slogutil.PrintStack(ctx, l, lvl) },
		name: "PrintStack",
	}, {
		f:    func() { slogutil.PrintStringerLines(ctx, l, lvl, msg, stringer) },
		name: "PrintStringerLines",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := testing.AllocsPerRun(N, tc.f)
			assert.Zero(t, got)
		})
	}
}
