package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPMap_allocs(t *testing.T) {
	t.Parallel()

	m := netutil.NewIPMap(0)
	m.Set(testIPv4, 42)

	t.Run("get", func(t *testing.T) {
		var v interface{}
		var ok bool
		allocs := testing.AllocsPerRun(100, func() {
			v, ok = m.Get(testIPv4)
		})

		require.True(t, ok)
		require.Equal(t, 42, v)

		assert.Equal(t, float64(0), allocs)
	})

	t.Run("len", func(t *testing.T) {
		var n int
		allocs := testing.AllocsPerRun(100, func() {
			n = m.Len()
		})

		require.Equal(t, 1, n)

		assert.Equal(t, float64(0), allocs)
	})
}

func TestIPMap(t *testing.T) {
	t.Parallel()

	val := 42

	t.Run("nil", func(t *testing.T) {
		var m *netutil.IPMap

		assert.NotPanics(t, func() {
			m.Clear()
		})

		assert.NotPanics(t, func() {
			m.Del(testIPv4)
			m.Del(testIPv6)
		})

		assert.NotPanics(t, func() {
			v, ok := m.Get(testIPv4)
			assert.Nil(t, v)
			assert.False(t, ok)

			v, ok = m.Get(testIPv6)
			assert.Nil(t, v)
			assert.False(t, ok)
		})

		assert.NotPanics(t, func() {
			assert.Equal(t, 0, m.Len())
		})

		assert.NotPanics(t, func() {
			n := 0
			m.Range(func(_ net.IP, _ interface{}) (cont bool) {
				n++

				return true
			})

			assert.Equal(t, 0, n)
		})

		assert.Panics(t, func() {
			m.Set(testIPv4, val)
		})

		assert.Panics(t, func() {
			m.Set(testIPv6, val)
		})

		assert.NotPanics(t, func() {
			sclone := m.ShallowClone()
			assert.Nil(t, sclone)
		})
	})

	testIPMap := func(t *testing.T, ip net.IP, s string) {
		m := netutil.NewIPMap(0)
		assert.Equal(t, 0, m.Len())

		v, ok := m.Get(ip)
		assert.Nil(t, v)
		assert.False(t, ok)

		m.Set(ip, val)
		v, ok = m.Get(ip)
		assert.Equal(t, val, v)
		assert.True(t, ok)

		n := 0
		m.Range(func(ipKey net.IP, v interface{}) (cont bool) {
			assert.Equal(t, ip.To16(), ipKey)
			assert.Equal(t, val, v)

			n++

			return false
		})
		assert.Equal(t, 1, n)

		sclone := m.ShallowClone()
		assert.Equal(t, m, sclone)

		assert.Equal(t, s, m.String())

		m.Del(ip)
		v, ok = m.Get(ip)
		assert.Nil(t, v)
		assert.False(t, ok)
		assert.Equal(t, 0, m.Len())
	}

	t.Run("ipv4", func(t *testing.T) {
		t.Parallel()

		testIPMap(t, testIPv4, "map[1.2.3.4:42]")
	})

	t.Run("ipv6", func(t *testing.T) {
		t.Parallel()

		testIPMap(t, testIPv6, "map[1234::cdef:42]")
	})
}
