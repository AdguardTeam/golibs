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

// Header keys and values for tests.
const (
	testHeaderKey = httphdr.AcceptLanguage
	testHeaderVal = "en"
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

func TestHeaderPoolEntry_Set(t *testing.T) {
	t.Parallel()

	want := http.Header{
		testHeaderKey: {testHeaderVal},
	}

	testCases := []struct {
		in   http.Header
		name string
	}{{
		in:   http.Header{},
		name: "empty",
	}, {
		in: http.Header{
			testHeaderKey: nil,
		},
		name: "nil_value",
	}, {
		in: http.Header{
			testHeaderKey: {testHeaderVal},
		},
		name: "value",
	}, {
		in: http.Header{
			testHeaderKey: {testHeaderVal, "other"},
		},
		name: "many_values",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := httputil.NewHeaderPool()

			e := p.Get(tc.in)
			require.NotNil(t, e)

			e.Set(testHeaderKey, testHeaderVal)

			got := e.Header()
			assert.Equal(t, want, got)
		})
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
	//	BenchmarkHeaderPool-24    	 4706522	       249.8 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkHeaderPoolEntry(b *testing.B) {
	p := httputil.NewHeaderPool()

	// Warmup the pool.
	e := p.Get(testHeader)
	e.Set(testHeaderKey, testHeaderVal)
	p.Put(e)

	b.ReportAllocs()
	for b.Loop() {
		e = p.Get(testHeader)
		e.Set(testHeaderKey, testHeaderVal)
		p.Put(e)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkHeaderPoolEntry-24    	 4787138	       250.2 ns/op	       0 B/op	       0 allocs/op
}
