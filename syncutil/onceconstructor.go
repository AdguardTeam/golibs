package syncutil

import "sync"

// OnceConstructor initializes a value for one key only once.
//
// TODO(a.garipov):  Add benchmarks.
type OnceConstructor[K comparable, V any] struct {
	loaders *sync.Map
	new     func(k K) (v V)
}

// NewOnceConstructor returns a new properly initialized *OnceConstructor that
// uses newFunc to construct a value for the given key.
func NewOnceConstructor[K comparable, V any](newFunc func(k K) (v V)) (c *OnceConstructor[K, V]) {
	return &OnceConstructor[K, V]{
		loaders: &sync.Map{},
		new:     newFunc,
	}
}

// Get returns a value for the given key.  If a value isn't available, it waits
// until it is.
func (c *OnceConstructor[K, V]) Get(key K) (v V) {
	// Step 1.  The fast track: check if there is already a value present.
	loaderVal, inited := c.loaders.Load(key)
	if inited {
		return loaderVal.(func() (v V))()
	}

	// Step 2.  Allocate a done channel and create a function that waits for one
	// single initialization.  Use the one returned from LoadOrStore regardless
	// of whether it's this one.
	//
	// TODO(a.garipov):  See if sync.Once or a similar stdlib API can be used.
	var cached V
	done := make(chan struct{}, 1)
	done <- struct{}{}
	loaderVal, _ = c.loaders.LoadOrStore(key, func() (loaded V) {
		_, ok := <-done
		if ok {
			// The only real receive.  Initialize the cached value and close the
			// channel so that other goroutines receive the same value.
			cached = c.new(key)
			close(done)
		}

		return cached
	})

	return loaderVal.(func() (v V))()
}
