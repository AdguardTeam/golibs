// Package bytecache provides a simple LRU cache implementation that takes size
// into account.
package bytecache

import "github.com/c2h5oh/datasize"

// Interface is the byte cache interface.  All methods must be safe for
// concurrent use.
type Interface interface {
	// Clear all data and statistics.
	Clear()

	// Del deletes the cached data.
	Del(key []byte)

	// Get returns the cached data.  val must not be modified.
	Get(key []byte) (val []byte, ok bool)

	// Set sets the data.  key and val must not be modified after calling Set.
	Set(key, val []byte)

	// Stats returns the current cache statistics.  The statistics are not
	// guaranteed to be consistent.
	Stats() (s Stats)
}

// Stats represents statistics for the cache.
type Stats struct {
	// Size is the current size of the byte cache.
	Size datasize.ByteSize

	// Count is the current number of items in the cache.
	Count uint64

	// Hit is the number of cache hits.
	Hit uint64

	// Miss is the number of cache misses.
	Miss uint64
}

// Empty is an [Interface] implementation that does nothing.
type Empty struct{}

// type check
var _ Interface = Empty{}

// Clear implements the [Interface] interface for Empty.
func (Empty) Clear() {}

// Get implements the [Interface] interface for Empty.  It always returns nil
// and false.
func (Empty) Get(key []byte) (val []byte, ok bool) { return nil, false }

// Del implements the [Interface] interface for Empty.
func (Empty) Del(key []byte) {}

// Set implements the [Interface] interface for Empty.
func (Empty) Set(key, val []byte) {}

// Stats implements the [Interface] interface for Empty.  s is always empty.
func (Empty) Stats() (s Stats) { return s }
