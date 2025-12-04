package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

func TestServerHeaderMiddleware(t *testing.T) {
	const srvHdr = "Test/1.0"
	mw := httputil.ServerHeaderMiddleware(srvHdr)

	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	h := mw.Wrap(httputil.HealthCheckHandler)
	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, srvHdr, w.Header().Get(httphdr.Server))
	assert.Equal(t, testBody, w.Body.String())
}
