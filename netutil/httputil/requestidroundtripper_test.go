package httputil_test

import (
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
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
		header:   "",
		reqID:    "",
		generate: false,
	}, {
		name:     "no_generation_reqid_in_context",
		header:   "",
		reqID:    testRequestID,
		generate: false,
	}, {
		name:     "generation_empty_context",
		header:   "",
		reqID:    "",
		generate: true,
	}, {
		name:     "request_with_reqid",
		reqID:    testRequestID,
		header:   testRequestIDStr,
		generate: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var header string
			onRoundTrip := func(req *http.Request) (resp *http.Response, err error) {
				header = req.Header.Get(httphdr.XRequestID)

				return nil, nil
			}

			transport := &fakehttp.RoundTripper{
				OnRoundTrip: onRoundTrip,
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
			} else {
				assert.Equal(t, string(tc.reqID), header)
			}
		})
	}
}

func TestRequestIDRoundTripper_race(t *testing.T) {
	t.Parallel()

	const (
		durWaitMax        = testTimeout / 10
		numGoroutine      = 1_000
		numServerRedirect = 10
		numServer         = 10
	)

	srvURLs := make([]string, numServer)
	for i := range numServer {
		s := httptest.NewTestServer(t, httputil.HealthCheckHandler)
		s.Start()
		srvURLs[i] = s.URL
	}

	srvRedirURLs := make([]string, numServerRedirect)
	for i := range numServerRedirect {
		srvURL := srvURLs[i]
		s := httptest.NewTestServer(t, http.RedirectHandler(srvURL, http.StatusTemporaryRedirect))
		s.Start()
		srvRedirURLs[i] = s.URL
	}

	rt := httputil.NewRequestIDRoundTripper(&httputil.RequestIDRoundTripperConfig{
		Transport: http.DefaultTransport,
		Generate:  true,
	})

	c := &http.Client{
		Transport: rt,
		// TODO(a.garipov):  Find ways to improve the performance on CI.
		Timeout: testTimeout * 10,
	}

	wg := &sync.WaitGroup{}
	for i := range numGoroutine {
		wg.Go(func() {
			time.Sleep(rand.N(durWaitMax))
			u := srvRedirURLs[i%numServerRedirect]

			resp, err := c.Get(u)
			require.NotNil(t, resp)
			require.NoError(t, err)

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}

	wg.Wait()
}

// TODO(a.garipov):  Improve tests and benchmarks.
func BenchmarkRequestIDRoundTripper_idPresentError(b *testing.B) {
	transport := &fakehttp.RoundTripper{
		OnRoundTrip: func(req *http.Request) (resp *http.Response, err error) {
			return nil, assert.AnError
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
	require.ErrorIs(b, err, assert.AnError)

	b.ReportAllocs()
	for b.Loop() {
		_, err = rt.RoundTrip(req)
		require.ErrorIs(b, err, assert.AnError)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkRequestIDRoundTripper_idPresentError-24    	 2626470	       407.3 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkRequestIDRoundTripper_noIDWithClose(b *testing.B) {
	var body io.ReadCloser = &struct {
		fakeio.Closer
		fakeio.Reader
	}{
		OnClose: func() (err error) { return nil },
		OnRead:  func(b []byte) (n int, err error) { panic(testutil.UnexpectedCall(b)) },
	}

	respWithBody := &http.Response{
		Body: body,
	}

	transport := &fakehttp.RoundTripper{
		OnRoundTrip: func(_ *http.Request) (resp *http.Response, err error) {
			return respWithBody, nil
		},
	}

	rt := httputil.NewRequestIDRoundTripper(
		&httputil.RequestIDRoundTripperConfig{
			Transport: transport,
			Generate:  true,
		},
	)

	ctx := testutil.ContextWithTimeout(b, testTimeout)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	require.NoError(b, err)

	// Warmup pool.
	_, err = rt.RoundTrip(req)
	require.NoError(b, err)

	err = respWithBody.Body.Close()
	require.NoError(b, err)

	respWithBody.Body = body

	b.ReportAllocs()
	for b.Loop() {
		_, err = rt.RoundTrip(req)
		require.NoError(b, err)

		err = respWithBody.Body.Close()
		require.NoError(b, err)

		respWithBody.Body = body
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil/httputil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkRequestIDRoundTripper_noIDWithClose-24    	 2154786	       510.0 ns/op	      24 B/op	       1 allocs/op
}
