package httputil_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/AdguardTeam/golibs/netutil/httputil"
)

// testMiddleware is a fake [httputil.Middleware] implementation for tests.
type testMiddleware struct {
	OnWrap func(h http.Handler) (wrapped http.Handler)
}

// type check
var _ httputil.Middleware = (*testMiddleware)(nil)

// Wrap implements the [Middleware] interface for *testMiddleware.
func (mw *testMiddleware) Wrap(h http.Handler) (wrapped http.Handler) {
	return mw.OnWrap(h)
}

func ExampleWrap() {
	mw1 := &testMiddleware{
		OnWrap: func(h http.Handler) (wrapped http.Handler) {
			f := func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("first middleware start")
				h.ServeHTTP(w, r)
				fmt.Println("first middleware end")
			}

			return http.HandlerFunc(f)
		},
	}

	mw2 := &testMiddleware{
		OnWrap: func(h http.Handler) (wrapped http.Handler) {
			f := func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("second middleware start")
				h.ServeHTTP(w, r)
				fmt.Println("second middleware end")
			}

			return http.HandlerFunc(f)
		},
	}

	var h http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		fmt.Println("handler")
	})

	h = httputil.Wrap(h, mw1, mw2)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, testPath, nil)

	h.ServeHTTP(w, r)

	// Output:
	// first middleware start
	// second middleware start
	// handler
	// second middleware end
	// first middleware end
}
