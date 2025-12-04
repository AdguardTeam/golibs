package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCodeRecorderResponseWriter(t *testing.T) {
	w := httputil.NewCodeRecorderResponseWriter(httptest.NewRecorder())
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	h := httputil.HealthCheckHandler
	h.ServeHTTP(w, r)
	w.SetImplicitSuccess()

	assert.Equal(t, http.StatusOK, w.Code())

	rr := testutil.RequireTypeAssert[*httptest.ResponseRecorder](t, w.Unwrap())

	assert.Equal(t, testBody, rr.Body.String())
}
