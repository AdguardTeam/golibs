// Package fakehttp contains fake implementations of interfaces from packages
// net/http and github.com/AdguardTeam/golibs/netutil/httputil.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic(testutil.UnexpectedCall(arg1, arg2))
package fakehttp

import "net/http"

// ResponseWriter is an [http.ResponseWriter] for tests.
type ResponseWriter struct {
	OnHeader      func() (hdr http.Header)
	OnWrite       func(b []byte) (n int, err error)
	OnWriteHeader func(code int)
}

// type check
var _ http.ResponseWriter = (*ResponseWriter)(nil)

// Header implements the [http.ResponseWriter] interface for *ResponseWriter.
func (rw *ResponseWriter) Header() (hdr http.Header) {
	return rw.OnHeader()
}

// Write implements the [http.ResponseWriter] interface for *ResponseWriter.
func (rw *ResponseWriter) Write(b []byte) (n int, err error) {
	return rw.OnWrite(b)
}

// WriteHeader implements the [http.ResponseWriter] interface for
// *ResponseWriter.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.OnWriteHeader(code)
}

// RoundTripper is an [http.RoundTripper] for tests.
type RoundTripper struct {
	OnRoundTrip func(req *http.Request) (resp *http.Response, err error)
}

// type check
var _ http.RoundTripper = (*RoundTripper)(nil)

// RoundTrip implements the [http.RoundTripper] interface for *RoundTripper.
func (r *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return r.OnRoundTrip(req)
}
