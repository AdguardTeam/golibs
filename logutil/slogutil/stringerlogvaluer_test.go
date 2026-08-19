package slogutil_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/testutil/fakefmt"
)

func Benchmark_StringerLogValuer(b *testing.B) {
	const (
		logMsg  = "msg"
		attrKey = "str"
	)

	stringer := fakefmt.NewStringer()

	data := []byte("a")
	stringer.OnString = func() (s string) {
		return string(data)
	}

	l := slogutil.New(&slogutil.Config{
		Format: slogutil.FormatJSONHybrid,
		Output: io.Discard,
		Level:  slog.LevelInfo,
	})

	benchCases := []struct {
		logArg   any
		name     string
		logLevel slog.Level
	}{{
		name:     "log_valuer_disabled_log_level",
		logLevel: slog.LevelDebug,
		logArg:   slogutil.NewStringerLogValuer(stringer),
	}, {
		name:     "log_valuer_enabled_log_level",
		logLevel: slog.LevelInfo,
		logArg:   slogutil.NewStringerLogValuer(stringer),
	}, {
		name:     "string_call_disabled_log_level",
		logLevel: slog.LevelDebug,
		logArg:   stringer.String(),
	}, {
		name:     "string_call_enabled_log_level",
		logLevel: slog.LevelInfo,
		logArg:   stringer.String(),
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				l.Log(context.Background(), bc.logLevel, logMsg, attrKey, bc.logArg)
			}
		})
	}

	// Most recent results:
	//
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/logutil/slogutil
	//	cpu: AMD Ryzen AI 7 PRO 350 w/ Radeon 860M
	//	Benchmark_StringerLogValuer/log_valuer_disabled_log_level-16            233588494                4.882 ns/op           0 B/op        0 allocs/op
	//	Benchmark_StringerLogValuer/log_valuer_enabled_log_level-16              1509172               790.9 ns/op            72 B/op        2 allocs/op
	//	Benchmark_StringerLogValuer/string_call_disabled_log_level-16           250471549                4.766 ns/op           0 B/op        0 allocs/op
	//	Benchmark_StringerLogValuer/string_call_enabled_log_level-16             1547262               763.6 ns/op            72 B/op        2 allocs/op
}
