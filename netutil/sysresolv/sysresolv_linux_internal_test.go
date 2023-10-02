//go:build linux

package sysresolv

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemResolvers_DialFunc_dockerEmbeddedDNS(t *testing.T) {
	t.Parallel()

	sr, err := NewSystemResolvers(nil, 53)
	require.NoError(t, err)

	got, err := sr.parse(net.JoinHostPort(dockerEmbeddedDNS, "53"))
	assert.Zero(t, got)
	assert.ErrorIs(t, err, errFakeDial)
}
