//go:build unix

package sysresolv

import (
	"context"
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestResolvers returns an underlying implementation of the default
// resolvers checker.
func newTestResolvers(t *testing.T, hostGenFunc HostGenFunc) (sr *systemResolvers) {
	t.Helper()

	r := NewTestResolvers(t, hostGenFunc)

	return testutil.RequireTypeAssert[*systemResolvers](t, r)
}

func TestSystemResolvers_Refresh(t *testing.T) {
	t.Run("expected_error", func(t *testing.T) {
		sr := newTestResolvers(t, nil)

		assert.NoError(t, sr.Refresh())
	})

	t.Run("unexpected_error", func(t *testing.T) {
		_, err := NewSystemResolvers(func() (host string) {
			return "127.0.0.1::123"
		})

		assert.Error(t, err)
	})
}

func TestSystemResolvers_DialFunc(t *testing.T) {
	imp := newTestResolvers(t, nil)

	testCases := []struct {
		want    error
		name    string
		address string
	}{{
		want:    errFakeDial,
		name:    "valid_ipv4",
		address: "127.0.0.1",
	}, {
		want:    errFakeDial,
		name:    "docker_ipv4_port",
		address: net.JoinHostPort(dockerEmbeddedDNS, "53"),
	}, {
		want:    errFakeDial,
		name:    "valid_ipv6_port",
		address: "[::1]:53",
	}, {
		want:    errFakeDial,
		name:    "valid_ipv6_zone_port",
		address: "[::1%lo0]:53",
	}, {
		want:    errBadAddrPassed,
		name:    "invalid_split_host",
		address: "127.0.0.1::123",
	}, {
		want:    errUnexpectedHostFormat,
		name:    "invalid_ipv6_zone_port",
		address: "[::1%%lo0]:53",
	}, {
		want:    errBadAddrPassed,
		name:    "invalid_parse_ip",
		address: "not-ip",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := imp.dialFunc(context.Background(), "", tc.address)
			require.Nil(t, conn)

			assert.ErrorIs(t, err, tc.want)
		})
	}
}
