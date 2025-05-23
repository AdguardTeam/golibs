package hostsfile_test

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/hostsfile"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validatingSet is a [hostsfile.Set] implementation that validates records in
// tests.
type validatingSet struct {
	// tb is used for assertions.
	tb testing.TB
}

// type check
var _ hostsfile.Set = validatingSet{}

// Add implements the [Set] interface for validatingSet.
func (s validatingSet) Add(r *hostsfile.Record) {
	validateRecord(s.tb, r, nil)
}

// type check
var _ hostsfile.HandleSet = validatingSet{}

// HandleInvalid implements the [HandleSet] interface for validatingSet.
func (s validatingSet) HandleInvalid(_ string, data []byte, err error) {
	rec := &hostsfile.Record{}

	var lineErr *hostsfile.LineError
	if errors.As(err, &lineErr) {
		err = lineErr.Unwrap()
	}
	require.Equal(s.tb, rec.UnmarshalText(data), err)

	validateRecord(s.tb, rec, err)
}

// validateRecord validates the given record considering the error as returned
// by [Record.UnmarshalText].
func validateRecord(t testing.TB, rec *hostsfile.Record, err error) {
	if helper, ok := t.(interface{ Helper() }); ok {
		helper.Helper()
	}

	const errPref = `ParseAddr("`

	var addrErr *netutil.AddrError
	switch {
	case errors.Is(err, hostsfile.ErrEmptyLine):
		// It's either a comment or an empty line.
		require.Nil(t, rec.Names)
		require.False(t, rec.Addr.IsValid())
	case errors.Is(err, hostsfile.ErrNoHosts):
		// The only field.
		require.Nil(t, rec.Names)
		require.False(t, rec.Addr.IsValid())
	case err != nil && strings.HasPrefix(err.Error(), errPref):
		// It's an invalid IP address.
		require.Nil(t, rec.Names)
		require.False(t, rec.Addr.IsValid())
	case errors.As(err, &addrErr):
		// It's a valid IP address, but some hostnames are invalid.
		require.NotNil(t, rec.Names)
		require.True(t, rec.Addr.IsValid())
	default:
		// Do not expect any other errors.
		require.NoError(t, err)

		require.True(t, rec.Addr.IsValid())
		require.NotEmpty(t, rec.Names)
	}
}

func TestRecord_UnmarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		data       string
		want       *hostsfile.Record
		wantErrMsg string
	}{{
		name:       "empty",
		data:       ``,
		want:       &hostsfile.Record{},
		wantErrMsg: "line is empty",
	}, {
		name:       "comment",
		data:       `# comment`,
		want:       &hostsfile.Record{},
		wantErrMsg: "line is empty",
	}, {
		name:       "comment_with_tab",
		data:       "\t# comment",
		want:       &hostsfile.Record{},
		wantErrMsg: "line is empty",
	}, {
		name:       "no_hosts",
		data:       "1.2.3.4 ",
		want:       &hostsfile.Record{},
		wantErrMsg: "no hostnames",
	}, {
		name:       "bad_addr",
		data:       `256.1.2.3 host1 host2`,
		want:       &hostsfile.Record{},
		wantErrMsg: `ParseAddr("256.1.2.3"): IPv4 field has value >255`,
	}, {
		name: "bad_hostname",
		data: "1.2.3.4 _host1",
		want: &hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{},
		},
		wantErrMsg: `name at index 0: bad domain name "_host1": ` +
			`bad top-level domain name label "_host1": ` +
			`bad top-level domain name label rune '_'`,
	}, {
		name: "bad_hostname_with_comment",
		data: "1.2.3.4 _host # this is bad host",
		want: &hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{},
		},
		wantErrMsg: `name at index 0: bad domain name "_host": ` +
			`bad top-level domain name label "_host": ` +
			`bad top-level domain name label rune '_'`,
	}, {
		name: "single_bad_host",
		data: "1.2.3.4 good.host bad._host",
		want: &hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{"good.host"},
		},
		wantErrMsg: `name at index 1: bad domain name "bad._host": ` +
			`bad top-level domain name label "_host": ` +
			`bad top-level domain name label rune '_'`,
	}, {
		name: "dot_host",
		data: ":: .",
		want: &hostsfile.Record{
			Addr:  netip.IPv6Unspecified(),
			Names: []string{},
		},
		wantErrMsg: `name at index 0: bad domain name ".": ` +
			`bad domain name label "": domain name label is empty`,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := &hostsfile.Record{}
			err := rec.UnmarshalText([]byte(tc.data))
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			assert.Equal(t, tc.want, rec)
		})
	}
}

