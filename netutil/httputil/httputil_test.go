package httputil_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/stretchr/testify/require"
)

// testTimeout is a common timeout for tests.
const testTimeout = 1 * time.Second

// Common constants for tests.
const (
	testPath = "/health-check"
	testBody = string(httputil.HealthCheckHandler)
)

// Common sinks for benchmark results.
var (
	reqSink *http.Request
)

func BenchmarkCopyRequestTo(b *testing.B) {
	ctx := context.Background()

	reqPool := syncutil.NewPool(func() (r *http.Request) { return &http.Request{} })

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		func() {
			reqSink = reqPool.Get()
			defer func() { reqPool.Put(reqSink) }()

			httputil.CopyRequestTo(ctx, reqSink, r)
		}()
	}

	wantBuf, gotBuf := &bytes.Buffer{}, &bytes.Buffer{}
	err = r.Write(wantBuf)
	require.NoError(b, err)

	err = reqSink.Write(gotBuf)
	require.NoError(b, err)

	require.Equal(b, wantBuf.Bytes(), gotBuf.Bytes())

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkCopyRequestTo-16    	34533667	        31.93 ns/op	       0 B/op	       0 allocs/op
}
