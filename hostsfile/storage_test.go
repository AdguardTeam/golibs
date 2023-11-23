package hostsfile_test

import (
	"io/fs"
	"net/netip"
	"path"
	"testing"

	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func TestNewDefaultStorage(t *testing.T) {
	t.Parallel()

	var ds *hostsfile.DefaultStorage
	var err error

	t.Run("good_file", func(t *testing.T) {
		var f fs.File

		f, err = testdata.Open(path.Join(t.Name(), "hosts"))
		require.NoError(t, err)
		testutil.CleanupAndRequireSuccess(t, f.Close)

		ds, err = hostsfile.NewDefaultStorage(f)
	})
	require.NoError(t, err)

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

	t.Run("RangeNames", func(t *testing.T) {
		t.Parallel()

		ds.RangeNames(func(addr netip.Addr, names []string) {
			assert.Equal(t, wantAddrs[addr], names)
		})
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

	t.Run("RangeAddrs", func(t *testing.T) {
		t.Parallel()

		ds.RangeAddrs(func(name string, addrs []netip.Addr) {
			assert.Equal(t, wantHosts[name], addrs)
		})
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

		ds.RangeAddrs(func(_ string, _ []netip.Addr) {
			require.Fail(t, "should not be called")
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
