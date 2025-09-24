// Package testutil contains utilities for common testing patterns.
package testutil

import (
	"testing"

	"github.com/AdguardTeam/golibs/internal/reflectutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertErrorMsg asserts that the error is not nil and that its message is
// equal to msg.  If msg is an empty string, AssertErrorMsg asserts that the
// error is nil instead.
func AssertErrorMsg(t assert.TestingT, msg string, err error) (ok bool) {
	if h, isHelper := t.(interface{ Helper() }); isHelper {
		h.Helper()
	}

	if msg == "" {
		return assert.NoError(t, err)
	}

	if !assert.Error(t, err) {
		return false
	}

	return assert.Equal(t, msg, err.Error())
}

// CleanupAndRequireSuccess sets a cleanup function which checks the error
// returned by f and fails the test using tb if there is one.
func CleanupAndRequireSuccess(tb testing.TB, f func() (err error)) {
	tb.Helper()

	tb.Cleanup(func() {
		err := f()
		require.NoError(tb, err)
	})
}

// RequireTypeAssert is a helper that first requires the desired type and then,
// if the type is correct, converts and returns the value.
func RequireTypeAssert[T any](t require.TestingT, v any) (res T) {
	if h, isHelper := t.(interface{ Helper() }); isHelper {
		h.Helper()
	}

	if reflectutil.IsInterface[T]() {
		require.Implements(t, (*T)(nil), v)

		return v.(T)
	}

	require.IsType(t, res, v)

	return v.(T)
}
