package httputil_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogMiddleware(t *testing.T) {
	logOutput := &bytes.Buffer{}
	logger := slogutil.New(&slogutil.Config{
		Output: logOutput,
		Format: slogutil.FormatJSON,
	})

	mw := httputil.NewLogMiddleware(logger, slog.LevelInfo)
	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l, ok := slogutil.LoggerFromContext(ctx)
		require.True(t, ok)

		l.InfoContext(ctx, "test", "attr", 123)

		_, err := io.WriteString(w, testBody)
		require.NoError(t, err)
	}))

	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequest(http.MethodGet, testPath, nil).WithContext(ctx)

	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testBody, w.Body.String())

	lines := bytes.Split(logOutput.Bytes(), []byte("\n"))

	// This includes an empty line at the end.
	require.Len(t, lines, 4)

	for i, line := range lines {
		if i == 3 && len(line) == 0 {
			continue
		}

		var obj map[string]any
		err := json.Unmarshal(line, &obj)
		require.NoError(t, err)

		assert.NotEmpty(t, "INFO", obj["msg"])
		assert.Equal(t, "INFO", obj["level"])
		assert.Equal(t, http.MethodGet, obj["method"])
		assert.Equal(t, testPath, obj["request_uri"])

		switch i {
		case 1:
			assert.Equal(t, float64(123), obj["attr"])
		case 2:
			assert.Equal(t, float64(http.StatusOK), obj["code"])

			// Make sure that the "elapsed" attribute is printed consistently.
			elapsedStr := testutil.RequireTypeAssert[string](t, obj["elapsed"])
			assert.Regexp(t, `[0-9.]+[a-zµ]+`, elapsedStr)
		}
	}
}

// testResponseWriter is a response writer that does nothing.
type testResponseWriter struct {
	// hdr is the header returned by [testResponseWriter.Header].
	hdr http.Header
}

// type check
var _ http.ResponseWriter = testResponseWriter{}

// Header returns [rw.hdr].
func (rw testResponseWriter) Header() (hdr http.Header) { return rw.hdr }

// Write does nothing and returns len(b), nil.
func (testResponseWriter) Write(b []byte) (_ int, _ error) { return len(b), nil }

// WriteHeader does nothing.
func (testResponseWriter) WriteHeader(_ int) {}

func BenchmarkLogMiddleware(b *testing.B) {
	w := testResponseWriter{
		hdr: http.Header{},
	}

	ctx := context.Background()
	r := httptest.NewRequest(http.MethodGet, testPath, nil).WithContext(ctx)

	b.Run("enabled", func(b *testing.B) {
		logHdlr := slogutil.NewLevelHandler(slog.LevelInfo, slog.DiscardHandler)
		mw := httputil.NewLogMiddleware(slog.New(logHdlr), slog.LevelInfo)
		h := mw.Wrap(httputil.HealthCheckHandler)

		b.ReportAllocs()
		for b.Loop() {
			h.ServeHTTP(w, r)
		}
	})

	b.Run("disabled", func(b *testing.B) {
		mw := httputil.NewLogMiddleware(slogutil.NewDiscardLogger(), slog.LevelInfo)
		h := mw.Wrap(httputil.HealthCheckHandler)

		b.ReportAllocs()
		for b.Loop() {
			h.ServeHTTP(w, r)
		}
	})

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkLogMiddleware
	//	BenchmarkLogMiddleware/enabled
	//	BenchmarkLogMiddleware/enabled-16         	 1000000	      1622 ns/op	     128 B/op	       6 allocs/op
	//	BenchmarkLogMiddleware/disabled
	//	BenchmarkLogMiddleware/disabled-16        	 1423170	       811.7 ns/op	      88 B/op	       4 allocs/op
}
