//go:build go1.21

package slogutil_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONHybridHandler_Handle(t *testing.T) {
	var (
		hybridOutput = &bytes.Buffer{}
		textOutput   = &bytes.Buffer{}
	)

	var (
		hybridHdlr = slogutil.NewJSONHybridHandler(hybridOutput, &slog.HandlerOptions{
			ReplaceAttr: slogutil.RemoveTime,
		})
		textHdlr = slog.NewTextHandler(textOutput, &slog.HandlerOptions{
			ReplaceAttr: slogutil.RemoveTime,
		})
	)

	var (
		hybridLogger = slog.New(hybridHdlr)
		textLogger   = slog.New(textHdlr)
	)

	// Test with multiple goroutines to be sure there are no races.
	const numGoroutine = 1_000

	wg := &sync.WaitGroup{}
	for i := 0; i < numGoroutine; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			hybridLogger.Info("test message", "i", i, "attr", "abc")
			textLogger.Info("test message", "i", i, "attr", "abc")
		}(i)
	}

	wg.Wait()

	hybridOutputStrings := strings.Split(hybridOutput.String(), "\n")
	require.Len(t, hybridOutputStrings, numGoroutine+1)

	textOutputStrings := strings.Split(textOutput.String(), "\n")
	require.Len(t, textOutputStrings, numGoroutine+1)

	slices.Sort(hybridOutputStrings)
	slices.Sort(textOutputStrings)

	for i := 0; i < numGoroutine; i++ {
		textString := textOutputStrings[i]
		expectedString := strings.Replace(textString, `level=INFO msg="test message" `, "", 1)

		jsonString := hybridOutputStrings[i]
		gotString := strings.Replace(jsonString, `{"level":"INFO","msg":"test message; attrs: `, "", 1)
		gotString = strings.Replace(gotString, `"}`, "", 1)

		assert.Equal(t, expectedString, gotString)
	}
}

func BenchmarkJSONHybridHandler_Handle(b *testing.B) {
	h := slogutil.NewJSONHybridHandler(io.Discard, nil)

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
	//	BenchmarkJSONHybridHandler_Handle-16       	 1035621	      1246 ns/op	      48 B/op	       1 allocs/op
}
