//go:build go1.21 && unix

package service_test

import "time"

// testTimeout is the common timeout for tests
const testTimeout = 1 * time.Second
