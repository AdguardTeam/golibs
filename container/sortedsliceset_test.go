package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/container"
	"github.com/stretchr/testify/require"
)

func BenchmarkSortedSliceSet_Add(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			var set *container.SortedSliceSet[string]

			values := newRandStrs(n, randStrLen)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				set = container.NewSortedSliceSet[string]()
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
	//	BenchmarkSortedSliceSet_Add
	//	BenchmarkSortedSliceSet_Add/10_strings
	//	BenchmarkSortedSliceSet_Add/10_strings-16         	 1000000	      1717 ns/op	       171.7 ns/add	     520 B/op	       6 allocs/op
	//	BenchmarkSortedSliceSet_Add/100_strings
	//	BenchmarkSortedSliceSet_Add/100_strings-16        	   66405	     16969 ns/op	       169.7 ns/add	    4488 B/op	       9 allocs/op
	//	BenchmarkSortedSliceSet_Add/1000_strings
	//	BenchmarkSortedSliceSet_Add/1000_strings-16       	    2820	    386448 ns/op	       386.4 ns/add	   35208 B/op	      12 allocs/op
	//	BenchmarkSortedSliceSet_Add/10000_strings
	//	BenchmarkSortedSliceSet_Add/10000_strings-16      	      97	  11721214 ns/op	      1172 ns/add	  665992 B/op	      19 allocs/op
	//	BenchmarkSortedSliceSet_Add/100000_strings
	//	BenchmarkSortedSliceSet_Add/100000_strings-16     	       1	1330335923 ns/op	     13303 ns/add	 8923656 B/op	      31 allocs/op
}

func BenchmarkSortedSliceSet_Has(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			set := container.NewSortedSliceSet(values...)
			value := values[n/2]

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				sinkBool = set.Has(value)
			}

			require.True(b, sinkBool)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/container
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkSortedSliceSet_Has
	//	BenchmarkSortedSliceSet_Has/10_strings
	//	BenchmarkSortedSliceSet_Has/10_strings-16         	75554308	        14.48 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/100_strings
	//	BenchmarkSortedSliceSet_Has/100_strings-16        	43907683	        28.66 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/1000_strings
	//	BenchmarkSortedSliceSet_Has/1000_strings-16       	29775063	        40.41 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/10000_strings
	//	BenchmarkSortedSliceSet_Has/10000_strings-16      	20591700	        57.83 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/100000_strings
	//	BenchmarkSortedSliceSet_Has/100000_strings-16     	13058570	        78.08 ns/op	       0 B/op	       0 allocs/op
}
