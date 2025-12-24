package httputil

import (
	"net/http"

	"github.com/AdguardTeam/golibs/syncutil"
)

// HeaderPool allows reducing allocations when making many clones of
// [http.Header] maps.
type HeaderPool struct {
	entries *syncutil.Pool[HeaderPoolEntry]
}

// HeaderPoolEntry is an entry in the [HeaderPool].
type HeaderPoolEntry struct {
	header http.Header
	values []string
}

// Header returns the cloned header.
func (e *HeaderPoolEntry) Header() (h http.Header) {
	return e.header
}

// headerMapLenEst is the estimate of how many headers are usually sent in an
// HTTP request or response.  It is based on the kinds of requests modern
// browsers send over HTTPS.
const headerMapLenEst = 16

// NewHeaderPool returns a new properly initialized *HeaderPool.
func NewHeaderPool() (p *HeaderPool) {
	return &HeaderPool{
		entries: syncutil.NewPool(func() (e *HeaderPoolEntry) {
			return &HeaderPoolEntry{
				header: make(http.Header, headerMapLenEst),
				values: make([]string, 0, headerMapLenEst),
			}
		}),
	}
}

// Get returns an entry that contains a deep clone of orig.  e should be
// returned to the pool by using [HeaderPool.Put].
func (p *HeaderPool) Get(orig http.Header) (e *HeaderPoolEntry) {
	e = p.entries.Get()
	clear(e.header)
	clear(e.values)
	e.values = e.values[:0]

	valuesStart := 0
	for k, values := range orig {
		if values == nil {
			// NOTE:  Preserve nil values, because [httputil.ReverseProxy]
			// distinguishes between nil and zero-length header values.
			e.header[k] = nil

			continue
		}

		l := len(values)
		e.values = append(e.values, values...)
		e.header[k] = e.values[valuesStart : valuesStart+l : valuesStart+l]
		valuesStart += l
	}

	return e
}

// Put returns e to the pool for later reuse.
func (p *HeaderPool) Put(e *HeaderPoolEntry) {
	p.entries.Put(e)
}
