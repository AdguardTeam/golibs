package randutil

import (
	"math/rand/v2"
	"sync"
)

// Reader is a ChaCha8-based cryptographically strong random number reader that
// is safe for concurrent use.
type Reader struct {
	// mu protects reader.
	mu     *sync.Mutex
	reader *rand.ChaCha8
}

// NewReader returns a new properly initialized *Reader seeded with the given
// seed.
func NewReader(seed [32]byte) (r *Reader) {
	return &Reader{
		mu:     &sync.Mutex{},
		reader: rand.NewChaCha8(seed),
	}
}

// Read generates len(p) random bytes and writes them into p.  It always returns
// len(p) and a nil error.  It's safe for concurrent use.
func (r *Reader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.reader.Read(p)
}
