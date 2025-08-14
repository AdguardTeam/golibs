package testutil_test

import (
	"fmt"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Common test constants.
const (
	testName   = "TestName"
	testErrMsg = "test error"
)

// testTB is a [testing.TB] for tests.
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

// newTestTB returns a new *testTB all settable methods of which panic.
func newTestTB() (t *testTB) {
	return &testTB{
		onCleanup: func(f func()) { panic(testutil.UnexpectedCall(f)) },
		onErrorf:  func(f string, args ...any) { panic(testutil.UnexpectedCall(f, args)) },
		onFailNow: func() { panic(testutil.UnexpectedCall()) },
		onHelper:  func() { panic(testutil.UnexpectedCall()) },
		onName:    func() (name string) { panic(testutil.UnexpectedCall()) },
	}
}

func TestAssertErrorMsg(t *testing.T) {
	t.Parallel()

	t.Run("msg", func(t *testing.T) {
		numHelper := 0
		gotFormat := ""
		var gotArgs []any

		tb := newTestTB()
		tb.onErrorf = func(format string, args ...any) {
			gotFormat = format
			gotArgs = args
		}
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		testutil.AssertErrorMsg(tb, testErrMsg, errors.Error(testErrMsg))

		assert.Greater(t, numHelper, 0)
		assert.Empty(t, gotFormat)
		assert.Empty(t, gotArgs)
	})

	t.Run("msg_mismatch", func(t *testing.T) {
		numHelper := 0
		gotFormat := ""
		var gotArgs []any

		tb := newTestTB()
		tb.onErrorf = func(format string, args ...any) {
			gotFormat = format
			gotArgs = args
		}
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		testutil.AssertErrorMsg(tb, testErrMsg, errors.Error("wrong test error"))

		assert.Greater(t, numHelper, 0)
		assert.NotEmpty(t, gotFormat)
		require.Len(t, gotArgs, 1)

		argStr := fmt.Sprint(gotArgs[0])
		assert.Contains(t, argStr, testName)
		assert.Contains(t, argStr, testErrMsg)
	})

	t.Run("empty_msg", func(t *testing.T) {
		numHelper := 0

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		testutil.AssertErrorMsg(tb, "", nil)

		assert.Greater(t, numHelper, 0)
	})
}

// goodCodec is a good encoding.TextMarshaler and encoding.TextUnmarshaler
// implementation.
type goodCodec struct {
	value []byte
}

// MarshalText implements encoding.TextMarshaler for goodCodec.
func (c goodCodec) MarshalText() (b []byte, err error) {
	return c.value, nil
}

// UnmarshalText implements encoding.TextUnmarshaler for goodCodec.
func (c *goodCodec) UnmarshalText(b []byte) (err error) {
	c.value = b

	return nil
}

// badCodec is a bad encoding.TextMarshaler and encoding.TextUnmarshaler
// implementation.
type badCodec struct {
	value []byte
}

// MarshalText implements encoding.TextMarshaler for badCodec.  It implements it
// badly, because it uses a pointer receiver.
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

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		require.NotPanics(t, func() {
			testutil.AssertMarshalText(tb, "good", &goodCodec{value: []byte("good")})
		})
		assert.Greater(t, numHelper, 0)
	})

	t.Run("bad", func(t *testing.T) {
		numHelper := 0
		numErrorf := 0

		tb := newTestTB()
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		require.NotPanics(t, func() {
			testutil.AssertMarshalText(tb, "bad", &badCodec{value: []byte("bad")})
		})
		assert.Greater(t, numErrorf, 0)
		assert.Greater(t, numHelper, 0)
	})
}

func TestAssertUnmarshalText(t *testing.T) {
	t.Parallel()

	t.Run("good", func(t *testing.T) {
		numHelper := 0

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tb, "good", &goodCodec{value: []byte("good")})
		})
		assert.Greater(t, numHelper, 0)
	})

	t.Run("bad", func(t *testing.T) {
		numHelper := 0
		numErrorf := 0

		tb := newTestTB()
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		require.NotPanics(t, func() {
			testutil.AssertUnmarshalText(tb, "bad", &badCodec{value: []byte("bad")})
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

	tb := newTestTB()
	tb.onCleanup = func(f func()) {
		gotFunc = f
	}
	tb.onHelper = func() { numHelper++ }
	tb.onName = func() (name string) { return testName }

	testutil.CleanupAndRequireSuccess(tb, cleanupFunc)

	assert.Greater(t, numHelper, 0)

	gotFunc()
	assert.True(t, cleanupFuncCalled)
}

func TestRequireTypeAssert(t *testing.T) {
	t.Parallel()

	require.True(t, t.Run("concrete_fail", func(t *testing.T) {
		var numErrorf, numFailNow, numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onFailNow = func() { numFailNow++ }
		tb.onName = func() (name string) { return testName }

		const wantErrMsg = "interface conversion: interface {} is string, not int"

		var v any = ""
		assert.PanicsWithError(t, wantErrMsg, func() {
			_ = testutil.RequireTypeAssert[int](tb, v)
		})

		assert.Equal(t, 1, numErrorf)
		assert.Equal(t, 1, numFailNow)
		assert.Greater(t, numHelper, 1)
	}))

	require.True(t, t.Run("concrete_success", func(t *testing.T) {
		var numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		var v any = 1
		var got int
		assert.NotPanics(t, func() {
			got = testutil.RequireTypeAssert[int](tb, v)
		})

		assert.Equal(t, v, got)
		assert.Greater(t, numHelper, 1)
	}))

	require.True(t, t.Run("interface_fail", func(t *testing.T) {
		var numErrorf, numFailNow, numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }
		tb.onErrorf = func(_ string, _ ...any) { numErrorf++ }
		tb.onFailNow = func() { numFailNow++ }
		tb.onName = func() (name string) { return testName }

		const wantErrMsg = "interface conversion: string is not error: missing method Error"

		var v any = ""
		assert.PanicsWithError(t, wantErrMsg, func() {
			_ = testutil.RequireTypeAssert[error](tb, v)
		})

		assert.Equal(t, 1, numErrorf)
		assert.Equal(t, 1, numFailNow)
		assert.Greater(t, numHelper, 1)
	}))

	require.True(t, t.Run("interface_success", func(t *testing.T) {
		var numHelper int

		tb := newTestTB()
		tb.onHelper = func() { numHelper++ }
		tb.onName = func() (name string) { return testName }

		var v any = errors.Error("")
		var got error
		assert.NotPanics(t, func() {
			got = testutil.RequireTypeAssert[error](tb, v)
		})

		assert.Equal(t, v, got)
		assert.Greater(t, numHelper, 1)
	}))
}
