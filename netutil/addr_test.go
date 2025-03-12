package netutil_test

import (
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exampleDomain = "example.com"

func TestCloneURL(t *testing.T) {
	t.Parallel()

	assert.Equal(t, (*url.URL)(nil), netutil.CloneURL(nil))
	assert.Equal(t, &url.URL{}, netutil.CloneURL(&url.URL{}))

	u, err := url.Parse("https://example.com/path?q=1&q=2#frag")
	require.NoError(t, err)

	clone := netutil.CloneURL(u)
	assert.Equal(t, u, clone)
	assert.NotSame(t, u, clone)
}

func TestValidateMAC(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		wantErrMsg string
		wantErrAs  any
		in         net.HardwareAddr
	}{{
		name:       "success_eui_48",
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.HardwareAddr{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
	}, {
		name:       "success_eui_64",
		wantErrMsg: "",
		wantErrAs:  nil,
		in:         net.HardwareAddr{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
	}, {
		name:       "success_infiniband",
		wantErrMsg: "",
		wantErrAs:  nil,
		in: net.HardwareAddr{
			0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
			0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x10, 0x11, 0x12, 0x13,
		},
	}, {
		name:       "error_nil",
		wantErrMsg: `bad mac address "": mac address is empty`,
		wantErrAs:  new(*netutil.LengthError),
		in:         nil,
	}, {
		name:       "error_empty",
		wantErrMsg: `bad mac address "": mac address is empty`,
		wantErrAs:  new(*netutil.LengthError),
		in:         net.HardwareAddr{},
	}, {
		name: "error_bad",
		wantErrMsg: `bad mac address "00:01:02:03": ` +
			`bad mac address length 4, allowed: [6 8 20]`,
		wantErrAs: new(*netutil.LengthError),
		in:        net.HardwareAddr{0x00, 0x01, 0x02, 0x03},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateMAC(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestJoinHostPort(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ":0", netutil.JoinHostPort("", 0))
	assert.Equal(t, "host:12345", netutil.JoinHostPort("host", 12345))
	assert.Equal(t, "1.2.3.4:12345", netutil.JoinHostPort("1.2.3.4", 12345))
	assert.Equal(t, "[1234::5678]:12345", netutil.JoinHostPort("1234::5678", 12345))
	assert.Equal(t, "[1234::5678]:12345", netutil.JoinHostPort("[1234::5678]", 12345))
	assert.Equal(t, "[1234::5678%lo]:12345", netutil.JoinHostPort("1234::5678%lo", 12345))
}

func TestSplitHostPort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantErrMsg string
		wantHost   string
		wantPort   uint16
	}{{
		name:       "success_ipv4",
		in:         "1.2.3.4:12345",
		wantErrMsg: "",
		wantHost:   "1.2.3.4",
		wantPort:   12345,
	}, {
		name:       "success_ipv6",
		in:         "[1234::5678]:12345",
		wantErrMsg: "",
		wantHost:   "1234::5678",
		wantPort:   12345,
	}, {
		name:       "success_ipv6_zone",
		in:         "[1234::5678%lo]:12345",
		wantErrMsg: "",
		wantHost:   "1234::5678%lo",
		wantPort:   12345,
	}, {
		name:       "success_host",
		in:         exampleDomain + ":12345",
		wantErrMsg: "",
		wantHost:   exampleDomain,
		wantPort:   12345,
	}, {
		name:       "bad_port",
		in:         exampleDomain + ":!!!",
		wantErrMsg: "parsing port: strconv.ParseUint: parsing \"!!!\": invalid syntax",
		wantHost:   "",
		wantPort:   0,
	}, {
		name:       "port_too_big",
		in:         exampleDomain + ":99999",
		wantErrMsg: "parsing port: strconv.ParseUint: parsing \"99999\": value out of range",
		wantHost:   "",
		wantPort:   0,
	}, {
		name:       "bad_syntax",
		in:         "[1234::5678:12345",
		wantErrMsg: "address [1234::5678:12345: missing ']' in address",
		wantHost:   "",
		wantPort:   0,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			host, port, err := netutil.SplitHostPort(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
			assert.Equal(t, tc.wantHost, host)
			assert.Equal(t, tc.wantPort, port)
		})
	}
}

var (
	longDomainName      = strings.Repeat("a", 255)
	longLabel           = strings.Repeat("a", 64)
	longLabelDomainName = longLabel + ".com"
)

func TestValidateDomainName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantErrAs  any
		wantErrMsg string
	}{{
		name:       "success",
		in:         exampleDomain,
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "success_idna",
		in:         "пример.рф",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "success_domain_name",
		in:         "_non-ldh-domain_.tld",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "bad_idna",
		in:         "xn---.com",
		wantErrAs:  nil,
		wantErrMsg: `bad domain name "xn---.com": idna: invalid label "-"`,
	}, {
		name:      "bad_tld",
		in:        exampleDomain + "-",
		wantErrAs: nil,
		wantErrMsg: `bad domain name "` + exampleDomain + `-": ` +
			`bad top-level domain name label "com-": ` +
			`bad top-level domain name label rune '-'`,
	}, {
		name:      "tld_too_long",
		in:        "example." + longLabel,
		wantErrAs: nil,
		wantErrMsg: `bad domain name "example.` + longLabel + `": ` +
			`bad top-level domain name label "` + longLabel + `": ` +
			`top-level domain name label is too long: got 64, max 63`,
	}, {
		name:      "numeric_tld",
		in:        "example.123",
		wantErrAs: new(*netutil.LabelError),
		wantErrMsg: `bad domain name "example.123": ` +
			`bad top-level domain name label "123": all octets are numeric`,
	}, {
		name:       "empty",
		in:         "",
		wantErrAs:  new(*netutil.LengthError),
		wantErrMsg: `bad domain name "": domain name is empty`,
	}, {
		name:      "tld_only",
		in:        "!!!",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad domain name "!!!": ` +
			`bad top-level domain name label "!!!": ` +
			`bad top-level domain name label rune '!'`,
	}, {
		name:      "bad_length",
		in:        longDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad domain name "` + longDomainName + `": ` +
			`domain name is too long: got 255, max 253`,
	}, {
		name:      "bad_label_length",
		in:        longLabelDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad domain name "` + longLabelDomainName + `": ` +
			`bad domain name label "` + longLabel + `": ` +
			`domain name label is too long: got 64, max 63`,
	}, {
		name:      "bad_label_empty",
		in:        "example..com",
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad domain name "example..com": ` +
			`bad domain name label "": domain name label is empty`,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateDomainName(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantErrAs  any
		wantErrMsg string
	}{{
		name:       "success",
		in:         exampleDomain,
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "success_idna",
		in:         "пример.рф",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "bad_idna",
		in:         "xn---.com",
		wantErrAs:  nil,
		wantErrMsg: `bad hostname "xn---.com": idna: invalid label "-"`,
	}, {
		name:       "success_one",
		in:         "e",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "empty",
		in:         "",
		wantErrAs:  new(*netutil.LengthError),
		wantErrMsg: `bad hostname "": hostname is empty`,
	}, {
		name:      "bad_symbol",
		in:        "!!!",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad hostname "!!!": ` +
			`bad top-level domain name label "!!!": ` +
			`bad top-level domain name label rune '!'`,
	}, {
		name:      "bad_length",
		in:        longDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad hostname "` + longDomainName + `": ` +
			`hostname is too long: got 255, max 253`,
	}, {
		name:      "bad_label_length",
		in:        longLabelDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad hostname "` + longLabelDomainName + `": ` +
			`bad hostname label "` + longLabel + `": ` +
			`hostname label is too long: got 64, max 63`,
	}, {
		name:      "bad_label_empty",
		in:        "example..com",
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad hostname "example..com": ` +
			`bad hostname label "": hostname label is empty`,
	}, {
		name:      "bad_label_first_symbol",
		in:        "example.-aa.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad hostname "example.-aa.com": ` +
			`bad hostname label "-aa": bad hostname label rune '-'`,
	}, {
		name:      "bad_label_last_symbol",
		in:        "example-.aa.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad hostname "example-.aa.com": ` +
			`bad hostname label "example-": bad hostname label rune '-'`,
	}, {
		name:      "bad_label_symbol",
		in:        "example.a!!!.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad hostname "example.a!!!.com": ` +
			`bad hostname label "a!!!": bad hostname label rune '!'`,
	}, {
		name:      "numeric_tld",
		in:        "example.123",
		wantErrAs: new(*netutil.LabelError),
		wantErrMsg: `bad hostname "example.123": ` +
			`bad top-level domain name label "123": all octets are numeric`,
	}, {
		name:      "bad_tld",
		in:        "example._bad",
		wantErrAs: new(*netutil.LabelError),
		wantErrMsg: `bad hostname "example._bad": ` +
			`bad top-level domain name label "_bad": ` +
			`bad top-level domain name label rune '_'`,
	}, {
		name:      "too_long_tld",
		in:        "example." + longLabel,
		wantErrAs: new(*netutil.LabelError),
		wantErrMsg: `bad hostname "example.` + longLabel + `": ` +
			`bad top-level domain name label "` + longLabel + `": ` +
			`top-level domain name label is too long: got 64, max 63`,
	}, {
		name:      "empty_tld",
		in:        "example.",
		wantErrAs: new(*netutil.LabelError),
		wantErrMsg: `bad hostname "example.": ` +
			`bad top-level domain name label "": ` +
			`top-level domain name label is empty`,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateHostname(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestValidateSRVDomainName(t *testing.T) {
	t.Parallel()

	longSRVLabel := "_" + strings.Repeat("a", 16)
	longSRVLabelDomainName := longSRVLabel + ".com"

	testCases := []struct {
		name       string
		in         string
		wantErrAs  any
		wantErrMsg string
	}{{
		name:       "success",
		in:         "_http." + exampleDomain,
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "success_idna",
		in:         "_http.пример.рф",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:      "bad_idna",
		in:        "xn---.com",
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad service domain name "xn---.com": ` +
			`idna: invalid label "-"`,
	}, {
		name:       "success_one",
		in:         "_u.tld",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "empty",
		in:         "",
		wantErrAs:  new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "": service domain name is empty`,
	}, {
		name:      "bad_symbol",
		in:        "_!",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service domain name "_!": ` +
			`bad top-level domain name label "_!": ` +
			`bad top-level domain name label rune '_'`,
	}, {
		name:      "bad_length",
		in:        longDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "` + longDomainName + `": ` +
			`service domain name is too long: got 255, max 253`,
	}, {
		name:      "bad_label_length",
		in:        longSRVLabelDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "` + longSRVLabelDomainName + `": ` +
			`bad service name label "` + longSRVLabel + `": ` +
			`service name label is too long: got 17, max 16`,
	}, {
		name:      "bad_label_empty",
		in:        "example..com",
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "example..com": ` +
			`bad hostname label "": hostname label is empty`,
	}, {
		name:      "bad_label_first_symbol",
		in:        "example._-a.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service domain name "example._-a.com": ` +
			`bad service name label "_-a": bad service name label rune '-'`,
	}, {
		name:      "bad_label_last_symbol",
		in:        "_example-.aa.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service domain name "_example-.aa.com": ` +
			`bad service name label "_example-": bad service name label rune '-'`,
	}, {
		name:      "bad_label_unexpected_underscore",
		in:        "example._ht_tp.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service domain name "example._ht_tp.com": ` +
			`bad service name label "_ht_tp": bad service name label rune '_'`,
	}, {
		name:      "bad_service_label_empty",
		in:        "example._.com",
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "example._.com": ` +
			`bad service name label "_": service name label is empty`,
	}, {
		name:      "bad_hostname_label",
		in:        "-srv.com",
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad service domain name "-srv.com": ` +
			`bad hostname label "-srv": bad hostname label rune '-'`,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateSRVDomainName(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			if tc.wantErrAs != nil {
				require.Error(t, err)

				assert.ErrorAs(t, err, new(*netutil.AddrError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestValidateServiceNameLabel_errors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		wantErrAs  any
		wantErrMsg string
		in         string
		name       string
	}{{
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service name label "non-service.com": ` +
			`bad service name label rune 'n'`,
		in:   "non-service.com",
		name: "bad_rune",
	}, {
		wantErrAs:  new(*netutil.LabelError),
		wantErrMsg: `bad service name label "": service name label is empty`,
		in:         "",
		name:       "empty",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateServiceNameLabel(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			assert.ErrorAs(t, err, new(*netutil.LabelError))
			assert.ErrorAs(t, err, tc.wantErrAs)
		})
	}
}

// Common long test cases for benchmarking.
var (
	testLongValidLabel    = strings.Repeat("a", 63)
	testLongValidHostname = strings.Repeat(testLongValidLabel+".", 3) + "com"

	testLongInvalidLabel    = strings.Repeat("a", 62) + "!"
	testLongInvalidHostname = strings.Repeat(testLongValidLabel+".", 3) + "123"
)

func BenchmarkValidateDomainName(b *testing.B) {
	benchCases := []struct {
		name string
		in   string
	}{{
		name: "common",
		in:   exampleDomain,
	}, {
		name: "long_names",
		in:   testLongValidLabel,
	}, {
		name: "long_labels",
		in:   testLongValidHostname,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for b.Loop() {
				err = netutil.ValidateDomainName(bc.in)
			}

			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkValidateDomainName
	//	BenchmarkValidateDomainName/common
	//	BenchmarkValidateDomainName/common-16         	13883317	        85.97 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateDomainName/long_names
	//	BenchmarkValidateDomainName/long_names-16     	 8274099	       142.7 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateDomainName/long_labels
	//	BenchmarkValidateDomainName/long_labels-16    	 5193171	       229.5 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkIsValidHostnameLabel(b *testing.B) {
	benchCases := []struct {
		want assert.BoolAssertionFunc
		name string
		in   string
	}{{
		want: assert.True,
		name: "valid_short",
		in:   "label",
	}, {
		want: assert.True,
		name: "valid_long",
		in:   testLongValidLabel,
	}, {
		want: assert.False,
		name: "invalid_short",
		in:   "_label",
	}, {
		want: assert.False,
		name: "invalid_long",
		in:   testLongInvalidLabel,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var got bool
			b.ReportAllocs()
			for b.Loop() {
				got = netutil.IsValidHostnameLabel(bc.in)
			}

			bc.want(b, got)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkIsValidHostnameLabel
	//	BenchmarkIsValidHostnameLabel/valid_short
	//	BenchmarkIsValidHostnameLabel/valid_short-16         	226464375	         5.114 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostnameLabel/valid_long
	//	BenchmarkIsValidHostnameLabel/valid_long-16          	24387102	        47.21 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostnameLabel/invalid_short
	//	BenchmarkIsValidHostnameLabel/invalid_short-16       	404997958	         2.557 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostnameLabel/invalid_long
	//	BenchmarkIsValidHostnameLabel/invalid_long-16        	21934094	        51.47 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkValidateSRVDomainName(b *testing.B) {
	benchCases := []struct {
		name string
		in   string
	}{{
		name: "common",
		in:   exampleDomain,
	}, {
		name: "long_names",
		in:   testLongValidLabel,
	}, {
		name: "long_labels",
		in:   "_" + strings.Repeat("a", 15) + strings.Repeat("."+testLongValidLabel, 3),
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for b.Loop() {
				err = netutil.ValidateSRVDomainName(bc.in)
			}

			require.NoError(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkValidateSRVDomainName
	//	BenchmarkValidateSRVDomainName/common
	//	BenchmarkValidateSRVDomainName/common-16         	12360364	        97.25 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateSRVDomainName/long_names
	//	BenchmarkValidateSRVDomainName/long_names-16     	 8374941	       141.2 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateSRVDomainName/long_labels
	//	BenchmarkValidateSRVDomainName/long_labels-16    	 2283667	       483.9 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkValidateHostname(b *testing.B) {
	benchCases := []struct {
		want require.ErrorAssertionFunc
		name string
		in   string
	}{{
		want: require.NoError,
		name: "common",
		in:   "domain.example",
	}, {
		want: require.NoError,
		name: "good_short",
		in:   "abc.xyz",
	}, {
		want: require.NoError,
		name: "good_long",
		in:   testLongValidHostname,
	}, {
		want: require.NoError,
		name: "good_idna",
		in:   "международный.пример",
	}, {
		want: require.Error,
		name: "bad_short",
		in:   "!!!",
	}, {
		want: require.Error,
		name: "bad_long",
		in:   testLongInvalidHostname,
	}, {
		want: require.Error,
		name: "bad_idna",
		in:   "xn---.com",
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for b.Loop() {
				err = netutil.ValidateHostname(bc.in)
			}

			bc.want(b, err)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkValidateHostname
	//	BenchmarkValidateHostname/common
	//	BenchmarkValidateHostname/common-16         	11950093	        98.83 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateHostname/good_short
	//	BenchmarkValidateHostname/good_short-16     	12399717	        94.78 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateHostname/good_long
	//	BenchmarkValidateHostname/good_long-16      	 2437651	       458.3 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkValidateHostname/good_idna
	//	BenchmarkValidateHostname/good_idna-16      	  237970	      4357 ns/op	     216 B/op	       6 allocs/op
	//	BenchmarkValidateHostname/bad_short
	//	BenchmarkValidateHostname/bad_short-16      	 2553793	       462.6 ns/op	     168 B/op	       4 allocs/op
	//	BenchmarkValidateHostname/bad_long
	//	BenchmarkValidateHostname/bad_long-16       	  958098	      1228 ns/op	      96 B/op	       2 allocs/op
	//	BenchmarkValidateHostname/bad_idna
	//	BenchmarkValidateHostname/bad_idna-16       	 4453160	       259.4 ns/op	      80 B/op	       2 allocs/op
}

func BenchmarkIsValidHostname(b *testing.B) {
	benchCases := []struct {
		want require.BoolAssertionFunc
		name string
		in   string
	}{{
		want: require.True,
		name: "common",
		in:   "domain.example",
	}, {
		want: require.True,
		name: "good_short",
		in:   "abc.xyz",
	}, {
		want: require.True,
		name: "good_long",
		in:   testLongValidHostname,
	}, {
		want: require.True,
		name: "good_idna",
		in:   "международный.пример",
	}, {
		want: require.False,
		name: "bad_short",
		in:   "!!!",
	}, {
		want: require.False,
		name: "bad_long",
		in:   testLongInvalidHostname,
	}, {
		want: require.False,
		name: "bad_idna",
		in:   "xn---.com",
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var got bool
			b.ReportAllocs()
			for b.Loop() {
				got = netutil.IsValidHostname(bc.in)
			}

			bc.want(b, got)
		})
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkIsValidHostname
	//	BenchmarkIsValidHostname/common
	//	BenchmarkIsValidHostname/common-16         	15910483	        75.20 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostname/good_short
	//	BenchmarkIsValidHostname/good_short-16     	15677844	        74.37 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostname/good_long
	//	BenchmarkIsValidHostname/good_long-16      	 3425068	       344.1 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostname/good_idna
	//	BenchmarkIsValidHostname/good_idna-16      	  311358	      3992 ns/op	     216 B/op	       6 allocs/op
	//	BenchmarkIsValidHostname/bad_short
	//	BenchmarkIsValidHostname/bad_short-16      	27769028	        40.98 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostname/bad_long
	//	BenchmarkIsValidHostname/bad_long-16       	 3439100	       348.4 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidHostname/bad_idna
	//	BenchmarkIsValidHostname/bad_idna-16       	 5983546	       180.6 ns/op	      32 B/op	       1 allocs/op
}

func FuzzIsValidHostname(f *testing.F) {
	for _, seed := range []string{
		"",
		" ",
		"\n",
		exampleDomain,
		"пример.рф",
		"xn---.com",
		"e",
		"!!!",
		longDomainName,
		longLabelDomainName,
		"example..com",
		"example.-aa.com",
		"example-.aa.com",
		"example.a!!!.com",
		"example.123",
		"example._bad",
		"example." + longLabel,
		"example.",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ok := netutil.IsValidHostname(input)
		err := netutil.ValidateHostname(input)

		require.Equal(t, err == nil, ok)
	})
}

func FuzzIsValidHostnameLabel(f *testing.F) {
	for _, seed := range []string{
		"",
		" ",
		"\n",
		exampleDomain,
		testLongValidHostname,
		testLongValidLabel,
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ok := netutil.IsValidHostnameLabel(input)
		err := netutil.ValidateHostnameLabel(input)

		require.Equal(t, err == nil, ok)
	})
}
