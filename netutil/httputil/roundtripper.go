package httputil

import (
	"net/http"

	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/requestid"
	"github.com/AdguardTeam/golibs/syncutil"
)

// RequestIDRoundTripperConfig is a configuration structure for
// [RequestIDRoundTripper].
type RequestIDRoundTripperConfig struct {
	// Transport is a round tripper that will be used to perform request after
	// ID substitution.  It must not be nil.
	Transport http.RoundTripper

	// Generate indicates whether a new request ID should be generated if the
	// context doesn't contain one.
	Generate bool
}

// RequestIDRoundTripper is an implementation of [http.RoundTripper] that puts
// the request ID from the context into the X-Request-ID header.  If the request
// already contains an X-Request-ID header, then the round tripper does nothing.
type RequestIDRoundTripper struct {
	transport http.RoundTripper
	reqPool   *syncutil.Pool[http.Request]
	generate  bool
}

// NewRequestIDRoundTripper returns properly initialized *RequestIDRoundTripper.
// c must be valid.
func NewRequestIDRoundTripper(c *RequestIDRoundTripperConfig) (r *RequestIDRoundTripper) {
	return &RequestIDRoundTripper{
		transport: c.Transport,
		generate:  c.Generate,
		reqPool: syncutil.NewPool(func() (r *http.Request) {
			return &http.Request{}
		}),
	}
}

// type check
var _ http.RoundTripper = (*RequestIDRoundTripper)(nil)

// RoundTrip implements the [http.RoundTripper] for *RequestIDRoundTripper.
func (r *RequestIDRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if req.Header.Get(httphdr.XRequestID) != "" {
		return r.transport.RoundTrip(req)
	}

	ctx := req.Context()

	id, ok := requestid.IDFromContext(ctx)
	if !ok && r.generate {
		id = requestid.New()
		ok = true
	}

	if !ok {
		return r.transport.RoundTrip(req)
	}

	newReq := r.reqPool.Get()
	defer r.reqPool.Put(newReq)

	CopyRequestTo(ctx, newReq, req)
	req = newReq

	req.Header.Set(httphdr.XRequestID, string(id))

	return r.transport.RoundTrip(req)
}
