package httputil_test

import (
	"net/http"
	"testing"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakenet/fakehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDRoundTripper_RoundTrip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		header   string
		reqID    requestid.ID
		generate bool
	}{{
		name:     "no_generation_empty_context",
		generate: false,
	}, {
		name:     "no_generation_reqid_in_context",
		generate: false,
		reqID:    testRequestID,
	}, {
		name:     "generation_empty_context",
		generate: true,
	}, {
		name:   "request_with_reqid",
		reqID:  testRequestID,
		header: testRequestIDStr,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var header string
			transport := &fakehttp.RoundTripper{
				OnRoundTrip: func(req *http.Request) (resp *http.Response, err error) {
					header = req.Header.Get(httphdr.XRequestID)

					return nil, nil
				},
			}

			conf := &httputil.RequestIDRoundTripperConfig{
				Transport: transport,
				Generate:  tc.generate,
			}

			rt := httputil.NewRequestIDRoundTripper(conf)

			ctx := testutil.ContextWithTimeout(t, testTimeout)
			if tc.reqID != "" {
				ctx = requestid.ContextWithRequestID(ctx, tc.reqID)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
			require.NoError(t, err)

			req.Header.Set(httphdr.XRequestID, tc.header)

			_, err = rt.RoundTrip(req)
			require.NoError(t, err)

			if tc.generate {
				assert.NotEmpty(t, header)
			} else if tc.reqID == "" {
				assert.Empty(t, header)
			}
		})
	}
}

func BenchmarkRequestIDRoundTripper(b *testing.B) {
	transport := &fakehttp.RoundTripper{
		OnRoundTrip: func(req *http.Request) (resp *http.Response, err error) {
			return nil, nil
		},
	}

	rt := httputil.NewRequestIDRoundTripper(
		&httputil.RequestIDRoundTripperConfig{
			Transport: transport,
		},
	)

	ctx := testutil.ContextWithTimeout(b, testTimeout)
	ctx = requestid.ContextWithRequestID(ctx, testRequestID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	require.NoError(b, err)

	// Warmup pool.
	_, err = rt.RoundTrip(req)
	require.NoError(b, err)

	b.ReportAllocs()
	for b.Loop() {
		_, err = rt.RoundTrip(req)
		require.NoError(b, err)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkRequestIDRoundTripper
	//	BenchmarkRequestIDRoundTripper-24    	 6023694	       204.6 ns/op	      16 B/op	       1 allocs/op
}
