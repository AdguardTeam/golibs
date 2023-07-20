// Package pprofutil contains utilities for pprof HTTP handlers.
package pprofutil

import (
	"net/http"
	"net/http/pprof"
	"net/url"
)

// BasePath is the default base path used by [RoutePprof].
//
// TODO(a.garipov): Consider adding the ability to configure the base path.
const BasePath = "/debug/pprof/"

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

// RoutePprof adds all pprof handlers to r under the paths within [BasePath].
func RoutePprof(r Router) {
	// See also profileSupportsDelta in src/net/http/pprof/pprof.go.
	routes := []struct {
		handler http.Handler
		pattern string
	}{{
		handler: http.HandlerFunc(pprof.Index),
		pattern: "/",
	}, {
		handler: pprof.Handler("allocs"),
		pattern: "/allocs",
	}, {
		handler: pprof.Handler("block"),
		pattern: "/block",
	}, {
		handler: http.HandlerFunc(pprof.Cmdline),
		pattern: "/cmdline",
	}, {
		handler: pprof.Handler("goroutine"),
		pattern: "/goroutine",
	}, {
		handler: pprof.Handler("heap"),
		pattern: "/heap",
	}, {
		handler: pprof.Handler("mutex"),
		pattern: "/mutex",
	}, {
		handler: http.HandlerFunc(pprof.Profile),
		pattern: "/profile",
	}, {
		handler: http.HandlerFunc(pprof.Symbol),
		pattern: "/symbol",
	}, {
		handler: pprof.Handler("threadcreate"),
		pattern: "/threadcreate",
	}, {
		handler: http.HandlerFunc(pprof.Trace),
		pattern: "/trace",
	}}

	for _, route := range routes {
		pattern, err := url.JoinPath(BasePath, route.pattern)
		if err != nil {
			// Generally shouldn't happen, as the list of patterns is static.
			panic(err)
		}

		r.Handle(pattern, route.handler)
	}
}
