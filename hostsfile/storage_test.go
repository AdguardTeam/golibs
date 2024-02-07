package hostsfile_test

import (
	"net/netip"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestDefaultStorage_lookup(t *testing.T) {
	t.Parallel()

	// Variables mirroring the testdata/TestDefaultStorage/*/hosts file.
	var (
		v4Addr1 = netip.MustParseAddr("0.0.0.1")
		v4Addr2 = netip.MustParseAddr("0.0.0.2")

		mappedAddr1 = netip.MustParseAddr("::ffff:0.0.0.1")
		mappedAddr2 = netip.MustParseAddr("::ffff:0.0.0.2")

		v6Addr1 = netip.MustParseAddr("::1")
		v6Addr2 = netip.MustParseAddr("::2")

		wantHosts = map[string][]netip.Addr{
			"host.one":       {v4Addr1, mappedAddr1, v6Addr1},
			"host.two":       {v4Addr2, mappedAddr2, v6Addr2},
			"host.new":       {v4Addr2, v4Addr1, mappedAddr2, mappedAddr1, v6Addr2, v6Addr1},
			"again.host.two": {v4Addr2, mappedAddr2, v6Addr2},
		}

		wantAddrs = map[netip.Addr][]string{
			v4Addr1:     {"Host.One", "host.new"},
			v4Addr2:     {"Host.Two", "Host.New", "Again.Host.Two"},
			mappedAddr1: {"Host.One", "host.new"},
			mappedAddr2: {"Host.Two", "Host.New", "Again.Host.Two"},
			v6Addr1:     {"Host.One", "host.new"},
			v6Addr2:     {"Host.Two", "Host.New", "Again.Host.Two"},
		}
	)

	f, err := testdata.Open(path.Join(t.Name(), "hosts"))
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, f.Close)

	ds, err := hostsfile.NewDefaultStorage(f)
	require.NoError(t, err)

	t.Run("ByAddr", func(t *testing.T) {
		t.Parallel()

		// Sort keys to make the test deterministic.
		addrs := maps.Keys(wantAddrs)
		slices.SortFunc(addrs, netip.Addr.Compare)

		for _, addr := range addrs {
			addr := addr

			t.Run(addr.String(), func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, wantAddrs[addr], ds.ByAddr(addr))
			})
		}
	})

	t.Run("ByHost", func(t *testing.T) {
		t.Parallel()

		// Sort keys to make the test deterministic.
		hosts := maps.Keys(wantHosts)
		slices.Sort(hosts)

		for _, host := range hosts {
			host := host

			t.Run(host, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, wantHosts[host], ds.ByName(host))
			})
		}
	})
}

func TestNewDefaultStorage_bad(t *testing.T) {
	t.Parallel()

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		f, err := testdata.Open(path.Join(t.Name(), "hosts"))
		require.NoError(t, err)
		testutil.CleanupAndRequireSuccess(t, f.Close)

		ds, err := hostsfile.NewDefaultStorage(f)
		require.NoError(t, err)
		assert.NotNil(t, ds)

		ds.RangeAddrs(func(_ string, _ []netip.Addr) (ok bool) {
			require.Fail(t, "should not be called")

			return false
		})
	})

	t.Run("reader", func(t *testing.T) {
		t.Parallel()

		r := &fakeio.Reader{
			OnRead: func(_ []byte) (n int, err error) {
				return 0, assert.AnError
			},
		}

		ds, err := hostsfile.NewDefaultStorage(r)
		require.ErrorIs(t, err, assert.AnError)

		assert.Nil(t, ds)
	})
}

func TestDefaultStorage_HandleInvalid(t *testing.T) {
	t.Parallel()

	ds, err := hostsfile.NewDefaultStorage()
	require.NoError(t, err)

	testCases := []struct {
		err  error
		name string
	}{{
		name: "empty_line",
		err:  hostsfile.ErrEmptyLine,
	}, {
		name: "unexpected_error",
		err:  assert.AnError,
	}, {
		name: "line_error",
		err:  &hostsfile.LineError{},
	}}

	for _, tc := range testCases {
		tc := tc

		assert.NotPanics(t, func() {
			ds.HandleInvalid(tc.name, nil, tc.err)
		})
	}
}

