package httputil

import (
	"net/http"
	"net/http/pprof"
	"net/url"
)

// PprofBasePath is the default base path used by [RoutePprof].
const PprofBasePath = "/debug/pprof/"

// RoutePprof adds all pprof handlers to r under the paths within
// [PprofBasePath].
func RoutePprof(r Router) {
	// See also profileSupportsDelta in src/net/http/pprof/pprof.go.
	routes := []struct {
		handler http.Handler
		method  string
		pattern string
	}{{
		handler: http.HandlerFunc(pprof.Index),
		pattern: "/",
		method:  http.MethodGet,
	}, {
		handler: pprof.Handler("allocs"),
		pattern: "/allocs",
		method:  http.MethodGet,
	}, {
		handler: pprof.Handler("block"),
		pattern: "/block",
		method:  http.MethodGet,
	}, {
		handler: http.HandlerFunc(pprof.Cmdline),
		pattern: "/cmdline",
		method:  http.MethodGet,
	}, {
		handler: pprof.Handler("goroutine"),
		pattern: "/goroutine",
		method:  http.MethodGet,
	}, {
		handler: pprof.Handler("heap"),
		pattern: "/heap",
		method:  http.MethodGet,
	}, {
		handler: pprof.Handler("mutex"),
		pattern: "/mutex",
		method:  http.MethodGet,
	}, {
		handler: http.HandlerFunc(pprof.Profile),
		pattern: "/profile",
		method:  http.MethodGet,
	}, {
		handler: http.HandlerFunc(pprof.Symbol),
		pattern: "/symbol",
		method:  http.MethodGet,
	}, {
		// NOTE:  The /symbol API can accept both GET and POST queries.
		handler: http.HandlerFunc(pprof.Symbol),
		pattern: "/symbol",
		method:  http.MethodPost,
	}, {
		handler: pprof.Handler("threadcreate"),
		pattern: "/threadcreate",
		method:  http.MethodGet,
	}, {
		handler: http.HandlerFunc(pprof.Trace),
		pattern: "/trace",
		method:  http.MethodGet,
	}}

	for _, route := range routes {
		pattern, err := url.JoinPath(PprofBasePath, route.pattern)
		if err != nil {
			// Generally shouldn't happen, as the list of patterns is static.
			panic(err)
		}

		r.Handle(route.method+" "+pattern, route.handler)
	}
}
