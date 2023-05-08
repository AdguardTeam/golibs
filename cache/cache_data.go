package cache

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type onDeleteType func(key []byte, val []byte)

type cache struct {
	items map[string]*item

	// LRU: for removing the least recently used item on reaching cache size limit
	// Note: slows down Get() due to an additional work with pointers
	// Example: [ (sentinel) <-> item1 <-> item2 <-> (sentinel) ]
	// When the item is accessed, it's moved to the end of the list.
	usage     listItem
	usageLock sync.Mutex

	lock          sync.RWMutex
	size          uint       // current size in bytes (keys+values)
	promotions    chan *item // channel to buffer item promotions
	promoterClose chan interface{}

	conf Config

	// stats:
	miss           int32 // number of misses
	hit            int32 // number of hits
	promotionFails int32 // number of failed promotions due to congestion
}

type item struct {
	key   []byte
	value []byte
	used  listItem
}

const maxUint = (1 << (unsafe.Sizeof(uint(0)) * 8)) - 1

func newCache(conf Config) *cache {
	c := cache{}
	c.Clear()
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conf = conf
	if c.conf.MaxSize == 0 {
		c.conf.MaxSize = maxUint
	}
	if c.conf.MaxCount == 0 {
		c.conf.MaxCount = maxUint
	}
	if c.conf.MaxElementSize == 0 {
		c.conf.MaxElementSize = c.conf.MaxSize
	}
	if c.conf.MaxElementSize > c.conf.MaxSize {
		c.conf.MaxElementSize = c.conf.MaxSize
	}
	return &c
}

func (c *cache) Clear() {
	c.lock.Lock()
	c.items = make(map[string]*item)
	c.size = 0
	c.lock.Unlock()

	c.usageLock.Lock()
	listInit(&c.usage)
	c.usageLock.Unlock()

	atomic.StoreInt32(&c.hit, 0)
	atomic.StoreInt32(&c.miss, 0)
	atomic.StoreInt32(&c.promotionFails, 0)
}

// Set value
func (c *cache) Set(key []byte, val []byte) bool {
	addSize := uint(len(key) + len(val))
	if addSize > c.conf.MaxElementSize {
		return false // too large data
	}

	it := item{}
	it.key = key
	it.value = val

	c.lock.Lock()
	defer c.lock.Unlock()
	c.usageLock.Lock()
	defer c.usageLock.Unlock()

	if !c.conf.EnableLRU &&
		(c.size+addSize > c.conf.MaxSize || uint(len(c.items)) == c.conf.MaxCount) {
		return false // cache is full
	}

	for c.size+addSize > c.conf.MaxSize || uint(len(c.items)) == c.conf.MaxCount {
		first := listFirst(&c.usage)
		it := (*item)(structPtr(unsafe.Pointer(first), unsafe.Offsetof(item{}.used)))
		c.size -= uint(len(it.key) + len(it.value))
		listUnlink(first)
		delete(c.items, string(it.key))

		if c.conf.OnDelete != nil {
			c.conf.OnDelete(it.key, it.value)
		}
	}

	if c.conf.EnableLRU {
		listAppend(&it.used, listLast(&c.usage))
	}

	it2, exists := c.items[string(key)]
	if exists {
		listUnlink(&it2.used)
		c.size -= uint(len(it2.key) + len(it2.value))
	}
	c.items[string(key)] = &it
	c.size += addSize

	return exists
}

// Get value
func (c *cache) Get(key []byte) []byte {

	c.lock.RLock()
	val, ok := c.items[string(key)]
	c.lock.RUnlock()

	if ok && c.conf.EnableLRU {
		c.promote(val)
	}

	if !ok {
		atomic.AddInt32(&c.miss, 1)
		return nil
	}
	atomic.AddInt32(&c.hit, 1)
	return val.value
}

// Del - delete element
func (c *cache) Del(key []byte) {
	c.lock.Lock()
	it, ok := c.items[string(key)]
	if !ok {
		c.lock.Unlock()
		return
	}
	listUnlink(&it.used)
	c.size -= uint(len(it.key) + len(it.value))
	delete(c.items, string(key))
	c.lock.Unlock()
}

// GetStats - get counters
func (c *cache) Stats() Stats {
	s := Stats{}
	s.Count = len(c.items)
	s.Size = int(c.size)
	s.Hit = int(atomic.LoadInt32(&c.hit))
	s.Miss = int(atomic.LoadInt32(&c.miss))
	s.PromotionFails = int(atomic.LoadInt32(&c.promotionFails))
	return s
}

// promote - promotes LRU values
func (c *cache) promote(item *item) {
	select {
	case c.promotions <- item:
	default:
		atomic.AddInt32(&c.promotionFails, 1)
	}
}

// promoter - executes promotions on usage
func (c *cache) promoter() {
	for {
		select {
		case item := <-c.promotions:
			c.lock.RLock()
			_, ok := c.items[string(item.key)]
			if ok {
				c.usageLock.Lock()
				listUnlink(&item.used)
				listAppend(&item.used, listLast(&c.usage))
				c.usageLock.Unlock()
			} else {
				atomic.AddInt32(&c.promotionFails, 1)
			}
			c.lock.RUnlock() // unlock last to prevent modification by Set()
		case <-c.promoterClose:
			return
		}
	}
}
