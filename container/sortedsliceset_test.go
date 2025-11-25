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
			set := container.NewSortedSliceSet(make([]string, 0, n)...)

			b.ReportAllocs()
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
	//	BenchmarkSortedSliceSet_Union/10_strings-8         	29344582	        40.86 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100_strings-8        	 3049860	       392.6 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/1000_strings-8       	  256148	      4663 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10000_strings-8      	   17222	     69520 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100000_strings-8     	    1365	    862403 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10_strings_receiver_is_x-8         	18862746	        65.19 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100_strings_receiver_is_x-8        	 2225612	       538.1 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/1000_strings_receiver_is_x-8       	  246709	      4876 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/10000_strings_receiver_is_x-8      	   24936	     48104 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Union/100000_strings_receiver_is_x-8     	    2233	    483712 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkSortedSliceSet_Intersection(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)
			set := container.NewSortedSliceSet(make([]string, 0, n)...)

			b.ReportAllocs()
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
	//	BenchmarkSortedSliceSet_Intersection/10_strings-8         	37729012	        31.54 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100_strings-8        	 3699164	       325.7 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/1000_strings-8       	  338103	      3585 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10000_strings-8      	   23232	     51780 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100000_strings-8     	    1566	    776892 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10_strings_receiver_is_x-8         	645125013	         1.856 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100_strings_receiver_is_x-8        	644168484	         1.859 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/1000_strings_receiver_is_x-8       	645357904	         1.857 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/10000_strings_receiver_is_x-8      	644380356	         1.859 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersection/100000_strings_receiver_is_x-8     	645454954	         1.858 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkSortedSliceSet_Intersects(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewSortedSliceSet(values[:n/2]...)
			y := container.NewSortedSliceSet(values[n/2:]...)

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
	//	BenchmarkSortedSliceSet_Intersects/10_strings-8         	40107934	        29.57 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersects/100_strings-8        	 3717073	       327.6 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersects/1000_strings-8       	  342010	      3553 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersects/10000_strings-8      	   21056	     55329 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkSortedSliceSet_Intersects/100000_strings-8     	    1566	    776398 ns/op	       0 B/op	       0 allocs/op
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
