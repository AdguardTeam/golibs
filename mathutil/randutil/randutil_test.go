package randutil_test

import "math/rand/v2"

// testGoroutinesNum is the number of goroutines for tests.
const testGoroutinesNum = 512

// testSeed is the common seed for tests.
var testSeed = [32]byte{}

// testRNG is the common random-number generator for tests.
var testRNG = rand.New(rand.NewChaCha8(testSeed))
