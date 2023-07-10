package cache

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Parallel()

	conf := Config{}
	var rmKey, rmVal atomic.Value
	conf.OnDelete = func(key []byte, val []byte) {
		rmKey.Store(key)
		rmVal.Store(val)
	}

	conf.MaxSize = 12
	conf.MaxElementSize = 12
	conf.MaxCount = 3
	conf.EnableLRU = true
	c := New(conf)

	var d []byte

	// get - not found
	assert.True(t, c.Get([]byte("k1")) == nil)

	// add new
	assert.True(t, !c.Set([]byte("k1"), []byte("v1")))
	assert.True(t, !c.Set([]byte("k2"), []byte("v2")))
	assert.True(t, c.Stats().Count == 2)

	// get added
	d = c.Get([]byte("k1"))
	assert.True(t, bytes.Equal(d, []byte("v1")))
	d = c.Get([]byte("k2"))
	assert.True(t, bytes.Equal(d, []byte("v2")))

	// replace existing
	assert.True(t, c.Set([]byte("k1"), []byte("v!")))
	d = c.Get([]byte("k1"))
	assert.True(t, bytes.Equal(d, []byte("v!")))

	// delete
	c.Del([]byte("k1"))
	assert.True(t, c.Get([]byte("k1")) == nil)
	c.Clear()

	// MaxCount limit
	assert.True(t, !c.Set([]byte("k1"), []byte("v1")))
	assert.True(t, !c.Set([]byte("k2"), []byte("v2")))
	assert.True(t, !c.Set([]byte("k3"), []byte("v3")))

	rmKey.Store([]byte(nil))
	rmVal.Store([]byte(nil))
	assert.True(t, !c.Set([]byte("k4"), []byte("v4"))) // "k1" is removed

	assert.Eventually(t, func() (ok bool) {
		return bytes.Equal(rmKey.Load().([]byte), []byte("k1")) &&
			bytes.Equal(rmVal.Load().([]byte), []byte("v1"))
	}, 1*time.Second, 100*time.Millisecond)

	c.Clear()

	// MaxSize limit
	assert.True(t, !c.Set([]byte("k1"), []byte("v1")))

	rmKey.Store([]byte(nil))
	rmVal.Store([]byte(nil))
	assert.True(t, !c.Set([]byte("k2"), []byte("1234567"))) // "k1" is removed
	assert.Eventually(t, func() (ok bool) {
		return bytes.Equal(rmKey.Load().([]byte), []byte("k1"))
	}, 1*time.Second, 100*time.Millisecond)

	c.Clear()

	// MaxElementSize limit
	assert.True(t, !c.Set([]byte("k1"), []byte("12345678901")))
	assert.True(t, c.Get([]byte("k1")) == nil)

	c.Del([]byte("k1"))
	assert.True(t, c.Stats().Count == 0)
	assert.True(t, c.Stats().Size == 0)
}

// Set, get, delete items in parallel
func TestParallel(t *testing.T) {
	t.Parallel()

	conf := Config{}
	conf.EnableLRU = true
	conf.MaxSize = 1024
	c := New(conf)

	wg := sync.WaitGroup{}
	N := 100
	for w := 0; w != 100; w++ {
		wg.Add(1)
		go func(wid int) {
			for i := 0; i != N; i++ {
				key := []byte(fmt.Sprintf("key-%d-%d", wid, i))
				val := []byte{1, 2, 3, byte(i % 255)}
				_ = c.Set(key, val)

				rval := c.Get(key)
				if rval != nil {
					assert.True(t, val[3] == rval[3])
				}

				c.Del(key)
			}
			wg.Done()
		}(w)
	}

	wg.Wait()
}
