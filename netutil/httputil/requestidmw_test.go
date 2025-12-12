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
)

func TestRequestIDMiddleware(t *testing.T) {
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	req.Header.Set(httphdr.XRequestID, testStrRequestID)
	mw := httputil.NewRequestIDMiddleware()

	var gotReqID string
	h := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		id := requestid.MustIDFromContext(rCtx)
		gotReqID = id.String()

		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	assert.Equal(t, testStrRequestID, gotReqID)
}

func BenchmarkRequestIDMiddleware(b *testing.B) {
	ctx := testutil.ContextWithTimeout(b, testTimeout)
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	req.Header.Set(httphdr.XRequestID, testStrRequestID)
	mw := httputil.NewRequestIDMiddleware()

	h := mw.Wrap(httputil.HealthCheckHandler)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	for b.Loop() {
		h.ServeHTTP(w, req)
	}

	// Most recent result:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: Apple M3
	//	BenchmarkRequestIDMiddleware-8   	 9917952	       103.6 ns/op	     400 B/op	       4 allocs/op
	//	PASS
	//	ok  	github.com/AdguardTeam/golibs/netutil/httputil	1.566s
}
