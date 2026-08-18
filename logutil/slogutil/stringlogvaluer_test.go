package slogutil_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/testutil"
)

func Benchmark_StringLogValuer(b *testing.B) {
	const (
		logMsg  = "msg"
		attrKey = "str"
	)

	stringer := testutil.NewStringer()

	data := []byte("a")
	stringer.OnString = func() (s string) {
		return string(data)
	}

	l := slogutil.New(&slogutil.Config{
		Format: slogutil.FormatJSONHybrid,
		Output: io.Discard,
		Level:  slog.LevelInfo,
	})

	b.Run("log_valuer_disabled_log_level", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			l.Debug(logMsg, attrKey, slogutil.NewStringLogValuer(stringer))
		}
	})

	b.Run("log_valuer_enabled_log_level", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			l.Info(logMsg, attrKey, slogutil.NewStringLogValuer(stringer))
		}
	})

	b.Run("string_call_disabled_log_level", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			l.Debug(logMsg, attrKey, stringer.String())
		}
	})

	b.Run("string_call_enabled_log_level", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			l.Info(logMsg, attrKey, stringer.String())
		}
	})

	// Most recent results:
	//
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/AdGuardDNS/internal/agdslog
	//	cpu: AMD Ryzen AI 7 PRO 350 w/ Radeon 860M
	//	Benchmark_StringLogValuer/log_valuer_disabled_log_level-16              247650633                4.612 ns/op           0 B/op        0 allocs/op
	//	Benchmark_StringLogValuer/log_valuer_enabled_log_level-16                1511684               779.5 ns/op            72 B/op        2 allocs/op
	//	Benchmark_StringLogValuer/string_call_disabled_log_level-16             47299754                23.82 ns/op           16 B/op        1 allocs/op
	//	Benchmark_StringLogValuer/string_call_enabled_log_level-16               1544602               781.5 ns/op            88 B/op        3 allocs/op
}
