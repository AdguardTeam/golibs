package httputil

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
)

// LogMiddleware adds a logger using [slogutil.ContextWithLogger] and logs the
// starts and ends of queries at a given level.
type LogMiddleware struct {
	logger *slog.Logger
	lvl    slog.Level
}

// NewLogMiddleware returns a new *LogMiddleware with l as the base logger.
func NewLogMiddleware(l *slog.Logger, lvl slog.Level) (mw *LogMiddleware) {
	return &LogMiddleware{
		logger: l,
		lvl:    lvl,
	}
}

// type check
var _ Middleware = (*LogMiddleware)(nil)

// Wrap implements the [Middleware] interface for *LogMiddleware.
func (mw *LogMiddleware) Wrap(h http.Handler) (wrapped http.Handler) {
	f := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		l := mw.logger.With(
			"raddr", r.RemoteAddr,
			"method", r.Method,
			"host", r.Host,
			"request_uri", r.RequestURI,
		)

		rw := NewCodeRecorderResponseWriter(w)

		ctx := slogutil.ContextWithLogger(r.Context(), l)
		r = r.WithContext(ctx)

		l.Log(ctx, mw.lvl, "started")
		defer func() {
			elapsed := time.Since(startTime)
			l.Log(ctx, mw.lvl, "finished", "code", rw.code, "elapsed", elapsed.String())
		}()

		h.ServeHTTP(rw, r)
		rw.SetImplicitSuccess()
	}

	return http.HandlerFunc(f)
}
