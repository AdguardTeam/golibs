package syncutil

import (
	"fmt"
	"sync"
)

// Pool is the strongly typed version of [sync.Pool] that manages pointers
// to T.
type Pool[T any] struct {
	pool *sync.Pool
}

// NewPool returns a new strongly typed pool.  newFunc must not be nil.
func NewPool[T any](newFunc func() (v *T)) (p *Pool[T]) {
	if newFunc == nil {
		panic(fmt.Errorf("nil newFunc in NewPool"))
	}

	return &Pool[T]{
		pool: &sync.Pool{
			New: func() (v any) { return newFunc() },
		},
	}
}

// NewSlicePool is a helper for constructing pools with pointers to slices of a
// type with the given length.
func NewSlicePool[T any](l int) (p *Pool[[]T]) {
	return NewPool(func() (v *[]T) {
		s := make([]T, l)

		return &s
	})
}

// Get selects an arbitrary item from the pool, removes it from the pool, and
// returns it to the caller.
//
// See [sync.Pool.Get].
func (p *Pool[T]) Get() (v *T) {
	return p.pool.Get().(*T)
}

// Put adds v to the pool.
//
// See [sync.Pool.Put].
func (p *Pool[T]) Put(v *T) {
	p.pool.Put(v)
}
