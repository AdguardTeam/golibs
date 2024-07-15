package httputil_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogMiddleware(t *testing.T) {
	logOutput := &bytes.Buffer{}
	l := slogutil.New(&slogutil.Config{
		Output: logOutput,
		Format: slogutil.FormatJSON,
	})

	mw := httputil.NewLogMiddleware(l, slog.LevelInfo)
	h := mw.Wrap(httputil.HealthCheckHandler)

	w := httptest.NewRecorder()
	ctx := testutil.ContextWithTimeout(t, testTimeout)
	r := httptest.NewRequest(http.MethodGet, testPath, nil).WithContext(ctx)

	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testBody, w.Body.String())

	lines := bytes.Split(logOutput.Bytes(), []byte("\n"))

	// This includes an empty line at the end.
	require.Len(t, lines, 3)

	for i, line := range lines {
		if i == 2 && len(line) == 0 {
			continue
		}

		var obj map[string]any
		err := json.Unmarshal(line, &obj)
		require.NoError(t, err)

		assert.Equal(t, "INFO", obj["level"])
		assert.Equal(t, http.MethodGet, obj["method"])
		assert.Equal(t, testPath, obj["request_uri"])

		if i == 1 {
			assert.Equal(t, float64(http.StatusOK), obj["code"])
		}
	}
}
