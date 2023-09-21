package sysresolv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// NewTestResolvers is a helper for creating resolvers in tests.
func NewTestResolvers(t testing.TB, hostGenFunc HostGenFunc) (r Resolvers) {
	t.Helper()

	r, err := NewSystemResolvers(hostGenFunc)
	require.NoError(t, err)
	require.NotNil(t, r)

	return r
}
