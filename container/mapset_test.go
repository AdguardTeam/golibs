package container_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/container"
	"github.com/stretchr/testify/require"
)

func BenchmarkMapSet_Intersection(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			set := container.NewMapSet(make([]string, n)...)
			set.Clear()

			b.ReportAllocs()
			for b.Loop() {
				set.Intersection(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
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
	//  BenchmarkMapSet_Intersection/10_strings-8         	14420484	     85.18 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100_strings-8        	 1947295	     622.4 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/1000_strings-8       	  199666	      6118 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10000_strings-8      	   15506	     77826 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100000_strings-8     	     913	   1328179 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10_strings_receiver_is_x-8         	233993522	         5.151 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100_strings_receiver_is_x-8        	233974759	         5.138 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/1000_strings_receiver_is_x-8       	233754374	         5.128 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/10000_strings_receiver_is_x-8      	232039011	         5.171 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Intersection/100000_strings_receiver_is_x-8     	231925400	         5.153 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMapSet_Union(b *testing.B) {
	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
			y := container.NewMapSet(values[n/2:]...)

			b.ReportAllocs()
			set := container.NewMapSet(make([]string, 0)...)
			for b.Loop() {
				set.Union(x, y)
			}
		})
	}

	for n := 10; n <= setMaxLen; n *= 10 {
		b.Run(fmt.Sprintf("%d_strings_receiver_is_x", n), func(b *testing.B) {
			values := newRandStrs(n, randStrLen)
			x := container.NewMapSet(values[:n/2]...)
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
	//	BenchmarkMapSet_Union/10_strings-8         	6467498	       186.9 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100_strings-8        	  797383	      1504 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/1000_strings-8       	   74584	     16300 ns/op	       1 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10000_strings-8      	    6807	    175736 ns/op	     128 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100000_strings-8     	     616	   1995085 ns/op	   11616 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10_strings_receiver_is_x-8         	14227329	        83.86 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100_strings_receiver_is_x-8        	 1601006	       746.0 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/1000_strings_receiver_is_x-8       	  155972	      7721 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/10000_strings_receiver_is_x-8      	   13354	     89850 ns/op	      32 B/op	       0 allocs/op
	//	BenchmarkMapSet_Union/100000_strings_receiver_is_x-8     	    1098	   1078387 ns/op	    3183 B/op	       0 allocs/op
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
