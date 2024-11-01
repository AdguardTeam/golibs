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
			b.ResetTimer()
			for range b.N {
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
	//	BenchmarkMapSet_Add/10_strings-16         	  702885	      1514 ns/op	       151.3 ns/add	     491 B/op	       4 allocs/op
	//	BenchmarkMapSet_Add/100_strings
	//	BenchmarkMapSet_Add/100_strings-16        	   66829	     15960 ns/op	       159.6 ns/add	    5615 B/op	      11 allocs/op
	//	BenchmarkMapSet_Add/1000_strings
	//	BenchmarkMapSet_Add/1000_strings-16       	    7400	    206750 ns/op	       206.7 ns/add	   85713 B/op	      36 allocs/op
	//	BenchmarkMapSet_Add/10000_strings
	//	BenchmarkMapSet_Add/10000_strings-16      	     592	   1982110 ns/op	       198.2 ns/add	  676473 B/op	     216 allocs/op
	//	BenchmarkMapSet_Add/100000_strings
	//	BenchmarkMapSet_Add/100000_strings-16     	      49	  23063097 ns/op	       230.6 ns/add	 5597995 B/op	    3903 allocs/op
}

func BenchmarkMapSet_Has(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			set := container.NewMapSet(values...)
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
	//	BenchmarkMapSet_Has
	//	BenchmarkMapSet_Has/10_strings
	//	BenchmarkMapSet_Has/10_strings-16         	171413164	         6.886 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/100_strings
	//	BenchmarkMapSet_Has/100_strings-16        	166819746	         6.607 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/1000_strings
	//	BenchmarkMapSet_Has/1000_strings-16       	179336127	         6.870 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/10000_strings
	//	BenchmarkMapSet_Has/10000_strings-16      	164002748	         6.831 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Has/100000_strings
	//	BenchmarkMapSet_Has/100000_strings-16     	170170257	         6.518 ns/op	       0 B/op	       0 allocs/op
}
