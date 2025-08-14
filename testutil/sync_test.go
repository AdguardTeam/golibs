package testutil_test

import (
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRequireSend(t *testing.T) {
	t.Parallel()

	const timeout = 100 * time.Millisecond

	t.Run("success", func(t *testing.T) {
		var numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		ch := make(chan struct{}, 1)
		testutil.RequireSend(tb, ch, struct{}{}, timeout)

		assert.Equal(t, 1, numHelper)
		assert.Len(t, ch, 1)
	})

	t.Run("fail", func(t *testing.T) {
		var numHelper, numErrorf, numFailNow int

		tb := newTestTB()
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onFailNow = func() { numFailNow++ }
		tb.onHelper = func() { numHelper++ }

		ch := make(chan struct{})
		testutil.RequireSend(tb, ch, struct{}{}, timeout)

		assert.Equal(t, 1, numHelper)
		assert.Equal(t, 1, numErrorf)
		assert.Equal(t, 1, numFailNow)
		assert.Len(t, ch, 0)
	})
}

func TestRequireReceive(t *testing.T) {
	const timeout = 100 * time.Millisecond

	t.Run("success", func(t *testing.T) {
		var numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		ch := make(chan struct{}, 1)
		ch <- struct{}{}

		_, ok := testutil.RequireReceive(tb, ch, timeout)
		assert.True(t, ok)
		assert.Equal(t, 1, numHelper)
		assert.Len(t, ch, 0)
	})

	t.Run("fail", func(t *testing.T) {
		var numHelper, numErrorf, numFailNow int

		tb := newTestTB()
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onFailNow = func() { numFailNow++ }
		tb.onHelper = func() { numHelper++ }

		ch := make(chan struct{})
		_, ok := testutil.RequireReceive(tb, ch, timeout)
		assert.False(t, ok)
		assert.Equal(t, 1, numHelper)
		assert.Equal(t, 1, numErrorf)
		assert.Equal(t, 1, numFailNow)
		assert.Len(t, ch, 0)
	})
}
