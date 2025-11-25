package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/container"
	"github.com/stretchr/testify/require"
)

func BenchmarkMapSet_Union(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)

			x := container.NewMapSet(values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			// Warmup to fill map.
			set := container.NewMapSet(values...)

			b.ReportAllocs()
			for b.Loop() {
				set.Union(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)

			// Warmup to fill map.
			x := container.NewMapSet(values...)
			x.Clear()

			appendToMapSet(x, values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			b.ReportAllocs()
			for b.Loop() {
				x.Union(x, y)
			}
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: Apple M3
	//	BenchmarkMapSet_Union/10_strings-8         	 6075039	       181.3 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100_strings-8        	  815840	      1499 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/1000_strings-8       	   75883	     15797 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10000_strings-8      	    7112	    172980 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100000_strings-8     	     597	   2019998 ns/op	     186 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10_strings_receiver_is_x-8         	14187638	        84.68 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100_strings_receiver_is_x-8        	 1696568	       708.3 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/1000_strings_receiver_is_x-8       	  151192	      7901 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10000_strings_receiver_is_x-8      	   10000	    102098 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100000_strings_receiver_is_x-8     	    1057	   1101593 ns/op	       0 B/op	       0 allocs/op
}

// appendToMapSet is a helper function that adds all given values to s.
func appendToMapSet(s *container.MapSet[string], values ...string) {
	for _, val := range values {
		s.Add(val)
	}
}

func BenchmarkMapSet_Intersection(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			// Warmup to fill map.
			set := container.NewMapSet(values...)

			b.ReportAllocs()
			for b.Loop() {
				set.Intersection(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)

			// Warmup to fill map.
			x := container.NewMapSet(values...)
			x.Clear()

			appendToMapSet(x, values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			b.ReportAllocs()
			for b.Loop() {
				x.Intersection(x, y)
			}
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: Apple M3
	//	BenchmarkMapSet_Intersection/10_strings-8         	14810011	        80.95 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100_strings-8        	 1999047	       601.8 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/1000_strings-8       	  203292	      5880 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10000_strings-8      	   15690	     76495 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100000_strings-8     	     961	   1286111 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10_strings_receiver_is_x-8         	234085954	         5.092 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100_strings_receiver_is_x-8        	235796456	         5.089 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/1000_strings_receiver_is_x-8       	236973238	         5.064 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10000_strings_receiver_is_x-8      	235984836	         5.110 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100000_strings_receiver_is_x-8     	232043947	         5.214 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMapSet_Intersects(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			b.ReportAllocs()
			for b.Loop() {
				x.Intersects(y)
			}
		})
	}

	// Most recent results:
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: Apple M3
	//	BenchmarkMapSet_Intersects/10_strings-8         	14720714	        81.79 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersects/100_strings-8        	 1891896	       611.2 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersects/1000_strings-8       	  206956	      5787 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersects/10000_strings-8      	   16372	     73291 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersects/100000_strings-8     	     974	   1271732 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMapSet_Add(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			var set *container.MapSet[string]

			values := newRandStrs(n, randStrLen)

			b.ReportAllocs()
			for b.Loop() {
				set = container.NewMapSet[string]()
				for _, v := range values {
					set.Add(v)
				}
			}

			perIter := b.Elapsed() / time.Duration(b.N)
			b.ReportMetric(float64(perIter)/float64(n), "ns/add")

			require.True(b, set.Has(values[0]))
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkMapSet_Add
	//	BenchmarkMapSet_Add/10_strings
	//	BenchmarkMapSet_Add/10_strings-16         	  941695	      1259 ns/op	       125.8 ns/add	     720 B/op	       6 allocs/op
	//	BenchmarkMapSet_Add/100_strings
	//	BenchmarkMapSet_Add/100_strings-16        	   95310	     12189 ns/op	       121.9 ns/add	    6960 B/op	      12 allocs/op
	//	BenchmarkMapSet_Add/1000_strings
	//	BenchmarkMapSet_Add/1000_strings-16       	    5937	    183340 ns/op	       183.3 ns/add	  109025 B/op	      23 allocs/op
	//	BenchmarkMapSet_Add/10000_strings
	//	BenchmarkMapSet_Add/10000_strings-16      	     645	   1698322 ns/op	       169.8 ns/add	  873551 B/op	      82 allocs/op
	//	BenchmarkMapSet_Add/100000_strings
	//	BenchmarkMapSet_Add/100000_strings-16     	      56	  21119326 ns/op	       211.2 ns/add	 6990508 B/op	     536 allocs/op
}

func BenchmarkMapSet_Has(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			set := container.NewMapSet(values...)
			value := values[n/2]

			var ok bool
			b.ReportAllocs()
			for b.Loop() {
				ok = set.Has(value)
			}

			require.True(b, ok)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkMapSet_Has
	//	BenchmarkMapSet_Has/10_strings
	//	BenchmarkMapSet_Has/10_strings-16         	110246420	        10.86 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/100_strings
	//	BenchmarkMapSet_Has/100_strings-16        	89053683	        11.71 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/1000_strings
	//	BenchmarkMapSet_Has/1000_strings-16       	89003828	        12.94 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/10000_strings
	//	BenchmarkMapSet_Has/10000_strings-16      	96195098	        12.63 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/100000_strings
	//	BenchmarkMapSet_Has/100000_strings-16     	88061529	        13.01 ns/op	       0 B/op	       0 allocs/op
}
