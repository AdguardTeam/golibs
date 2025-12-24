package httputil_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testHeader is the common header for tests.
var testHeader = http.Header{
	"Nil-Header":         nil,
	"Zero-Length-Header": {},
	httphdr.Accept:       {"*/*"},
	httphdr.ContentType:  {"application/json"},
	httphdr.Host:         {"example.com"},
	httphdr.UserAgent:    {"MyHTTPClient/1.0"},
	httphdr.XRequestID:   {testRequestIDStr},
}

func TestHeaderPool(t *testing.T) {
	t.Parallel()

	p := httputil.NewHeaderPool()

	e := p.Get(testHeader)
	require.NotNil(t, e)

	h := e.Header()
	assert.Equal(t, testHeader, h)

	// TODO(a.garipov):  See if there are other safe ways to compare two
	// map-header pointers.
	origPtrStr := fmt.Sprintf("%p", testHeader)
	newPtrStr := fmt.Sprintf("%p", h)
	assert.NotEqual(t, origPtrStr, newPtrStr)

	for k, v := range testHeader {
		switch {
		case v == nil:
			assert.Nilf(t, h[k], "at index %q", k)
		case len(v) == 0:
			assert.NotNilf(t, h[k], "at index %q", k)
			assert.Empty(t, h[k], "at index %q", k)
		default:
			assert.NotSamef(t, &testHeader[k][0], &h[k][0], "at index %q", k)
		}
	}
}

func BenchmarkHeaderPool(b *testing.B) {
	p := httputil.NewHeaderPool()

	// Warmup the pool.
	e := p.Get(testHeader)
	p.Put(e)

	b.ReportAllocs()
	for b.Loop() {
		e = p.Get(testHeader)
		p.Put(e)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkHeaderPool-24    	 8087578	       145.6 ns/op	       0 B/op	       0 allocs/op
}
