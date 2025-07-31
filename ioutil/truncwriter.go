package ioutil

import "io"

// TruncatedWriter is an [io.Writer] that writes up to a certain limit of bytes
// to its underlying writer and then ignores the rest.
type TruncatedWriter struct {
	w      io.Writer
	limit  uint64
	offset uint64
}

// NewTruncatedWriter returns a new truncated writer.  It wraps w so that it
// writes up to limit bytes and then ignores the rest.
func NewTruncatedWriter(w io.Writer, limit uint64) (tw *TruncatedWriter) {
	return &TruncatedWriter{
		w:     w,
		limit: limit,
	}
}

// type check
var _ io.Writer = (*TruncatedWriter)(nil)

// Write implements the [io.Writer] interface for *TruncatedWriter.  n is always
// len(b).
func (w *TruncatedWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	remaining := w.limit - w.offset
	if remaining == 0 {
		return n, nil
	}

	idx := min(uint64(n), remaining)

	// TODO(e.burkov): As the actual number of written bytes could be less then
	// idx, consider returning this actual number.
	_, err = w.w.Write(b[:idx])

	w.offset += idx

	return n, err
}
