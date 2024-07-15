package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

// TODO(a.garipov):  Add error test once package testhttp lands in golibs.
func TestPlainTextHandler_ServeHTTP(t *testing.T) {
	h := httputil.PlainTextHandler(testBody)

	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequest(http.MethodGet, testPath, nil).WithContext(ctx)

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testBody, w.Body.String())
}
