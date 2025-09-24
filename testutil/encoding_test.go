package testutil_test

import (
	"encoding"
	"testing"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// goodCodec is a good [encoding.TextMarshaler] and [encoding.TextUnmarshaler]
// implementation.
type goodCodec struct {
	value []byte
}

// type check
var _ encoding.TextMarshaler = goodCodec{}

// MarshalText implements [encoding.TextMarshaler] for goodCodec.
func (c goodCodec) MarshalText() (b []byte, err error) {
	return c.value, nil
}

// type check
var _ encoding.TextUnmarshaler = (*goodCodec)(nil)

// UnmarshalText implements [encoding.TextUnmarshaler] for goodCodec.
func (c *goodCodec) UnmarshalText(b []byte) (err error) {
	c.value = b

	return nil
}

// badCodec is a bad [encoding.TextMarshaler] and [encoding.TextUnmarshaler]
// implementation.
type badCodec struct {
	value []byte
}

// type check
var _ encoding.TextMarshaler = (*badCodec)(nil)

// MarshalText implements encoding.TextMarshaler for badCodec.  It implements it
// badly, because it uses a pointer receiver.
func (c *badCodec) MarshalText() (b []byte, err error) {
	return c.value, nil
}

// type check
var _ encoding.TextUnmarshaler = badCodec{}

// UnmarshalText implements encoding.TextUnmarshaler for badCodec.  It
// implements it badly, because it uses a non-pointer receiver.
func (c badCodec) UnmarshalText(b []byte) (err error) {
	c.value = b
	_ = c.value

	return nil
}

func TestAssertMarshalText(t *testing.T) {
	t.Parallel()

	numHelper := 0

	tb := newTestTB()
	tb.onHelper = func() { numHelper++ }

	require.NotPanics(t, func() {
		testutil.AssertMarshalText(tb, "good", &goodCodec{value: []byte("good")})
	})

	assert.Greater(t, numHelper, 0)

	// TODO(a.garipov):  Consider ways to test the bad case, as the type system
	// currently prevents it.
}

func TestAssertUnmarshalText(t *testing.T) {
	t.Parallel()

	require.True(t, t.Run("good", func(t *testing.T) {
		t.Parallel()

		numHelper := 0

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tb, "good", &goodCodec{value: []byte("good")})
		})
		assert.Greater(t, numHelper, 0)
	}))

	require.True(t, t.Run("bad", func(t *testing.T) {
		t.Parallel()

		numHelper := 0
		numErrorf := 0

		tb := newTestTB()
		tb.onErrorf = func(s string, _ ...any) { numErrorf++ }
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tb, "bad", &badCodec{value: []byte("bad")})
		})
		assert.Greater(t, numErrorf, 0)
		assert.Greater(t, numHelper, 0)
	}))
}
