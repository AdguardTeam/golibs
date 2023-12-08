package syncutil

import "context"

// Semaphore is the semaphore interface.
type Semaphore interface {
	// Acquire gets the resource, will block until the resource can be acquired.
	// ctx is used for cancellation.
	Acquire(ctx context.Context) (err error)

	// Release the resource, never blocks.
	Release()
}

// EmptySemaphore is a semaphore that has no limit.
type EmptySemaphore struct{}

// type check
var _ Semaphore = EmptySemaphore{}

// Acquire implements the [Semaphore] interface for EmptySemaphore.  It always
// returns nil.
func (EmptySemaphore) Acquire(_ context.Context) (err error) { return nil }

// Release implements the [Semaphore] interface for EmptySemaphore.
func (EmptySemaphore) Release() {}

// unit is a convenient alias for struct{}.
type unit = struct{}

// ChanSemaphore is a channel-based semaphore.
//
// It must be initialized with [NewChanSemaphore].
type ChanSemaphore struct {
	c chan unit
}

// type check
var _ Semaphore = (*ChanSemaphore)(nil)

// Acquire implements the [Semaphore] interface for *ChanSemaphore.
func (c *ChanSemaphore) Acquire(ctx context.Context) (err error) {
	select {
	case c.c <- unit{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release implements the [Semaphore] interface for *ChanSemaphore.
func (c *ChanSemaphore) Release() {
	select {
	case <-c.c:
	default:
	}
}

// NewChanSemaphore returns a new *ChanSemaphore with the provided maximum
// resource number.
func NewChanSemaphore(maxRes uint) (c *ChanSemaphore) {
	return &ChanSemaphore{
		c: make(chan unit, maxRes),
	}
}
