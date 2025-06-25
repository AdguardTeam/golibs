package randutil

import (
	"math/rand/v2"
	"sync"
)

// LockedSource is an implementation of [rand.Source] that is concurrency-safe.
type LockedSource struct {
	// mu protects src.
	mu  *sync.Mutex
	src rand.Source
}

// NewLockedSource returns new properly initialized *LockedSource.
func NewLockedSource(src rand.Source) (s *LockedSource) {
	return &LockedSource{
		mu:  &sync.Mutex{},
		src: src,
	}
}

// type check
var _ rand.Source = (*LockedSource)(nil)

// Uint64 implements the [rand.Source] interface for *LockedSource.
func (s *LockedSource) Uint64() (r uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.src.Uint64()
}