func TestRecord_MarshalText(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		want []byte
		rec  hostsfile.Record
	}{{
		name: "empty",
		want: []byte{},
		rec:  hostsfile.Record{},
	}, {
		name: "no_hosts",
		want: []byte(testIPv4.String()),
		rec:  hostsfile.Record{Addr: testIPv4},
	}, {
		name: "single_host",
		want: []byte(testIPv4.String() + " host1"),
		rec: hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{"host1"},
		},
	}, {
		name: "multiple_hosts",
		want: []byte(testIPv4.String() + " host1 host2"),
		rec: hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{"host1", "host2"},
		},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := tc.rec.MarshalText()
			require.NoError(t, err)

			assert.Equal(t, tc.want, data)
		})
	}
}

func BenchmarkRecord_UnmarshalText(b *testing.B) {
	benchCases := []struct {
		name      string
		label     string
		labelsNum int
		hostsNum  int
	}{{
		name:      "two_hosts",
		label:     "label",
		labelsNum: 1,
		hostsNum:  2,
	}, {
		name:      "many_labels",
		label:     "label",
		labelsNum: 10,
		hostsNum:  2,
	}, {
		name:      "many_hosts",
		label:     "label",
		labelsNum: 1,
		hostsNum:  20,
	}, {
		name:      "many_labels_and_hosts",
		label:     "label",
		labelsNum: 10,
		hostsNum:  20,
	}, {
		name:      "two_large_hosts",
		label:     "really-wide-label-that-is-just-long-enough-to-fit-into-63-bytes",
		labelsNum: 4,
		hostsNum:  2,
	}, {
		name:      "many_large_hosts",
		label:     "really-wide-label-that-is-just-long-enough-to-fit-into-63-bytes",
		labelsNum: 4,
		hostsNum:  20,
	}, {
		name:      "many_hosts_tiny_labels",
		label:     "a",
		labelsNum: 1,
		hostsNum:  256,
	}, {
		name:      "many_hosts_many_tiny_labels",
		label:     "a",
		labelsNum: 127,
		hostsNum:  256,
	}}

	var rec hostsfile.Record
	for _, bc := range benchCases {
		host := strings.Repeat(bc.label+".", bc.labelsNum)[:len(bc.label)*bc.labelsNum]
		input := []byte(testIPv6.StringExpanded() + " " + strings.Repeat(host+" ", bc.hostsNum))

		b.Run(bc.name, func(b *testing.B) {
			var err error

			b.ReportAllocs()
			for b.Loop() {
				err = rec.UnmarshalText(input)
			}

			require.NoError(b, err)
		})

		b.Run(bc.name+"_with_allocs", func(b *testing.B) {
			b.Skip("Comment this line to run the benchmark with allocs")

			var err error

			b.ReportAllocs()
			for b.Loop() {
				err = rec.UnmarshalTextEachSublice(input)
			}

			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/hostsfile
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkRecord_UnmarshalText
	//	BenchmarkRecord_UnmarshalText/two_hosts
	//	BenchmarkRecord_UnmarshalText/two_hosts-16         	  993870	      1192 ns/op	      96 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_labels
	//	BenchmarkRecord_UnmarshalText/many_labels-16       	 1000000	      2577 ns/op	     192 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_hosts
	//	BenchmarkRecord_UnmarshalText/many_hosts-16        	  254593	      5107 ns/op	     496 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_labels_and_hosts
	//	BenchmarkRecord_UnmarshalText/many_labels_and_hosts-16         	   47509	     22245 ns/op	    1392 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/two_large_hosts
	//	BenchmarkRecord_UnmarshalText/two_large_hosts-16               	  259664	      3964 ns/op	     592 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_large_hosts
	//	BenchmarkRecord_UnmarshalText/many_large_hosts-16              	   54634	     29194 ns/op	    5744 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_hosts_tiny_labels
	//	BenchmarkRecord_UnmarshalText/many_hosts_tiny_labels-16        	   26436	     47916 ns/op	    5424 B/op	       3 allocs/op
	//	BenchmarkRecord_UnmarshalText/many_hosts_many_tiny_labels
	//	BenchmarkRecord_UnmarshalText/many_hosts_many_tiny_labels-16   	    1159	   1018552 ns/op	   37680 B/op	       3 allocs/op
}

func BenchmarkRecord_MarshalText(b *testing.B) {
	benchCases := []struct {
		name string
		rec  hostsfile.Record
	}{{
		name: "empty",
		rec:  hostsfile.Record{},
	}, {
		name: "ipv4_only",
		rec:  hostsfile.Record{Addr: testIPv4},
	}, {
		name: "ipv6_only",
		rec:  hostsfile.Record{Addr: testIPv6},
	}, {
		name: "ipv4_with_hosts",
		rec: hostsfile.Record{
			Addr:  testIPv4,
			Names: []string{"host1", "host2", "host3", "host4", "host5"},
		},
	}, {
		name: "ipv6_with_hosts",
		rec: hostsfile.Record{
			Addr:  testIPv6,
			Names: []string{"host1", "host2", "host3", "host4", "host5"},
		},
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = bc.rec.MarshalText()
			}
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/hostsfile
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkRecord_MarshalText
	//	BenchmarkRecord_MarshalText/empty
	//	BenchmarkRecord_MarshalText/empty-16         	167696192	         6.848 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkRecord_MarshalText/ipv4_only
	//	BenchmarkRecord_MarshalText/ipv4_only-16     	13294310	        79.69 ns/op	      16 B/op	       1 allocs/op
	//	BenchmarkRecord_MarshalText/ipv6_only
	//	BenchmarkRecord_MarshalText/ipv6_only-16     	 7429681	       211.1 ns/op	      48 B/op	       1 allocs/op
	//	BenchmarkRecord_MarshalText/ipv4_with_hosts
	//	BenchmarkRecord_MarshalText/ipv4_with_hosts-16         	 5413416	       206.0 ns/op	      64 B/op	       2 allocs/op
	//	BenchmarkRecord_MarshalText/ipv6_with_hosts
	//	BenchmarkRecord_MarshalText/ipv6_with_hosts-16         	 4154754	       308.1 ns/op	     144 B/op	       2 allocs/op
}

func FuzzRecord_UnmarshalText(f *testing.F) {
	for _, seed := range []string{
		"",
		"\n",
		"1.0.0.1 host1 host2",
		"1.0.0.1 host1.domain host2.domain",
		"127.0.0.1 localhost",
		"::1 localhost",
		"1234:5678:90ab:cdef:1234:5678:90ab:cdef host1 host2",
		"1234:5678:90ab:cdef:: host1.domain host2.domain",
		"256.256.256.256 bad.host",
		"fe80::1 localhost # comment",
		"fe80::1 # comment",
		"# comment",
		"1.2.3.4 -123-",
		"1.1.1.1 abc-",
		"1.1.1.1 -abc",
		"  1.2.3.4   spaced.hosts  ",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		data := []byte(input)

		rec := &hostsfile.Record{}
		err := rec.UnmarshalText(data)
		validateRecord(t, rec, err)

		// TODO(e.burkov):  Add MarshalText subtest.
	})
}
