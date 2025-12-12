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
	ctx := testutil.ContextWithTimeout(t, testTimeout)

	var header string
	transport := &fakehttp.RoundTripper{
		OnRoundTrip: func(req *http.Request) (resp *http.Response, err error) {
			header = req.Header.Get(httphdr.XRequestID)

			return nil, nil
		},
	}

	conf := &httputil.RequestIDRoundTripperConfig{
		Transport: transport,
		Generate:  false,
	}

	rt := httputil.NewRequestIDRoundTripper(conf)

	t.Run("without_generation_empty_context", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
		require.NoError(t, err)

		_, err = rt.RoundTrip(req)
		require.NoError(t, err)

		assert.Empty(t, header)
	})

	t.Run("without_generation", func(t *testing.T) {
		ctx = requestid.ContextWithRequestID(ctx, testRequestID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
		require.NoError(t, err)

		_, err = rt.RoundTrip(req)
		require.NoError(t, err)

		assert.Equal(t, testStrRequestID, header)
	})

	conf.Generate = true
	rt = httputil.NewRequestIDRoundTripper(conf)

	t.Run("with_generation_empty_context", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
		require.NoError(t, err)

		_, err = rt.RoundTrip(req)
		require.NoError(t, err)

		assert.NotEmpty(t, header)
	})
}
