package httputil

import (
	"net/http"

	"github.com/AdguardTeam/golibs/httphdr"
)

// ServerHeaderMiddleware adds a Server HTTP header to all responses.
type ServerHeaderMiddleware string

// type check
var _ Middleware = ServerHeaderMiddleware("")

// Wrap implements the [Middleware] interface for *ServerHeaderMiddleware.
func (mw ServerHeaderMiddleware) Wrap(h http.Handler) (wrapped http.Handler) {
	f := func(w http.ResponseWriter, r *http.Request) {
		respHdr := w.Header()
		respHdr.Add(httphdr.Server, string(mw))

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
