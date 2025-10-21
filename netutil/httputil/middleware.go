package httputil

import (
	"net/http"
	"slices"
)

// Middleware is a common HTTP middleware interface.
type Middleware interface {
	// Wrap returns a new handler that wraps the specified handler.
	Wrap(handler http.Handler) (wrapped http.Handler)
}

// Wrap is a helper function that wraps h with the specified middlewares.
// middlewares are wrapped, and thus called, in the same order in which they
// were specified.  That is, the first middleware will be the first to receive
// the request, and so on.
func Wrap(h http.Handler, middlewares ...Middleware) (wrapped http.Handler) {
	wrapped = h

	// Wrap the handler with the middlewares in the reverse order.  This way the
	// middleware that was specified first is also the first to receive the
	// request.
	for _, mw := range slices.Backward(middlewares) {
		wrapped = mw.Wrap(wrapped)
	}

	return wrapped
}

// MiddlewareFunc is a function that implements the [Middleware] interface.
type MiddlewareFunc func(h http.Handler) (wrapped http.Handler)

// type check
var _ Middleware = MiddlewareFunc(nil)

// Wrap implements the [Middleware] interface for MiddlewareFunc.
func (f MiddlewareFunc) Wrap(h http.Handler) (wrapped http.Handler) {
	return f(h)
}

// type check
var _ MiddlewareFunc = PassThrough

// PassThrough is a [MiddlewareFunc] that returns h as-is.
func PassThrough(h http.Handler) (wrapped http.Handler) {
	return h
}
