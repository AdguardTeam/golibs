package bytecache

import (
	"cmp"
	"math"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/c2h5oh/datasize"
)

// DefaultConfig defines the configuration for the default byte cache.
type DefaultConfig struct {
	// OnDelete, if non-nil, is called after an element has been evicted during
	// [Default.Set].  OnDelete must not call any methods on cache.
	OnDelete OnDeleteFunc

	// MaxSize is the maximum size of the cache for both keys and values.  Zero
	// means no limit.
	MaxSize datasize.ByteSize

	// MaxElementSize is the maximum size of an element, including the key and
	// the value.  Elements larger than this are dropped.  Zero means equal to
	// [DefaultConfig.MaxSize].  It must be no greater than
	// [DefaultConfig.MaxSize] unless the latter is zero.
	MaxElementSize datasize.ByteSize

	// MaxCount is the maximum number of elements.  Zero means no limit.
	MaxCount uint64
}

// OnDeleteFunc is a type for callback functions called after an element has
// been deleted automatically.  The function should not be long-running.  The
// function must not retain key and val.  The function must not call any methods
// on the cache.
//
// TODO(a.garipov):  Consider making a larger interface to make the cache more
// observable.
type OnDeleteFunc func(key, val []byte)

// Default is the default byte cache implementation.  It is safe for concurrent
// use.
type Default struct {
	// mu protects all mutable fields below.  It must not be nil.
	mu *sync.Mutex

	// items is the map of all items in the byte cache.  It must not be nil.
	items map[string]*item

	// head is the head of the least-recently used (LRU) list of elements.  All
	// elements must be present in items.
	head *item

	// onDelete is called on element eviction.  It must not be nil.
	onDelete OnDeleteFunc

	// maxSize is the maximum size of the byte cache.
	maxSize datasize.ByteSize

	// maxElemSize is the maximum size of an element in the byte cache.
	maxElemSize datasize.ByteSize

	// size is the current size of the byte cache.
	size datasize.ByteSize

	// numHit is the number of cache hits.
	numHit uint64

	// numMiss is the number of cache misses.
	numMiss uint64

	// maxCount is the maximum number of items in the byte cache.
	maxCount uint64
}

// NewDefault creates a new properly initialized default byte cache.  If conf is
// nil, defaults are used.  If conf is not nil, it must be valid.
func NewDefault(conf *DefaultConfig) (c *Default) {
	conf = cmp.Or(conf, &DefaultConfig{})

	c = &Default{
		items: map[string]*item{},
		head:  newListHead(),
		mu:    &sync.Mutex{},

		maxSize:  cmp.Or(conf.MaxSize, math.MaxUint64),
		maxCount: cmp.Or(conf.MaxCount, math.MaxUint64),
	}

	if conf.OnDelete == nil {
		c.onDelete = func(_, _ []byte) {}
	} else {
		c.onDelete = conf.OnDelete
	}

	c.maxElemSize = cmp.Or(conf.MaxElementSize, c.maxSize)

	errors.Check(validate.NoGreaterThan("conf.MaxElementSize", c.maxElemSize, c.maxSize))

	return c
}

// type check
var _ Interface = (*Default)(nil)

// Clear implements the [Interface] interface for *Default.
func (c *Default) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.items)

	c.head = newListHead()

	c.size = 0

	c.numHit = 0
	c.numMiss = 0
}

// Get implements the [Interface] interface for *Default.
func (c *Default) Get(key []byte) (val []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// NOTE:  Here and below, assume the Go compiler optimizes string([]byte)
	// conversions when used directly in a map operation.  See the benchmarks.
	it, ok := c.items[string(key)]
	if ok {
		it.unlink()
		it.insertAfter(c.head.back())

		c.numHit++

		return it.value, true
	}

	c.numMiss++

	return nil, false
}

// Del implements the [Interface] interface for *Default.
func (c *Default) Del(key []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	it, ok := c.items[string(key)]
	if !ok {
		return
	}

	it.unlink()
	c.size -= it.size()
	delete(c.items, string(key))
}

// Set implements the [Interface] interface for *Default.
func (c *Default) Set(key, val []byte) {
	addSize := datasize.ByteSize(len(key) + len(val))
	if addSize > c.maxElemSize {
		// TODO(a.garipov):  Make this and other corner cases more observable.
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	itPrev, ok := c.items[string(key)]
	if ok {
		itPrev.unlink()
		delete(c.items, string(key))
		c.size -= itPrev.size()
	}

	for c.size+addSize > c.maxSize || uint64(len(c.items)) >= c.maxCount {
		it := c.head.front()
		c.size -= it.size()
		it.unlink()
		delete(c.items, string(it.key))
		c.onDelete(it.key, it.value)
	}

	it := &item{
		key:   key,
		value: val,
	}

	it.insertAfter(c.head.back())
	c.items[string(key)] = it
	c.size += addSize
}

// Stats implements the [Interface] interface for *Default.
func (c *Default) Stats() (s Stats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return Stats{
		Size:  c.size,
		Count: uint64(len(c.items)),
		Hit:   c.numHit,
		Miss:  c.numMiss,
	}
}
