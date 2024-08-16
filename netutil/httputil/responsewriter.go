package httputil

import (
	"cmp"
	"net/http"
)

// Wrapper is a copy of the hidden rwUnwrapper interface from the Go standard
// library.  It is added here for tests, linting, etc.
type Wrapper interface {
	Unwrap() (rw http.ResponseWriter)
}

// CodeRecorderResponseWriter wraps an [http.ResponseWriter] allowing to save
// the response code.
type CodeRecorderResponseWriter struct {
	rw   http.ResponseWriter
	code int
}

// NewCodeRecorderResponseWriter returns a new *CodeRecorderResponseWriter which
// uses the given response writer as its base.
func NewCodeRecorderResponseWriter(rw http.ResponseWriter) (w *CodeRecorderResponseWriter) {
	return &CodeRecorderResponseWriter{
		rw: rw,
	}
}

// Code returns the status code that was set.  It expects that
// [CodeRecorderResponseWriter.SetImplicitSuccess] has already been called.
func (w *CodeRecorderResponseWriter) Code() (code int) {
	return w.code
}

// SetImplicitSuccess should be called after the handler has finished to set the
// status code to [http.StatusOK] if it hadn't been set explicitly.  This can be
// used to detect panics within handlers, as when a handler panics before
// calling w.WriteHeader, SetImplicitSuccess isn't reached, and w.Code returns 0
// and false.
func (w *CodeRecorderResponseWriter) SetImplicitSuccess() {
	w.code = cmp.Or(w.code, http.StatusOK)
}

// type check
var _ Wrapper = (*CodeRecorderResponseWriter)(nil)

// Unwrap implements the [Wrapper] interface for *CodeRecorderResponseWriter.
func (w *CodeRecorderResponseWriter) Unwrap() (rw http.ResponseWriter) {
	return w.rw
}

// type check
var _ http.ResponseWriter = (*CodeRecorderResponseWriter)(nil)

// Header implements [http.ResponseWriter] for *CodeRecorderResponseWriter.
func (w *CodeRecorderResponseWriter) Header() (h http.Header) {
	return w.rw.Header()
}

// Write implements [http.ResponseWriter] for *CodeRecorderResponseWriter.
func (w *CodeRecorderResponseWriter) Write(b []byte) (n int, err error) {
	return w.rw.Write(b)
}

// WriteHeader implements [http.ResponseWriter] for *CodeRecorderResponseWriter.
func (w *CodeRecorderResponseWriter) WriteHeader(code int) {
	w.code = code

	w.rw.WriteHeader(code)
}
