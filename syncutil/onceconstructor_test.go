package syncutil_test

import (
	"sync/atomic"
	"testing"

	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

func TestOnceConstructor(t *testing.T) {
	t.Parallel()

	numCalls := atomic.Uint32{}
	c := syncutil.NewOnceConstructor(func(k int) (v int) {
		numCalls.Add(1)

		return k + 1
	})

	const (
		n = 10_000

		key  = 1
		want = key + 1
	)

	results := make(chan int, n)

	for range n {
		go func() {
			results <- c.Get(key)
		}()
	}

	for range n {
		got, _ := testutil.RequireReceive(t, results, testTimeout)
		assert.Equal(t, want, got)
	}

	assert.Equal(t, uint32(1), numCalls.Load())
}
