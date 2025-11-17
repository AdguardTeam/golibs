package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSortedSliceSet(t *testing.T) {
	values := []int{1, 1, 1}

	set := container.NewSortedSliceSet(values...)
	assert.Equal(t, []int{1}, set.Values())

	set.Add(2)
	assert.Equal(t, []int{1, 2}, set.Values())

	set.Delete(2)
	assert.Equal(t, []int{1}, set.Values())

	set.Clear()
	assert.Equal(t, []int{}, set.Values())
}

func BenchmarkSortedSliceSet_Union(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)

			b.ReportAllocs()
			set := container.NewSortedSliceSet(make([]string, 0, n)...)
			for b.Loop() {
				set.Union(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)

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
	//	BenchmarkSortedSliceSet_Union/10_strings-8         	26431711	        45.33 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100_strings-8        	 2654437	       460.3 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/1000_strings-8       	  203306	      5891 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10000_strings-8      	   14985	     80245 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100000_strings-8     	    1292	    936201 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10_strings_receiver_is_x-8         	18278089	        65.88 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100_strings_receiver_is_x-8        	 2210296	       543.9 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/1000_strings_receiver_is_x-8       	  244483	      4887 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10000_strings_receiver_is_x-8      	   24802	     48247 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100000_strings_receiver_is_x-8     	    2238	    483468 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkSortedSliceSet_Intersection(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)

			b.ReportAllocs()
			set := container.NewSortedSliceSet(make([]string, 0, n)...)
			for b.Loop() {
				set.Intersection(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)

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
	//	BenchmarkSortedSliceSet_Intersection/10_strings-8         	34232856	        34.76 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100_strings-8        	 3094461	       389.1 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/1000_strings-8       	  256555	      4692 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10000_strings-8      	   18898	     64293 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100000_strings-8     	    1426	    848021 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10_strings_receiver_is_x-8         	500120686	         2.396 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100_strings_receiver_is_x-8        	500622649	         2.397 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/1000_strings_receiver_is_x-8       	501104431	         2.396 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10000_strings_receiver_is_x-8      	500923489	         2.395 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100000_strings_receiver_is_x-8     	497481928	         2.421 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkSortedSliceSet_Add(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			var set *container.SortedSliceSet[string]

			values := newRandStrs(n, randStrLen)

			b.ReportAllocs()
			for b.Loop() {
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
	//	BenchmarkSortedSliceSet_Add/10_strings-16         	  946683	      1303 ns/op	       130.3 ns/add	     520 B/op	       6 allocs/op
	//	BenchmarkSortedSliceSet_Add/100_strings
	//	BenchmarkSortedSliceSet_Add/100_strings-16        	   82279	     14888 ns/op	       148.9 ns/add	    4488 B/op	       9 allocs/op
	//	BenchmarkSortedSliceSet_Add/1000_strings
	//	BenchmarkSortedSliceSet_Add/1000_strings-16       	    7930	    316216 ns/op	       316.2 ns/add	   35208 B/op	      12 allocs/op
	//	BenchmarkSortedSliceSet_Add/10000_strings
	//	BenchmarkSortedSliceSet_Add/10000_strings-16      	      88	  14064210 ns/op	      1406 ns/add	  666001 B/op	      19 allocs/op
	//	BenchmarkSortedSliceSet_Add/100000_strings
	//	BenchmarkSortedSliceSet_Add/100000_strings-16     	       1	1426843470 ns/op	     14268 ns/add	 8923624 B/op	      33 allocs/op
}

func BenchmarkSortedSliceSet_Has(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			set := container.NewSortedSliceSet(values...)
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
	//	BenchmarkSortedSliceSet_Has
	//	BenchmarkSortedSliceSet_Has/10_strings
	//	BenchmarkSortedSliceSet_Has/10_strings-16         	70456480	        15.81 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/100_strings
	//	BenchmarkSortedSliceSet_Has/100_strings-16        	38889938	        29.56 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/1000_strings
	//	BenchmarkSortedSliceSet_Has/1000_strings-16       	26026233	        43.81 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/10000_strings
	//	BenchmarkSortedSliceSet_Has/10000_strings-16      	18938230	        61.65 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Has/100000_strings
	//	BenchmarkSortedSliceSet_Has/100000_strings-16     	14932011	        79.37 ns/op	       0 B/op	       0 allocs/op
}
