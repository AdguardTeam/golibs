package httputil

import (
	"net/http"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/requestid"
)

// RequestIDRoundTripperConfig is a configuration structure for
// [RequestIDRoundTripper].
type RequestIDRoundTripperConfig struct {
	// Transport is a RoundTripper that will be used to perform request after ID
	// substitution.  It must not be nil.
	Transport http.RoundTripper

	// Generate indicates whether a new request ID should be generated if the
	// context doesn't contain one.
	Generate bool
}

// RequestIDRoundTripper is an implementation of [http.RoundTripper] that puts
// request ID from context to X-Request-ID header.
type RequestIDRoundTripper struct {
	transport http.RoundTripper
	generate  bool
}

// NewRequestIDRoundTripper returns properly initialized *RequestIDRoundTripper.
// c must be valid.
func NewRequestIDRoundTripper(c *RequestIDRoundTripperConfig) (r *RequestIDRoundTripper) {
	return &RequestIDRoundTripper{
		transport: c.Transport,
		generate:  c.Generate,
	}
}

// type check
var _ http.RoundTripper = (*RequestIDRoundTripper)(nil)

// RoundTrip implements the [http.RoundTripper] for *RequestIDRoundTripper.
func (r *RequestIDRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()

	id, ok := requestid.IDFromContext(ctx)
	if !ok && r.generate {
		id = requestid.New()
		ok = true
	}

	if ok {
		req = req.Clone(ctx)
		req.Header.Set(httphdr.XRequestID, id.String())
	}

	return r.transport.RoundTrip(req)
}
