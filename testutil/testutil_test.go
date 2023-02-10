package testutil_test

import (
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testutil.DiscardLogOutput(m)
}

// Common test constants.
const (
	testName   = "TestName"
	testErrMsg = "test error"
)

// testTB is a testing.TB for tests.
type testTB struct {
	// TB is embedded here simply to make *testTB a testing.TB without actually
	// implementing all methods.
	testing.TB

	onCleanup func(f func())
	onErrorf  func(format string, args ...any)
	onFailNow func()
	onHelper  func()
	onName    func() (name string)
}

// Cleanup implements the [testing.TB] interface for *testTB.
func (t *testTB) Cleanup(f func()) {
	t.onCleanup(f)
}

// Errorf implements the [testing.TB] interface for *testTB.
func (t *testTB) Errorf(format string, args ...any) {
	t.onErrorf(format, args...)
}

// FailNow implements the [testing.TB] interface for *testTB.
func (t *testTB) FailNow() {
	t.onFailNow()
}

// Helper implements the [testing.TB] interface for *testTB.
func (t *testTB) Helper() {
	t.onHelper()
}

// Name implements the [testing.TB] interface for *testTB.
func (t *testTB) Name() (name string) {
	return t.onName()
}

func TestAssertErrorMsg(t *testing.T) {
	t.Parallel()

	t.Run("msg", func(t *testing.T) {
		numHelper := 0
		gotFormat := ""
		var gotArgs []any
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf: func(format string, args ...any) {
				gotFormat = format
				gotArgs = args
			},
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { return testName },
		}

		testutil.AssertErrorMsg(tt, testErrMsg, errors.Error(testErrMsg))

		assert.Greater(t, numHelper, 0)
		assert.Empty(t, gotFormat)
		assert.Empty(t, gotArgs)
	})

	t.Run("msg_mismatch", func(t *testing.T) {
		numHelper := 0
		gotFormat := ""
		var gotArgs []any
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf: func(format string, args ...any) {
				gotFormat = format
				gotArgs = args
			},
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { return testName },
		}

		testutil.AssertErrorMsg(tt, testErrMsg, errors.Error("wrong test error"))

		assert.Greater(t, numHelper, 0)
		assert.NotEmpty(t, gotFormat)
		require.Len(t, gotArgs, 1)

		arg := gotArgs[0]
		assert.Contains(t, arg, testName)
		assert.Contains(t, arg, testErrMsg)
	})

	t.Run("empty_msg", func(t *testing.T) {
		numHelper := 0
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { panic("not implemented") },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		testutil.AssertErrorMsg(tt, "", nil)

		assert.Greater(t, numHelper, 0)
	})

	t.Run("unexpected_success", func(t *testing.T) {
		numHelper := 0
		gotFormat := ""
		var gotArgs []any
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf: func(format string, args ...any) {
				gotFormat = format
				gotArgs = args
			},
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { return testName },
		}

		testutil.AssertErrorMsg(tt, testErrMsg, nil)

		assert.Greater(t, numHelper, 0)
		assert.NotEmpty(t, gotFormat)
		require.Len(t, gotArgs, 1)

		assert.Contains(t, gotArgs[0], testName)
	})
}

// goodCodec is a good [encoding.TextMarshaler] and [encoding.TextUnmarshaler]
// implementation.
type goodCodec struct {
	value []byte
}

// MarshalText implements [encoding.TextMarshaler] for goodCodec.
func (c goodCodec) MarshalText() (b []byte, err error) {
	return c.value, nil
}

// UnmarshalText implements encoding.TextUnmarshaler for goodCodec.
func (c *goodCodec) UnmarshalText(b []byte) (err error) {
	c.value = b

	return nil
}

// badCodec is a bad [encoding.TextMarshaler] and [encoding.TextUnmarshaler]
// implementation.
type badCodec struct {
	value []byte
}

// MarshalText implements [encoding.TextMarshaler] for badCodec.  It implements
// it badly, because it uses a pointer receiver.
func (c *badCodec) MarshalText() (b []byte, err error) {
	return c.value, nil
}

// UnmarshalText implements encoding.TextUnmarshaler for badCodec.  It
// implements it badly, because it uses a non-pointer receiver.
func (c badCodec) UnmarshalText(b []byte) (err error) {
	c.value = b
	_ = c.value

	return nil
}

