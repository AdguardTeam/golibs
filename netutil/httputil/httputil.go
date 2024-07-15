// Package httputil contains common constants, functions, and types for working
// with HTTP.
package httputil

import "net/http"

// Middleware is a common HTTP middleware interface.
type Middleware interface {
	Wrap(handler http.Handler) (wrapped http.Handler)
}

// Wrap is a helper function that attaches the specified middlewares to the
// Handler.  Middlewares will be called in the same order in which they were
// specified.  That is, the first middleware will be the first to receive the
// request, and so on.
func Wrap(h http.Handler, middlewares ...Middleware) (wrapped http.Handler) {
	wrapped = h

	// Wrap the handler with the middlewares in the reverse order.  This way the
	// middleware that was specified first is also the first to receive the
	// request.
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		wrapped = m.Wrap(wrapped)
	}

	return wrapped
}

// Router is the interface for HTTP routers, such as [*http.ServeMux].
type Router interface {
	Handle(pattern string, h http.Handler)
}

// RouterFunc is a functional implementation of the [Router] interface.
type RouterFunc func(pattern string, h http.Handler)

// Handle implements the [Router] interface for RouterFunc.
func (f RouterFunc) Handle(pattern string, h http.Handler) {
	f(pattern, h)
}
