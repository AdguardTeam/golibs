package randutil_test

import (
	"math/rand/v2"
	"sync"
	"testing"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
)

func TestLockedSource_race(t *testing.T) {
	t.Parallel()

	src := randutil.NewLockedSource(rand.NewPCG(0, 0))

	wg := &sync.WaitGroup{}
	wg.Add(testGoroutinesNum)

	startCh := make(chan struct{})
	for range testGoroutinesNum {
		go func() {
			defer wg.Done()

			<-startCh
			for range 1_000 {
				_ = src.Uint64()
			}
		}()
	}

	close(startCh)

	wg.Wait()
}

func BenchmarkLockedSource_Uint64(b *testing.B) {
	src := randutil.NewLockedSource(rand.NewChaCha8(testSeed))

	b.ReportAllocs()
	for b.Loop() {
		_ = src.Uint64()
	}

	// Most recent results:
	//
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/mathutil/randutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkLockedSource_Uint64-16    	65685742	        17.07 ns/op	       0 B/op	       0 allocs/op
}
