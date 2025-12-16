// Package requestid contains utilities for working with request ids.
package requestid

import (
	"encoding/base64"
	"fmt"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
)

// idLen is the length of a request ID.
const idLen = 16

// encodedIDLen is the length of encoded request ID.
const encodedIDLen = 22

// ID is the ID of a request.  It is an opaque, randomly generated string.
type ID string

// requestIDRand is used to create [ID]s.
//
// TODO(a.garipov): Consider making a struct instead of using one global source.
var requestIDRand = randutil.NewReader(randutil.MustNewSeed())

// New generates new request ID.
func New() (id ID) {
	reqID := make([]byte, idLen)
	_, err := requestIDRand.Read(reqID)
	if err != nil {
		panic(fmt.Errorf("generating random request id: %w", err))
	}

	enc := base64.URLEncoding.WithPadding(base64.NoPadding)

	encoded := make([]byte, encodedIDLen)
	enc.Encode(encoded, reqID)

	return ID(encoded)
}
