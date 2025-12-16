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
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testTimeout is common timeout for tests.
	testTimeout = 1 * time.Second

	// testRequestID is a common request ID for tests.
	testRequestID requestid.ID = "abcdefghijklmnop"
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

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	ctx = requestid.ContextWithRequestID(ctx, testRequestID)

	wg := &sync.WaitGroup{}
	for i := range numGoroutine {
		wg.Go(func() {
			hybridLogger.InfoContext(ctx, "test message", "i", i, "attr", "abc")
			textLogger.InfoContext(ctx, "test message", "i", i, "attr", "abc")
		})
	}

	wg.Wait()

	hybridOutputStrings := strings.Split(hybridOutput.String(), "\n")
	require.Len(t, hybridOutputStrings, numGoroutine+1)

	textOutputStrings := strings.Split(textOutput.String(), "\n")
	require.Len(t, textOutputStrings, numGoroutine+1)

	slices.Sort(hybridOutputStrings)
	slices.Sort(textOutputStrings)

	var (
		prefix = `{"request_id":"` +
			string(testRequestID) +
			`","severity":"NORMAL","message":"level=INFO msg=\"test message\" `
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
	ctx = requestid.ContextWithRequestID(ctx, testRequestID)

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
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/logutil/slogutil
	//	cpu: Apple M3
	//	BenchmarkJSONHybridHandler_Handle-8   	 2292276	       515.2 ns/op	      80 B/op	       2 allocs/op
}
