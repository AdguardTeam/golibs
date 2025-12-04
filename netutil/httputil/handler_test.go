package httputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

// TODO(a.garipov):  Add error test once package testhttp lands in golibs.
func TestPlainTextHandler_ServeHTTP(t *testing.T) {
	h := httputil.PlainTextHandler(testBody)

	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testBody, w.Body.String())
}

func TestPanicHandler(t *testing.T) {
	const testErr errors.Error = "foo"

	h := httputil.PanicHandler(testErr)
	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	assert.PanicsWithValue(t, testErr, func() {
		h.ServeHTTP(w, r)
	})
}