func TestDefaultStorage_range(t *testing.T) {
	t.Parallel()

	const hostsStr = `` +
		"1.2.3.4 host.example another.example\n" +
		"4.3.2.1 yet.another.example\n"

	var (
		v4Addr1 = netip.MustParseAddr("1.2.3.4")
		v4Addr2 = netip.MustParseAddr("4.3.2.1")

		wantHosts = map[string][]netip.Addr{
			"host.example":        {v4Addr1},
			"another.example":     {v4Addr1},
			"yet.another.example": {v4Addr2},
		}

		wantAddrs = map[netip.Addr][]string{
			v4Addr1: {"host.example", "another.example"},
			v4Addr2: {"yet.another.example"},
		}
	)

	ds, err := hostsfile.NewDefaultStorage(strings.NewReader(hostsStr))
	require.NoError(t, err)

	empty, err := hostsfile.NewDefaultStorage()
	require.NoError(t, err)

	t.Run("RangeAddrs", func(t *testing.T) {
		t.Parallel()

		names := maps.Clone(wantHosts)

		ds.RangeAddrs(func(name string, addrs []netip.Addr) (ok bool) {
			got, ok := names[name]
			require.True(t, ok)
			require.Equal(t, got, addrs)

			delete(names, name)

			return len(names) > 0
		})

		require.Empty(t, names)
	})

	t.Run("RangeNames", func(t *testing.T) {
		t.Parallel()

		addrs := maps.Clone(wantAddrs)

		ds.RangeNames(func(addr netip.Addr, names []string) (ok bool) {
			got, ok := addrs[addr]
			require.True(t, ok)
			require.Equal(t, got, names)

			delete(addrs, addr)

			return len(addrs) > 0
		})

		require.Empty(t, addrs)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		empty.RangeAddrs(func(_ string, _ []netip.Addr) (ok bool) {
			return assert.Fail(t, "should not be called")
		})

		empty.RangeNames(func(_ netip.Addr, _ []string) (ok bool) {
			return assert.Fail(t, "should not be called")
		})
	})
}

func TestDefaultStorage_Equal(t *testing.T) {
	t.Parallel()

	const hosts1 = `` +
		"1.2.3.4 host.example another.example\n" +
		"4.3.2.1 yet.another.example\n"

	const hosts2 = `` +
		"5.6.7.8 host.example another.example\n" +
		"8.7.6.5 yet.another.example\n"

	hs1, err := hostsfile.NewDefaultStorage(strings.NewReader(hosts1))
	require.NoError(t, err)

	hs2, err := hostsfile.NewDefaultStorage(strings.NewReader(hosts2))
	require.NoError(t, err)

	empty, err := hostsfile.NewDefaultStorage()
	require.NoError(t, err)

	testCases := []struct {
		a    *hostsfile.DefaultStorage
		b    *hostsfile.DefaultStorage
		want assert.BoolAssertionFunc
		name string
	}{{
		name: "equal",
		a:    hs1,
		b:    hs1,
		want: assert.True,
	}, {
		name: "not_equal",
		a:    hs1,
		b:    hs2,
		want: assert.False,
	}, {
		name: "nils",
		a:    nil,
		b:    nil,
		want: assert.True,
	}, {
		name: "nil_receiver",
		a:    nil,
		b:    empty,
		want: assert.False,
	}, {
		name: "nil_argument",
		a:    empty,
		b:    nil,
		want: assert.False,
	}, {
		name: "empty",
		a:    empty,
		b:    empty,
		want: assert.True,
	}, {
		name: "one_empty",
		a:    empty,
		b:    hs1,
		want: assert.False,
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.want(t, tc.a.Equal(tc.b))
		})
	}
}
