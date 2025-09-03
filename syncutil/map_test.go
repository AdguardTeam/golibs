package syncutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/syncutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Parallel()

	m := syncutil.NewMap[string, int]()

	const (
		key1 = "a"
		key2 = "b"
	)

	require.True(t, t.Run("store_and_delete", func(t *testing.T) {
		m.Store(key1, 1)
		assertMapKey(t, m, key1, 1, assert.True)
		assertMapKey(t, m, key2, 0, assert.False)

		n, loaded := m.LoadOrStore(key1, 2)
		assert.Equal(t, n, 1)
		assert.True(t, loaded)

		n, loaded = m.LoadOrStore(key2, 3)
		assert.Equal(t, n, 3)
		assert.False(t, loaded)

		m.Delete(key1)
		assertMapKey(t, m, key1, 0, assert.False)

		n, loaded = m.LoadAndDelete(key2)
		assert.Equal(t, n, 3)
		assert.True(t, loaded)
		assertMapKey(t, m, key2, 0, assert.False)

		n, loaded = m.LoadAndDelete(key2)
		assert.Equal(t, n, 0)
		assert.False(t, loaded)
	}))

	require.True(t, t.Run("swap", func(t *testing.T) {
		n, swapped := m.Swap(key1, 1)
		assert.Equal(t, n, 0)
		assert.False(t, swapped)

		n, swapped = m.Swap(key1, 2)
		assert.Equal(t, n, 1)
		assert.True(t, swapped)

		swapped = m.CompareAndSwap(key1, 1, 3)
		assert.False(t, swapped)

		swapped = m.CompareAndSwap(key1, 2, 3)
		assert.True(t, swapped)
		assertMapKey(t, m, key1, 3, assert.True)

		deleted := m.CompareAndDelete(key1, 1)
		assert.False(t, deleted)

		deleted = m.CompareAndDelete(key1, 3)
		assert.True(t, deleted)
	}))

	require.True(t, t.Run("range_and_clear", func(t *testing.T) {
		m.Store(key1, 1)

		var n int
		for range m.Range {
			n++
		}

		assert.Equal(t, 1, n)

		m.Clear()

		n = 0
		for range m.Range {
			n++
		}

		assert.Equal(t, 0, n)
	}))
}

// assertMapKey is a helper for checking the presence or absence of a key in a
// map.
func assertMapKey[K comparable, V any](
	tb testing.TB,
	m *syncutil.Map[K, V],
	key K,
	want V,
	wantOK assert.BoolAssertionFunc,
) {
	tb.Helper()

	v, ok := m.Load(key)
	assert.Equal(tb, v, want)
	wantOK(tb, ok)
}

func BenchmarkMap_int(b *testing.B) {
	const (
		key = 1
		val = 2
	)

	m := syncutil.NewMap[int, int]()

	require.True(b, b.Run("store", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			m.Store(key, val)
		}
	}))

	require.True(b, b.Run("load", func(b *testing.B) {
		var got int

		b.ReportAllocs()
		for b.Loop() {
			got, _ = m.Load(key)
		}

		require.Equal(b, val, got)
	}))

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/syncutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkMap_int/store-16   	 9240952	        162.7 ns/op	      48 B/op	       1 allocs/op
	//	BenchmarkMap_int/load-16    	63999354	        18.00 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkMap_struct(b *testing.B) {
	type K struct {
		str string
		num int
	}

	type V struct {
		str string
		num int
	}

	key := K{
		str: "key",
		num: 123,
	}

	val := &V{
		str: "value",
		num: 456,
	}

	m := syncutil.NewMap[K, *V]()

	require.True(b, b.Run("store", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			m.Store(key, val)
		}
	}))

	require.True(b, b.Run("load", func(b *testing.B) {
		var got *V

		b.ReportAllocs()
		for b.Loop() {
			got, _ = m.Load(key)
		}

		require.Equal(b, val, got)
	}))

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/syncutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkMap_struct/store-16   	 4702305	        247.0 ns/op	      72 B/op	       2 allocs/op
	//	BenchmarkMap_struct/load-16    	31171255	        34.20 ns/op	       0 B/op	       0 allocs/op
}
