package hostsfile_test

import (
	"bytes"
	"fmt"
	"io"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// testIPv4 is an IPv4 address common for tests.  Do not mutate.
	testIPv4 = netip.AddrFrom4([4]byte{1, 2, 3, 4})

	// testIPv6 is an IPv6 address common for tests.  Do not mutate.
	testIPv6 = netip.AddrFrom16([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
)

// testdata is a virtual filesystem with data for testing purposes.
var testdata = os.DirFS("./testdata")

// sliceSet is a [hostsfile.Set] implementation based that stores records in a
// slice.
type sliceSet []hostsfile.Record

// Add implements the [Set] interface for *sliceSet.
func (s *sliceSet) Add(r *hostsfile.Record) {
	*s = append(*s, *r)
}

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		source     io.Reader
		wantErrMsg string
		want       []hostsfile.Record
	}{{
		name:       "empty",
		source:     strings.NewReader(``),
		want:       nil,
		wantErrMsg: "",
	}, {
		name:       "empty_line",
		source:     strings.NewReader("\n"),
		want:       nil,
		wantErrMsg: `parsing: line 1: line is empty`,
	}, {
		name:       "comment_line",
		source:     strings.NewReader(`# comment`),
		want:       nil,
		wantErrMsg: `parsing: line 1: line is empty`,
	}, {
		name:       "no_hosts",
		source:     strings.NewReader(`1.2.3.4 `),
		want:       nil,
		wantErrMsg: `parsing: line 1: no hostnames`,
	}, {
		name:   "single_record",
		source: strings.NewReader(`1.2.3.4 host1 host2`),
		want: []hostsfile.Record{{
			Addr:  testIPv4,
			Names: []string{"host1", "host2"},
		}},
		wantErrMsg: "",
	}, {
		name: "with_comment",
		source: strings.NewReader(`
			# comment
			1.2.3.4 host1 host2`,
		),
		want: []hostsfile.Record{{
			Addr:  testIPv4,
			Names: []string{"host1", "host2"},
		}},
		wantErrMsg: "parsing: line 1: line is empty\nline 2: line is empty",
	}, {
		name: "two_records",
		source: strings.NewReader(`
			1.2.3.4 host1 host2
			4.3.2.1 host3 host4`,
		),
		want: []hostsfile.Record{{
			Addr:  testIPv4,
			Names: []string{"host1", "host2"},
		}, {
			Addr:  netip.MustParseAddr("4.3.2.1"),
			Names: []string{"host3", "host4"},
		}},
		wantErrMsg: "parsing: line 1: line is empty",
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var recs sliceSet
			err := hostsfile.Parse(&recs, tc.source, nil)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			assert.Equal(t, tc.want, []hostsfile.Record(recs))
		})
	}
}

func TestParse_fileSource(t *testing.T) {
	t.Parallel()

	f, err := testdata.Open("named_hosts")
	require.NoError(t, err)
	testutil.CleanupAndRequireSuccess(t, f.Close)

	recs := sliceSet{}
	err = hostsfile.Parse(&recs, f, nil)
	require.NoError(t, err)
	require.NotEmpty(t, recs)

	rec := recs[0]
	assert.Equal(t, testIPv4, rec.Addr)
	assert.Equal(t, []string{"host1", "host2"}, rec.Names)

	_, fileName := filepath.Split(rec.Source)
	assert.Equal(t, "named_hosts", fileName)
}

func TestParse_badReader(t *testing.T) {
	t.Parallel()

	const readErr errors.Error = "reading error"

	r := &fakeio.Reader{
		OnRead: func(p []byte) (n int, err error) {
			return 0, readErr
		},
	}

	err := hostsfile.Parse(hostsfile.DiscardSet{}, r, nil)
	require.ErrorIs(t, err, readErr)
}

func BenchmarkParse(b *testing.B) {
	// linesNum defines the number of line in the file for benchmarking.
	const linesNum = 1024

	data := &bytes.Buffer{}
	for i, addr := 0, netip.MustParseAddr("0.0.0.0"); i < linesNum; i++ {
		addr = addr.Next()
		fmt.Fprintf(data, "%s host%d\n", addr, i)
	}

	// Length of lines shouldn't exceed 64 bytes.
	buf := make([]byte, 0, 64)
	set := hostsfile.DiscardSet{}

	b.Run("run", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var err error
		for i := 0; i < b.N; i++ {
			err = hostsfile.Parse(set, data, buf)
		}

		require.NoError(b, err)
	})

	// goos: darwin
	// goarch: amd64
	// pkg: github.com/AdguardTeam/golibs/hostsfile
	// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	// BenchmarkParse/run-12		37397535		33.20 ns/op		0 B/op		0 allocs/op
}

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, seed []byte) {
		err := hostsfile.Parse(&validatingSet{tb: t}, bytes.NewReader(seed), nil)
		if err != nil {
			// TODO(e.burkov):  Check each error when it is possible to unwrap
			// those after migration to errors.Join.
			require.True(t, strings.HasPrefix(err.Error(), "parsing: "))
		}
	})
}
