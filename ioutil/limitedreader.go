package ioutil

import (
	"fmt"
	"io"

	"github.com/AdguardTeam/golibs/mathutil"
)

// LimitError is returned when the Limit is reached.
type LimitError struct {
	// Limit is the limit that triggered the error.
	Limit uint64
}

// type check
var _ error = (*LimitError)(nil)

// Error implements the error interface for *LimitError.
func (err *LimitError) Error() (msg string) {
	return fmt.Sprintf("cannot read more than %d bytes", err.Limit)
}

// limitedReader is a wrapper for io.Reader that has a reading limit.
type limitedReader struct {
	r     io.Reader
	limit uint64
	n     uint64
}

// type check
var _ io.Reader = (*limitedReader)(nil)

// Read implements the [io.Reader] interface for *limitedReader.
func (lr *limitedReader) Read(p []byte) (n int, err error) {
	if lr.n == 0 {
		return 0, &LimitError{
			Limit: lr.limit,
		}
	}

	// TODO(a.garipov): Use the builtin in Go 1.21.
	l := mathutil.Min(uint64(len(p)), lr.n)
	p = p[:l]

	n, err = lr.r.Read(p)
	lr.n -= uint64(n)

	return n, err
}

// LimitReader returns an io.Reader that reads up to n bytes.  Once that limit
// is reached, [ErrLimit] is returned from limited.Read.  limited.Read is not
// safe for concurrent use.  n must be non-negative.
func LimitReader(r io.Reader, n uint64) (limited io.Reader) {
	return &limitedReader{
		r:     r,
		limit: n,
		n:     n,
	}
}
