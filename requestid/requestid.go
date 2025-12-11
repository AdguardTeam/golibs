// Package requestid contains utilities for working with request ids.
package requestid

import (
	"context"
	"math/rand/v2"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/mathutil/randutil"
)

// ctxKey is the type for context keys.
type ctxKey int

// Context key values.
const (
	ctxKeyRequestID ctxKey = iota
)

// ID is the ID of a request.  It is an opaque, randomly generated string.
type ID [16]byte

// FromString converts string s to id.
func FromString(s string) (id ID) {
	copy(id[:], []byte(s))

	return id
}

// New generates new request ID.
func New() (id ID) {
	// #nosec G404 -- We do not need a real random here.
	rng := rand.New(rand.NewChaCha8(randutil.MustNewSeed()))

	return FromString(randutil.StringAlphabet(rng, 16, randutil.AlphabetLowercase))
}

// String implements the [fmt.Stringer] interface for ID.
func (i ID) String() (s string) {
	return string(i[:])
}

// ContextWithRequestID returns a new context with the given ID.
func ContextWithRequestID(parent context.Context, id ID) (ctx context.Context) {
	return context.WithValue(parent, ctxKeyRequestID, id)
}

// IDFromContext returns ID for this request, if any.
func IDFromContext(ctx context.Context) (id ID, ok bool) {
	v := ctx.Value(ctxKeyRequestID)
	if v == nil {
		return id, false
	}

	return v.(ID), true
}

// MustIDFromContext returns ID for this request and panics if there
// is no ID.
func MustIDFromContext(ctx context.Context) (id ID) {
	v := ctx.Value(ctxKeyRequestID)
	if v == nil {
		panic(errors.Error("requestid: no request id in context"))
	}

	return v.(ID)
}
