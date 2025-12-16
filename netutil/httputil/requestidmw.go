package httputil

import (
	"net/http"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/syncutil"
)

// RequestIDMiddleware reads request ID from headers and puts it in context.
type RequestIDMiddleware struct {
	reqPool *syncutil.Pool[http.Request]
}

// NewRequestIDMiddleware returns properly initialized RequestIDMiddleware.
func NewRequestIDMiddleware() (r *RequestIDMiddleware) {
	return &RequestIDMiddleware{
		reqPool: syncutil.NewPool(func() (r *http.Request) {
			return &http.Request{}
		}),
	}
}

// type check
var _ Middleware = (*RequestIDMiddleware)(nil)

// Wrap implements the [Middleware] interface for RequestIDMiddleware.
func (mw *RequestIDMiddleware) Wrap(h http.Handler) (wrapped http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqID := r.Header.Get(httphdr.XRequestID)
		if reqID == "" {
			h.ServeHTTP(w, r)

			return
		}

		ctx = requestid.ContextWithRequestID(ctx, requestid.ID(reqID))

		newReq := mw.reqPool.Get()
		defer mw.reqPool.Put(newReq)

		CopyRequestTo(ctx, newReq, r)
		h.ServeHTTP(w, newReq)
	})
}
