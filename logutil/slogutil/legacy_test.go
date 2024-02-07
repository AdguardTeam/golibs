package slogutil_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	aglog "github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/stretchr/testify/require"
)

func BenchmarkAdGuardLegacyHandler_Handle(b *testing.B) {
	aglog.SetOutput(io.Discard)
	h := slogutil.NewAdGuardLegacyHandler(slog.LevelInfo)

	ctx := context.Background()
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	r.AddAttrs(
		slog.Int("int", 123),
		slog.String("string", "abc"),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errSink = h.Handle(ctx, r)
	}

	require.NoError(b, errSink)

	// Most recent results, on a ThinkPad X13 with a Ryzen Pro 7 CPU:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/logutil/slogutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkAdGuardLegacyHandler_Handle-16    	  100437	     13357 ns/op	     361 B/op	      14 allocs/op
}
