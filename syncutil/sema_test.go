package syncutil_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func TestChanSemaphore(t *testing.T) {
	const (
		maxRes       = 3
		numGoroutine = 10_000
	)

	s := syncutil.NewChanSemaphore(maxRes)

	ctx := context.Background()
	current := &atomic.Int64{}

	for range numGoroutine {
		err := s.Acquire(ctx)
		require.NoError(t, err)

		go func() {
			defer s.Release()
			defer current.Add(-1)

			pt := &testutil.PanicT{}
			newCurrent := current.Add(1)
			require.LessOrEqual(pt, newCurrent, int64(maxRes))
		}()
	}
}

func TestEmptySemaphore(t *testing.T) {
	const (
		numGoroutine = 10
	)

	s := syncutil.EmptySemaphore{}

	ctx := context.Background()
	current := &atomic.Int64{}

	for range numGoroutine {
		err := s.Acquire(ctx)
		require.NoError(t, err)

		go func() {
			defer s.Release()
			defer current.Add(-1)

			pt := &testutil.PanicT{}
			newCurrent := current.Add(1)
			require.LessOrEqual(pt, newCurrent, int64(numGoroutine))
		}()
	}
}
