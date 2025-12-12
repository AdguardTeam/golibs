// Package requestid contains utilities for working with request ids.
package requestid

import (
	"encoding"
	"encoding/base64"
	"fmt"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
)

// ctxKey is the type for context keys.
type ctxKey int

// Context key values.
const (
	ctxKeyRequestID ctxKey = iota
)

// idLen is a length of request ID.
const idLen = 16

// ID is the ID of a request.  It is an opaque, randomly generated string.
type ID [idLen]byte

// FromString converts string s to id.
func FromString(s string) (id ID, err error) {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)

	_, err = enc.Decode(id[:], []byte(s))
	if err != nil {
		return id, fmt.Errorf("decoding request id: %w", err)
	}

	return id, nil
}

// requestIDRand is used to create [ID]s.
//
// TODO(a.garipov): Consider making a struct instead of using one global source.
var requestIDRand = randutil.NewReader(randutil.MustNewSeed())

// New generates new request ID.
func New() (id ID) {
	_, err := requestIDRand.Read(id[:])
	if err != nil {
		panic(fmt.Errorf("generating random request id: %w", err))
	}

	return id
}

// type check
var _ fmt.Stringer = (*ID)(nil)

// String implements the [fmt.Stringer] interface for ID.
func (i ID) String() (s string) {
	enc := base64.URLEncoding.WithPadding(base64.NoPadding)
	n := enc.EncodedLen(idLen)

	idData64 := make([]byte, n)
	enc.Encode(idData64, i[:])

	return string(idData64)
}

// type check
var _ encoding.TextMarshaler = (*ID)(nil)

// MarshalText implements the [encoding.TextMarshaler] interface for ID.
func (i ID) MarshalText() (data []byte, err error) {
	return []byte(i.String()), nil
}
