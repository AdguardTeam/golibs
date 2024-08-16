package container

// RingBuffer is the generic implementation of ring buffer data structure.
type RingBuffer[T any] struct {
	buf  []T
	cur  uint
	full bool
}

// NewRingBuffer initializes a new ring buffer with the given size.
func NewRingBuffer[T any](size uint) (rb *RingBuffer[T]) {
	return &RingBuffer[T]{
		buf: make([]T, size),
	}
}

// Push adds an element to the buffer and sets the current position to the next
// element.
func (rb *RingBuffer[T]) Push(e T) {
	if len(rb.buf) == 0 {
		return
	}

	rb.buf[rb.cur] = e
	rb.cur = (rb.cur + 1) % uint(cap(rb.buf))
	if rb.cur == 0 {
		rb.full = true
	}
}

// Current returns the element at the current position.  It returns zero value
// of T if rb is nil or empty.
func (rb *RingBuffer[T]) Current() (e T) {
	if rb == nil || len(rb.buf) == 0 {
		return e
	}

	return rb.buf[rb.cur]
}

// Range calls f for each element of the buffer starting from the current
// position until f returns false.
func (rb *RingBuffer[T]) Range(f func(T) (cont bool)) {
	before, after := rb.splitCur()

	for _, e := range before {
		if !f(e) {
			return
		}
	}

	for _, e := range after {
		if !f(e) {
			return
		}
	}
}

// ReverseRange calls f for each element of the buffer in reverse order ending
// with the current position until f returns false.
func (rb *RingBuffer[T]) ReverseRange(f func(T) (cont bool)) {
	before, after := rb.splitCur()

	for i := len(after) - 1; i >= 0; i-- {
		if !f(after[i]) {
			return
		}
	}

	for i := len(before) - 1; i >= 0; i-- {
		if !f(before[i]) {
			return
		}
	}
}

// splitCur splits the buffer in two, before and after current position in
// chronological order.  If buffer is not full, after is nil.
func (rb *RingBuffer[T]) splitCur() (before, after []T) {
	if len(rb.buf) == 0 {
		return nil, nil
	}

	cur := rb.cur
	if !rb.full {
		return rb.buf[:cur], nil
	}

	return rb.buf[cur:], rb.buf[:cur]
}

// Len returns a length of the buffer.
func (rb *RingBuffer[T]) Len() (l uint) {
	if !rb.full {
		return rb.cur
	}

	return uint(cap(rb.buf))
}

// Clear clears the buffer.
func (rb *RingBuffer[T]) Clear() {
	rb.full = false
	rb.cur = 0
}
