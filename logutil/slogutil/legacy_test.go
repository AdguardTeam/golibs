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

	var err error
	b.ReportAllocs()
	for b.Loop() {
		err = h.Handle(ctx, r)
	}

	require.NoError(b, err)

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/logutil/slogutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkAdGuardLegacyHandler_Handle
	//	BenchmarkAdGuardLegacyHandler_Handle-16    	 1000000	      2070 ns/op	     208 B/op	       9 allocs/op
}
