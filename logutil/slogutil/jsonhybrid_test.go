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
	for i := range numGoroutine {
		wg.Go(func() {
			hybridLogger.Info("test message", "i", i, "attr", "abc")
			textLogger.Info("test message", "i", i, "attr", "abc")
		})
	}

	wg.Wait()

	hybridOutputStrings := strings.Split(hybridOutput.String(), "\n")
	require.Len(t, hybridOutputStrings, numGoroutine+1)

	textOutputStrings := strings.Split(textOutput.String(), "\n")
	require.Len(t, textOutputStrings, numGoroutine+1)

	slices.Sort(hybridOutputStrings)
	slices.Sort(textOutputStrings)

	const (
		prefix = `{"severity":"NORMAL","message":"level=INFO msg=\"test message\" `
		suffix = `"}`
	)

	for i := range numGoroutine {
		textString := textOutputStrings[i]
		expectedString := strings.Replace(textString, `level=INFO msg="test message" `, "", 1)

		gotString := hybridOutputStrings[i]
		gotString = strings.TrimPrefix(gotString, prefix)
		gotString = strings.TrimSuffix(gotString, suffix)

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
	//	BenchmarkJSONHybridHandler_Handle
	//	BenchmarkJSONHybridHandler_Handle-16    	 1000000	      1453 ns/op	      48 B/op	       1 allocs/op
}
