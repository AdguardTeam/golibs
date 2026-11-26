package arp_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/arp"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatic_Table(t *testing.T) {
	t.Parallel()

	const (
		testHostV4 = "v4.test.example"
		testHostV6 = "v6.test.example"
	)

	t.Run("non_empty", func(t *testing.T) {
		t.Parallel()

		wantNeighbors := []*arp.Neighbor{{
			Name: testHostV4,
			IP:   testIPv4,
			MAC:  net.HardwareAddr{0xAB, 0xCD, 0xEF, 0xAB, 0xCD, 0xEF},
		}, {
			Name: testHostV6,
			IP:   testIPv6,
			MAC:  net.HardwareAddr{0xEF, 0xCD, 0xAB, 0xEF, 0xCD, 0xAB},
		}}

		s := arp.NewStatic(&arp.StaticConfig{
			Neighbors: wantNeighbors,
		})

		got, err := s.Table(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Equal(t, wantNeighbors, got.Neighbors)

		// Verify that the underlying slice was deeply copied.
		got.Neighbors[0].Name = "foo.example"
		assert.Equal(t, testHostV4, wantNeighbors[0].Name)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		s := arp.NewStatic(nil)

		got, err := s.Table(testutil.ContextWithTimeout(t, testTimeout))
		require.NoError(t, err)
		require.NotNil(t, got)

		assert.Empty(t, got.Neighbors)
	})
}
