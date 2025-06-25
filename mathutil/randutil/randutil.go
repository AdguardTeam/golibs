// Package randutil contains utilities for random numbers.
package randutil

import (
	cryptorand "crypto/rand"
)

// MustNewSeed returns new 32-byte seed for pseudorandom generators.
func MustNewSeed() (seed [32]byte) {
	// NOTE:  crypto/rand.Read crashes the program if there are any errors.
	_, _ = cryptorand.Read(seed[:])

	return seed
}
