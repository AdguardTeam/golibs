package httputil

import (
	"io"
	"net/http"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// PlainTextHandler is a simple handler that returns the value of the underlying
// string with a "text/plain" content type.  If there is an error during the
// response writing, it gets the logger from the context, if any, and writes the
// error there at the debug level.
type PlainTextHandler string

// type check
var _ http.Handler = PlainTextHandler("")

// ServeHTTP implements the [http.Handler] interface for PlainTextHandler.
func (text PlainTextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(httphdr.ContentType, "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err := io.WriteString(w, string(text))
	if err != nil {
		ctx := r.Context()
		l, ok := slogutil.LoggerFromContext(ctx)
		if ok {
			l.DebugContext(ctx, "writing plain-text response", slogutil.KeyError, err)
		}
	}
}

// HealthCheckHandler is the common healthcheck HTTP handler that writes the
// text "OK\n" into the response.  These are typically used for the GET
// /health-check HTTP API.
const HealthCheckHandler PlainTextHandler = "OK\n"
