package optslog_test

import (
	"context"
	"io"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/optslog"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	const N = 1_000

	l := slogutil.New(&slogutil.Config{
		Output: io.Discard,
	})

	ctx := context.Background()

	const testMessage = "test message"

	testCases := []struct {
		f    func()
		name string
	}{{
		f:    func() { optslog.Debug1(ctx, l, testMessage, "1", 1) },
		name: "Debug1",
	}, {
		f:    func() { optslog.Debug2(ctx, l, testMessage, "1", 1, "2", 2) },
		name: "Debug2",
	}, {
		f:    func() { optslog.Debug3(ctx, l, testMessage, "1", 1, "2", 2, "3", 3) },
		name: "Debug3",
	}, {
		f:    func() { optslog.Debug4(ctx, l, testMessage, "1", 1, "2", 2, "3", 3, "4", 4) },
		name: "Debug4",
	}, {
		f:    func() { optslog.Trace1(ctx, l, testMessage, "1", 1) },
		name: "Trace1",
	}, {
		f:    func() { optslog.Trace2(ctx, l, testMessage, "1", 1, "2", 2) },
		name: "Trace2",
	}, {
		f:    func() { optslog.Trace3(ctx, l, testMessage, "1", 1, "2", 2, "3", 3) },
		name: "Trace3",
	}, {
		f:    func() { optslog.Trace4(ctx, l, testMessage, "1", 1, "2", 2, "3", 3, "4", 4) },
		name: "Trace4",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := testing.AllocsPerRun(N, tc.f)
			assert.Zero(t, got)
		})
	}
}
