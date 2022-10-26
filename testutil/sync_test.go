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
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { panic("not implemented") },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		ch := make(chan struct{}, 1)
		testutil.RequireSend(tt, ch, struct{}{}, timeout)

		assert.Equal(t, 1, numHelper)
		assert.Len(t, ch, 1)
	})

	t.Run("fail", func(t *testing.T) {
		var numHelper, numErrorf, numFailNow int
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { numErrorf++ },
			onFailNow: func() { numFailNow++ },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		ch := make(chan struct{})
		testutil.RequireSend(tt, ch, struct{}{}, timeout)

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
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { panic("not implemented") },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		ch := make(chan struct{}, 1)
		ch <- struct{}{}

		_, ok := testutil.RequireReceive(tt, ch, timeout)
		assert.True(t, ok)
		assert.Equal(t, 1, numHelper)
		assert.Len(t, ch, 0)
	})

	t.Run("fail", func(t *testing.T) {
		var numHelper, numErrorf, numFailNow int
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { numErrorf++ },
			onFailNow: func() { numFailNow++ },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		ch := make(chan struct{})
		_, ok := testutil.RequireReceive(tt, ch, timeout)
		assert.False(t, ok)
		assert.Equal(t, 1, numHelper)
		assert.Equal(t, 1, numErrorf)
		assert.Equal(t, 1, numFailNow)
		assert.Len(t, ch, 0)
	})
}
