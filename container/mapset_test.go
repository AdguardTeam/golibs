package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/container"
	"github.com/stretchr/testify/require"
)

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
