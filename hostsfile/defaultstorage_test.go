package hostsfile_test

import (
	"io"
	"maps"
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

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	ds, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger:  testLogger,
			Readers: []io.Reader{f},
		},
	)
	require.NoError(t, err)

	t.Run("ByAddr", func(t *testing.T) {
		t.Parallel()

		for _, addr := range slices.SortedStableFunc(maps.Keys(wantAddrs), netip.Addr.Compare) {
			addr := addr

			t.Run(addr.String(), func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, wantAddrs[addr], ds.ByAddr(addr))
			})
		}
	})

	t.Run("ByHost", func(t *testing.T) {
		t.Parallel()

		for _, host := range slices.Sorted(maps.Keys(wantHosts)) {
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

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		ds, err := hostsfile.NewDefaultStorage(
			ctx,
			&hostsfile.DefaultStorageConfig{
				Logger:  testLogger,
				Readers: []io.Reader{f},
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, ds)

		for range ds.RangeAddrs {
			require.Fail(t, "should not be called")
		}
	})

	t.Run("reader", func(t *testing.T) {
		t.Parallel()

		r := &fakeio.Reader{
			OnRead: func(_ []byte) (n int, err error) {
				return 0, assert.AnError
			},
		}

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		ds, err := hostsfile.NewDefaultStorage(
			ctx,
			&hostsfile.DefaultStorageConfig{
				Logger:  testLogger,
				Readers: []io.Reader{r},
			},
		)
		require.ErrorIs(t, err, assert.AnError)

		assert.Nil(t, ds)
	})
}

func TestDefaultStorage_HandleInvalid(t *testing.T) {
	t.Parallel()

	ctx := testutil.ContextWithTimeout(t, testTimeout)
	ds, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger: testLogger,
		},
	)
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
			ds.HandleInvalid(ctx, tc.name, nil, tc.err)
		})
	}
}

func TestDefaultStorage_range(t *testing.T) {
	t.Parallel()

	const hostsStr = `` +
		"1.2.3.4 host.example another.example\n" +
		"4.3.2.1 yet.another.example\n"

	var (
		ctx = testutil.ContextWithTimeout(t, testTimeout)

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

	ds, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger:  testLogger,
			Readers: []io.Reader{strings.NewReader(hostsStr)},
		},
	)
	require.NoError(t, err)

	empty, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger: testLogger,
		},
	)
	require.NoError(t, err)

	t.Run("RangeAddrs", func(t *testing.T) {
		t.Parallel()

		names := maps.Clone(wantHosts)

		for name, addrs := range ds.RangeAddrs {
			got, ok := names[name]
			require.True(t, ok)
			require.Equal(t, got, addrs)

			delete(names, name)
		}

		require.Empty(t, names)
	})

	t.Run("RangeNames", func(t *testing.T) {
		t.Parallel()

		addrs := maps.Clone(wantAddrs)

		for addr, names := range ds.RangeNames {
			got, ok := addrs[addr]
			require.True(t, ok)
			require.Equal(t, got, names)

			delete(addrs, addr)
		}

		require.Empty(t, addrs)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		for range empty.RangeAddrs {
			assert.Fail(t, "should not be called")
		}

		for range empty.RangeNames {
			assert.Fail(t, "should not be called")
		}
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

	ctx := testutil.ContextWithTimeout(t, testTimeout)

	hs1, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger:  testLogger,
			Readers: []io.Reader{strings.NewReader(hosts1)},
		},
	)
	require.NoError(t, err)

	hs2, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger:  testLogger,
			Readers: []io.Reader{strings.NewReader(hosts2)},
		},
	)
	require.NoError(t, err)

	empty, err := hostsfile.NewDefaultStorage(
		ctx,
		&hostsfile.DefaultStorageConfig{
			Logger: testLogger,
		},
	)
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
