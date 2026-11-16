package bytecache_test

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/AdguardTeam/golibs/cache/bytecache"
	"github.com/c2h5oh/datasize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Common cache keys and values for tests.
var (
	testKey1 = []byte("k1")
	testKey2 = []byte("k2")

	testVal1 = []byte("v1")
	testVal2 = []byte("v2")
)

func TestDefault(t *testing.T) {
	t.Parallel()

	var evictedKey, evictedVal []byte
	c := bytecache.NewDefault(&bytecache.DefaultConfig{
		MaxSize:        12,
		MaxElementSize: 12,
		MaxCount:       3,
		OnDelete: func(key, val []byte) {
			evictedKey = slices.Clone(key)
			evictedVal = slices.Clone(val)
		},
	})

	require.True(t, t.Run("get_not_found", func(t *testing.T) {
		assertGet(t, c, testKey1, nil, assert.False)
	}))

	require.True(t, t.Run("get", func(t *testing.T) {
		c.Set(testKey1, testVal1)
		c.Set(testKey2, testVal2)

		assertGet(t, c, testKey1, testVal1, assert.True)
		assertGet(t, c, testKey2, testVal2, assert.True)

		assert.Equal(t, uint64(2), c.Stats().Count)
	}))

	require.True(t, t.Run("replace", func(t *testing.T) {
		c.Set(testKey1, []byte("v1a"))

		assertGet(t, c, testKey1, []byte("v1a"), assert.True)
	}))

	require.True(t, t.Run("delete", func(t *testing.T) {
		c.Del(testKey1)

		assertGet(t, c, testKey1, nil, assert.False)
	}))

	require.True(t, t.Run("clear", func(t *testing.T) {
		c.Clear()

		assert.Equal(t, bytecache.Stats{}, c.Stats())
	}))

	require.True(t, t.Run("max_count", func(t *testing.T) {
		c.Set(testKey1, testVal1)
		c.Set(testKey2, testVal2)
		c.Set([]byte("k3"), []byte("v3"))

		evictedKey, evictedVal = nil, nil

		c.Set([]byte("k4"), []byte("v4"))
		assert.Equal(t, testKey1, evictedKey)
		assert.Equal(t, testVal1, evictedVal)
	}))

	require.True(t, t.Run("max_size", func(t *testing.T) {
		c.Clear()

		c.Set(testKey1, testVal1)

		evictedKey, evictedVal = nil, nil

		c.Set(testKey2, []byte("1234567"))
		assert.Equal(t, testKey1, evictedKey)
	}))

	require.True(t, t.Run("max_elem_size", func(t *testing.T) {
		c.Clear()

		c.Set(testKey1, []byte("12345678901"))

		assertGet(t, c, testKey1, nil, assert.False)

		c.Del(testKey1)

		assert.Zero(t, c.Stats().Count)
		assert.Zero(t, c.Stats().Size)
	}))
}

func TestDefault_parallel(t *testing.T) {
	t.Parallel()

	const (
		numGoroutine = 100
		numIter      = 100
	)

	c := bytecache.NewDefault(&bytecache.DefaultConfig{
		MaxSize: 1000,
	})

	useCache := func(goID int) {
		for i := range numIter {
			key := fmt.Appendf(nil, "key-%d-%d", goID, i)
			val := []byte{1, 2, 3, byte(i % 255)}
			c.Set(key, val)

			got, ok := c.Get(key)
			if ok {
				assert.Equal(t, val[3], got[3])
			}

			c.Del(key)
		}
	}

	wg := sync.WaitGroup{}
	for g := range numGoroutine {
		wg.Go(func() { useCache(g) })
	}

	wg.Wait()
}

