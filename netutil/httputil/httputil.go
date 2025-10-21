// Package httputil contains common constants, functions, and types for working
// with HTTP.
package httputil

import (
	"context"
	"net/http"
)

// Router is the interface for HTTP routers, such as [http.ServeMux].
type Router interface {
	// Handle registers the handler for the given pattern.
	Handle(pattern string, h http.Handler)
}

// type check
var _ Router = (*http.ServeMux)(nil)

// RouterFunc is a functional implementation of the [Router] interface.
type RouterFunc func(pattern string, h http.Handler)

// type check
var _ Router = RouterFunc(nil)

// Handle implements the [Router] interface for RouterFunc.
func (f RouterFunc) Handle(pattern string, h http.Handler) {
	f(pattern, h)
}

// CopyRequestTo is an optimized version of [http.Request.WithContext] that uses
// compiler optimizations to allow reducing allocations with a pool.  ctx, dst,
// and src must not be nil.
//
// See https://github.com/golang/go/issues/68501#issuecomment-2234069762.
func CopyRequestTo(ctx context.Context, dst, src *http.Request) {
	// NOTE:  Just putting this into the code doesn't work most of the time.
	// Most likely because it gets harder for the compiler to inline.
	*dst = *src.WithContext(ctx)
}
