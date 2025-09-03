package httputil

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/AdguardTeam/golibs/timeutil"
)

// LogMiddleware adds a logger using [slogutil.ContextWithLogger] and logs the
// starts and ends of queries at a given level.
type LogMiddleware struct {
	attrPool *syncutil.Pool[[]slog.Attr]
	reqPool  *syncutil.Pool[http.Request]
	rwPool   *syncutil.Pool[CodeRecorderResponseWriter]
	logger   *slog.Logger
	lvl      slog.Level
}

// logMwAttrNum is the number of attributes used by the logger set by
// [LogMiddleware].
const logMwAttrNum = 4

// NewLogMiddleware returns a new *LogMiddleware with l as the base logger.
func NewLogMiddleware(l *slog.Logger, lvl slog.Level) (mw *LogMiddleware) {
	return &LogMiddleware{
		attrPool: syncutil.NewSlicePool[slog.Attr](logMwAttrNum),
		reqPool: syncutil.NewPool(func() (r *http.Request) {
			return &http.Request{}
		}),
		rwPool: syncutil.NewPool(func() (rw *CodeRecorderResponseWriter) {
			return &CodeRecorderResponseWriter{}
		}),
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

		attrsPtr := mw.attrsSlicePtr(r)
		defer mw.attrPool.Put(attrsPtr)

		logHdlr := mw.logger.Handler().WithAttrs(*attrsPtr)
		l := slog.New(logHdlr)
		ctx := slogutil.ContextWithLogger(r.Context(), l)

		nextReq := mw.reqPool.Get()
		defer mw.reqPool.Put(nextReq)

		CopyRequestTo(ctx, nextReq, r)

		rw := mw.rwPool.Get()
		defer mw.rwPool.Put(rw)

		rw.Reset(w)

		l.Log(ctx, mw.lvl, "started")
		defer mw.logFinished(ctx, l, rw, startTime)

		h.ServeHTTP(rw, nextReq)
		rw.SetImplicitSuccess()
	}

	return http.HandlerFunc(f)
}

// logFinished is called at the end of handling of a query.
func (mw *LogMiddleware) logFinished(
	ctx context.Context,
	l *slog.Logger,
	rw *CodeRecorderResponseWriter,
	startTime time.Time,
) {
	if l.Enabled(ctx, mw.lvl) {
		l.Log(
			ctx,
			mw.lvl,
			"finished",
			"code", rw.code,
			"elapsed", timeutil.Duration(time.Since(startTime)),
		)
	}
}

// attrsSlicePtr returns a pointer to a slice with the attributes from the
// request set.  Callers should defer returning attrsPtr back to the pool.
func (mw *LogMiddleware) attrsSlicePtr(r *http.Request) (attrsPtr *[]slog.Attr) {
	attrsPtr = mw.attrPool.Get()

	attrs := *attrsPtr

	// Optimize bounds checking.
	_ = attrs[logMwAttrNum-1]

	attrs[0] = slog.String("host", r.Host)
	attrs[1] = slog.String("method", r.Method)
	attrs[2] = slog.String("raddr", r.RemoteAddr)
	attrs[3] = slog.String("request_uri", r.RequestURI)

	return attrsPtr
}