func TestDefault_Set_evict(t *testing.T) {
	t.Parallel()

	var evictedKeys [][]byte
	c := bytecache.NewDefault(&bytecache.DefaultConfig{
		OnDelete: func(key, _ []byte) {
			evictedKeys = append(evictedKeys, slices.Clone(key))
		},
		MaxSize: 8,
	})

	c.Set(testKey1, testVal1)
	c.Set(testKey2, testVal2)

	require.Equal(t, uint64(2), c.Stats().Count)

	c.Set([]byte("1234"), []byte("5678"))

	assert.Equal(t, uint64(1), c.Stats().Count)
	assert.Equal(t, [][]byte{testKey1, testKey2}, evictedKeys)
}

func TestDefault_Set_replaceAtCapacity(t *testing.T) {
	t.Parallel()

	var evictedKeys [][]byte
	c := bytecache.NewDefault(&bytecache.DefaultConfig{
		OnDelete: func(key, _ []byte) {
			evictedKeys = append(evictedKeys, slices.Clone(key))
		},
		MaxCount: 2,
	})

	c.Set(testKey1, testVal1)
	c.Set(testKey2, testVal2)

	// Make testKey2 the most recently used entry.  testKey1 is now the least
	// recently used entry and must survive a replacement of testKey2.
	assertGet(t, c, testKey2, testVal2, assert.True)

	newVal2 := append(slices.Clone(testVal2), 'a')
	c.Set(testKey2, newVal2)
	assertGet(t, c, testKey1, testVal1, assert.True)
	assertGet(t, c, testKey2, newVal2, assert.True)

	assert.Empty(t, evictedKeys)
	assert.Equal(t, uint64(2), c.Stats().Count)
}

func TestDefault_Set_replaceExceedsMaxSize(t *testing.T) {
	t.Parallel()

	const maxSize datasize.ByteSize = 10

	c := bytecache.NewDefault(&bytecache.DefaultConfig{MaxSize: maxSize})

	c.Set(testKey1, []byte("123"))
	c.Set(testKey2, []byte("456"))
	require.Equal(t, maxSize, c.Stats().Size)

	// The new value is larger than the old one and pushes the cache over
	// MaxSize.
	c.Set(testKey2, []byte("123456"))

	assert.Equal(t, datasize.ByteSize(8), c.Stats().Size)
	assert.Equal(t, uint64(1), c.Stats().Count)
}

func TestDefault_Set_replaceLRUKey(t *testing.T) {
	t.Parallel()

	var evictedKeys [][]byte
	c := bytecache.NewDefault(&bytecache.DefaultConfig{
		OnDelete: func(key, _ []byte) {
			evictedKeys = append(evictedKeys, slices.Clone(key))
		},
		MaxCount: 2,
	})

	c.Set(testKey1, testVal1)
	c.Set(testKey2, testVal2)

	// testKey1 is the least recently used entry.  Replacing it must not evict
	// it as a side effect.
	newVal1 := append(slices.Clone(testVal1), 'a')
	c.Set(testKey1, newVal1)

	assertGet(t, c, testKey1, newVal1, assert.True)
	assertGet(t, c, testKey2, testVal2, assert.True)

	assert.Empty(t, evictedKeys)
}

func BenchmarkDefault_Get(b *testing.B) {
	c := bytecache.NewDefault(nil)

	c.Set(testKey1, testVal1)

	var got []byte
	var ok bool
	b.ReportAllocs()
	for b.Loop() {
		got, ok = c.Get(testKey1)
	}

	require.Equal(b, testVal1, got)
	require.True(b, ok)

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/cache/bytecache
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkDefault_Get-24    	41814622	        28.05 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkDefault_Set(b *testing.B) {
	c := bytecache.NewDefault(nil)

	b.ReportAllocs()
	for b.Loop() {
		c.Set(testKey1, testVal1)
	}

	got, ok := c.Get(testKey1)
	require.Equal(b, testVal1, got)
	require.True(b, ok)

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/cache/bytecache
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkDefault_Set-24    	13244890	        86.65 ns/op	      66 B/op	       3 allocs/op
}
