package container_test

import (
	"math/rand"
	"time"
)

// Common constants for tests.
const (
	randStrLen = 8
	setMaxLen  = 100_000
)

// newRandStrs returns a slice of random strings of length l with each string
// being strLen bytes long.
func newRandStrs(l, strLen int) (strs []string) {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	strs = make([]string, 0, l)
	for range l {
		data := make([]byte, strLen)
		_, _ = rng.Read(data)

		strs = append(strs, string(data))
	}

	return strs
}
