package ioutil

import (
	"bytes"
	"io"
	"sync"
)

// LockedBuffer is a wrapper around [bytes.Buffer] that makes it safe for
// concurrent use.
type LockedBuffer struct {
	buffer *bytes.Buffer
	mu     sync.Mutex
}

// NewLockedBuffer creates a new instance of *LockedBuffer.  The size of the
// underlying buffer is set to bufSize.
func NewLockedBuffer(bufSize uint) (b *LockedBuffer) {
	return &LockedBuffer{
		buffer: bytes.NewBuffer(make([]byte, 0, bufSize)),
		mu:     sync.Mutex{},
	}
}

// Buffer returns the underlying bytes buffer.
func (b *LockedBuffer) Buffer() (buf *bytes.Buffer) {
	return b.buffer
}

// type check
var _ io.Writer = (*LockedBuffer)(nil)

// Write implements the [io.Writer] interface for *LockedBuffer.
func (b *LockedBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Write(p)
}

// type check
var _ io.Reader = (*LockedBuffer)(nil)

// Read implements the [io.Reader] interface for *LockedBuffer.
func (b *LockedBuffer) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Read(p)
}

// Bytes returns a copy of the buffer contents.
func (b *LockedBuffer) Bytes() (bytes []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Bytes()
}
