package bytecache_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/cache/bytecache"
	"github.com/stretchr/testify/assert"
)

// assertGet is a helper that wraps [bytecache.Interface.Get].  wantOK must not
// be nil.
func assertGet(
	tb testing.TB,
	c bytecache.Interface,
	key []byte,
	want []byte,
	wantOK assert.BoolAssertionFunc,
) {
	tb.Helper()

	got, ok := c.Get(key)
	assert.Equal(tb, want, got)
	wantOK(tb, ok)
}

func TestEmpty(t *testing.T) {
	c := bytecache.Empty{}

	c.Clear()

	c.Set([]byte("a"), []byte("1"))
	assert.Zero(t, c.Stats())
	assertGet(t, c, []byte("a"), nil, assert.False)

	c.Set([]byte("a"), []byte("2"))
	c.Del([]byte("a"))
	assert.Zero(t, c.Stats())
}
