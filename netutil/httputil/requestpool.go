package httputil

import (
	"io"
	"net/http"

	"github.com/AdguardTeam/golibs/syncutil"
)

// RequestPool is a specialized pool that can be used by [http.RoundTripper]
// implementations to safely reuse requests and their headers.
//
// NOTE:  It is currently unsafe to reuse [http.Request]s in [http.RoundTrip]
// implementations due to races in stdlib.  See https://golang.org/issues/81445.
type RequestPool struct {
	bodies *syncutil.Pool[reqPoolRespBody]
	hdrs   *HeaderPool
	reqs   *syncutil.Pool[http.Request]
}

// NewRequestPool returns a new, properly initialized *RequestPool.
func NewRequestPool() (p *RequestPool) {
	return &RequestPool{
		bodies: syncutil.NewPool(func() (b *reqPoolRespBody) {
			return &reqPoolRespBody{}
		}),
		hdrs: NewHeaderPool(),
		reqs: syncutil.NewPool(func() (r *http.Request) {
			return &http.Request{}
		}),
	}
}

// Get returns a shallow clone of req that should be returned to the pool with
// either [RequestPool.Put] or [RequestPool.PutOnClose].  req must not be nil.
//
// If the clone is made only to change the headers, [RequestPool.GetWithHeader]
// is recommended instead.
func (p *RequestPool) Get(req *http.Request) (clone *http.Request) {
	clone = p.reqs.Get()

	ctx := req.Context()
	CopyRequestTo(ctx, clone, req)

	//lint:ignore SA1019 Set the deprecated Cancel field to nil to prevent races
	// when reusing requests.  Cancel is automatically set by the default HTTP
	// client for compatibility, but some [http.RoundTripper]s don't support it.
	clone.Cancel = nil

	return clone
}

// GetWithHeader is like [RequestPool.Get] but also clones the header and sets
// it in clone.  clone and hpe should be returned to the pool with either
// [RequestPool.Put] or [RequestPool.PutOnClose].  req must not be nil.
func (p *RequestPool) GetWithHeader(req *http.Request) (clone *http.Request, hpe *HeaderPoolEntry) {
	clone = p.Get(req)

	// NOTE:  CopyRequestTo creates a shallow copy, so in order to modify the
	// outgoing request without too many new allocations, a new header from the
	// pool must be used.
	hpe = p.hdrs.Get(clone.Header)
	clone.Header = hpe.Header()

	return clone, hpe
}

// Put returns req and, if non-nil, hpe, to the pools.  req must not be nil.
func (p *RequestPool) Put(req *http.Request, hpe *HeaderPoolEntry) {
	p.reqs.Put(req)
	if hpe != nil {
		p.hdrs.Put(hpe)
	}
}

// PutOnClose wraps resp.Body to return req and hpe (if present) to the pool
// when closed.  If resp is nil, it assumes that there was an error and returns
// req and hpe to the pool immediately.  See the [http.RoundTripper.RoundTrip]
// documentation:
//
// “[…] RoundTrip must return err == nil if it obtained a response, regardless
// of the response's HTTP status code.  A non-nil err should be reserved for
// failure to obtain a response.  […]”
//
// req must not be nil.  If resp is not nil, resp.Body must not be nil and
// resp.Request must not be used after calling resp.Body.Close.
//
// NOTE:  There is a data race in the Go standard library when the client's
// redirect check returns an error and the uerr helper reads
// [http.Response.Request] to obtain the URL which is immediately discarded.
//
// TODO(a.garipov):  File an issue about said race.
func (p *RequestPool) PutOnClose(resp *http.Response, req *http.Request, hpe *HeaderPoolEntry) {
	if resp == nil {
		p.Put(req, hpe)

		return
	}

	body := p.bodies.Get()
	body.rc = resp.Body
	body.pool = p
	body.req = req
	body.hpe = hpe

	resp.Body = body
}

// reqPoolRespBody is a wrapper around a body of an HTTP response that returns
// the reusable data to its pools once the body is closed.
//
// TODO(a.garipov):  Consider returning [http.ErrBodyReadAfterClose] from
// [reqPoolRespBody.Read] after [reqPoolRespBody.Close].
type reqPoolRespBody struct {
	// rc is the original response body.  It must not be nil.
	rc io.ReadCloser

	// hpe is the header-pool entry to return to its pool, if any.
	hpe *HeaderPoolEntry

	// pool is the pool where [reqPoolRespBody.req] and [reqPoolRespBody.hpe]
	// are returned.  It must not be nil unless returned to a pool.
	pool *RequestPool

	// req is the request to return to the pool.  It must not be nil unless
	// returned to a pool.
	req *http.Request
}

// type check
var _ io.ReadCloser = (*reqPoolRespBody)(nil)

// Read implements the [io.ReadCloser] interface for *reqPoolRespBody.
//
// TODO(a.garipov):  Consider returning [http.ErrBodyReadAfterClose] if this
// body is currently closed.
func (b *reqPoolRespBody) Read(buf []byte) (n int, err error) {
	return b.rc.Read(buf)
}

// Close implements the [io.ReadCloser] interface for *reqPoolRespBody.  The
// body must not be used after a call to Close until it's reinitialized by the
// pool.
func (b *reqPoolRespBody) Close() (err error) {
	err = b.rc.Close()

	b.pool.Put(b.req, b.hpe)

	bodies := b.pool.bodies

	// Zero out all reference-type fields before putting the body back into the
	// pool so as to not confuse the garbage collector with pointer cycles.
	b.rc = nil
	b.hpe = nil
	b.pool = nil
	b.req = nil
	bodies.Put(b)

	return err
}