func TestAssertMarshalText(t *testing.T) {
	t.Parallel()

	t.Run("good", func(t *testing.T) {
		numHelper := 0
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { panic("not implemented") },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		require.NotPanics(t, func() {
			testutil.AssertMarshalText(tt, "good", &goodCodec{value: []byte("good")})
		})
		assert.Greater(t, numHelper, 0)
	})

	t.Run("bad", func(t *testing.T) {
		numHelper := 0
		numErrorf := 0
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { numErrorf++ },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { return testName },
		}

		require.NotPanics(t, func() {
			testutil.AssertMarshalText(tt, "bad", &badCodec{value: []byte("bad")})
		})
		assert.Greater(t, numErrorf, 0)
		assert.Greater(t, numHelper, 0)
	})
}

func TestAssertUnmarshalText(t *testing.T) {
	t.Parallel()

	t.Run("good", func(t *testing.T) {
		numHelper := 0
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { panic("not implemented") },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { panic("not implemented") },
		}

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tt, "good", &goodCodec{value: []byte("good")})
		})
		assert.Greater(t, numHelper, 0)
	})

	t.Run("bad", func(t *testing.T) {
		numHelper := 0
		numErrorf := 0
		tt := &testTB{
			onCleanup: func(_ func()) { panic("not implemented") },
			onErrorf:  func(_ string, _ ...any) { numErrorf++ },
			onFailNow: func() { panic("not implemented") },
			onHelper:  func() { numHelper++ },
			onName:    func() (name string) { return testName },
		}

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tt, "bad", &badCodec{value: []byte("bad")})
		})
		assert.Greater(t, numErrorf, 0)
		assert.Greater(t, numHelper, 0)
	})
}

func TestCleanupAndRequireSuccess(t *testing.T) {
	t.Parallel()

	cleanupFuncCalled := false
	cleanupFunc := func() (err error) {
		cleanupFuncCalled = true

		return nil
	}

	var gotFunc func()
	numHelper := 0
	tt := &testTB{
		onCleanup: func(f func()) {
			gotFunc = f
		},
		onErrorf: func(_ string, _ ...any) { panic("not implemented") },
		onHelper: func() { numHelper++ },
		onName:   func() (name string) { return testName },
	}

	testutil.CleanupAndRequireSuccess(tt, cleanupFunc)

	assert.Greater(t, numHelper, 0)

	gotFunc()
	assert.True(t, cleanupFuncCalled)
}

func TestRequireTypeAssert(t *testing.T) {
	numHelper := 0
	numErrorf := 0
	numFailNow := 0
	tt := &testTB{
		onCleanup: func(_ func()) { panic("not implemented") },
		onErrorf:  func(_ string, _ ...any) { numErrorf++ },
		onFailNow: func() { numFailNow++ },
		onHelper:  func() { numHelper++ },
		onName:    func() (name string) { return testName },
	}

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() { numHelper = 0 })

		assert.Equal(t, "foo", testutil.RequireTypeAssert[string](tt, "foo"))
		assert.Equal(t, 42, testutil.RequireTypeAssert[int](tt, 42))
		assert.Equal(t, true, testutil.RequireTypeAssert[bool](tt, true))
		assert.Equal(t, struct{}{}, testutil.RequireTypeAssert[struct{}](tt, struct{}{}))
		assert.Equal(t, []int{1, 2, 3}, testutil.RequireTypeAssert[[]int](tt, []int{1, 2, 3}))

		assert.Greater(t, numHelper, 5)
		assert.Zero(t, numErrorf)
		assert.Zero(t, numFailNow)
	})

	t.Run("func", func(t *testing.T) {
		t.Cleanup(func() { numHelper = 0 })

		called := false
		f := func() {
			called = true
		}

		ff := testutil.RequireTypeAssert[func()](tt, f)
		ff()

		assert.True(t, called)
	})

	// TODO(e.burkov):  !! remove this??
	t.Run("interfaces", func(t *testing.T) {
		t.Cleanup(func() {
			numHelper = 0
			numErrorf = 0
			numFailNow = 0
		})

		require.Panics(t, func() {
			testutil.RequireTypeAssert[interface{}](tt, nil)
		})

		require.NotPanics(t, func() {
			testutil.RequireTypeAssert[interface{}](tt, t)
			testutil.RequireTypeAssert[testing.TB](tt, t)
		})

		assert.Greater(t, numHelper, 2)
		assert.Greater(t, numErrorf, 1)
		assert.Greater(t, numFailNow, 1)
	})
}
