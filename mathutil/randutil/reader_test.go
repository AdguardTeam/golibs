package randutil_test

import (
	"sync"
	"testing"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
	"github.com/stretchr/testify/require"
)

func TestReader_race(t *testing.T) {
	t.Parallel()

	const length = 128

	reader := randutil.NewReader(randutil.MustNewSeed())

	wg := &sync.WaitGroup{}
	startCh := make(chan struct{})
	for range testGoroutinesNum {
		wg.Go(func() {
			<-startCh
			for range 1_000 {
				buf := make([]byte, length)
				_, _ = reader.Read(buf)
			}
		})
	}

	close(startCh)

	wg.Wait()
}

func BenchmarkReader_Read(b *testing.B) {
	const length = 16

	reader := randutil.NewReader(testSeed)

	var n int
	var err error

	b.ReportAllocs()
	buf := make([]byte, length)
	for b.Loop() {
		n, err = reader.Read(buf)
	}

	require.Equal(b, length, n)
	require.NoError(b, err)

	// Most recent results:
	//
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/mathutil/randutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkReader_Read-16    	49254069	        24.56 ns/op	       0 B/op	       0 allocs/op
}
