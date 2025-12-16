package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware(t *testing.T) {
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	req.Header.Set(httphdr.XRequestID, string(testRequestID))
	mw := httputil.NewRequestIDMiddleware()

	var gotReqID requestid.ID
	var ok bool
	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		gotReqID, ok = requestid.IDFromContext(rCtx)
		if !ok {
			gotReqID = ""
		}

		w.WriteHeader(http.StatusOK)
	}))

	require.True(t, t.Run("hdr_contains_request_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Equal(t, testRequestID, gotReqID)
	}))

	require.True(t, t.Run("hdr_does_not_contain_request_id", func(t *testing.T) {
		req.Header.Del(httphdr.XRequestID)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.Empty(t, gotReqID)
	}))
}

func BenchmarkRequestIDMiddleware(b *testing.B) {
	ctx := testutil.ContextWithTimeout(b, testTimeout)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	req.Header.Set(httphdr.XRequestID, string(testRequestID))
	mw := httputil.NewRequestIDMiddleware()

	h := mw.Wrap(httputil.HealthCheckHandler)
	w := httptest.NewRecorder()

	// Warmup pool
	h.ServeHTTP(w, req)

	b.ReportAllocs()
	for b.Loop() {
		h.ServeHTTP(w, req)
	}

	// Most recent result:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: Apple M3
	//	BenchmarkRequestIDMiddleware-8   	11447746	       105.1 ns/op	      91 B/op	       3 allocs/op
}
