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

func TestValidateDomainName(t *testing.T) {
	t.Parallel()

	longDomainName := strings.Repeat("a", 255)
	longLabel := strings.Repeat("a", 64)
	longLabelDomainName := longLabel + ".com"

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

	longDomainName := strings.Repeat("a", 255)
	longLabel := strings.Repeat("a", 64)
	longLabelDomainName := longLabel + ".com"

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

	longDomainName := strings.Repeat("a", 255)
	longLabel := "_" + strings.Repeat("a", 16)
	longLabelDomainName := longLabel + ".com"

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
		in:        longLabelDomainName,
		wantErrAs: new(*netutil.LengthError),
		wantErrMsg: `bad service domain name "` + longLabelDomainName + `": ` +
			`bad service name label "` + longLabel + `": ` +
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
	testLongValidLabel = strings.Repeat("a", 63)
	testLongValidName  = strings.Repeat(testLongValidLabel+".", 3) + "com"
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
		in:   testLongValidName,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				errSink = netutil.ValidateDomainName(bc.in)
			}

			require.NoError(b, errSink)
		})
	}

	// goos: darwin
	// goarch: amd64
	// pkg: github.com/AdguardTeam/golibs/netutil
	// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	// BenchmarkValidateDomainName/common-12		10058613	109.4 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateDomainName/long_names-12	4830151		246.2 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateDomainName/long_labels-12	4775589		246.7 ns/op		0 B/op	0 allocs/op
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
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				errSink = netutil.ValidateSRVDomainName(bc.in)
			}

			require.NoError(b, errSink)
		})
	}

	// goos: darwin
	// goarch: amd64
	// pkg: github.com/AdguardTeam/golibs/netutil
	// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	// BenchmarkValidateSRVDomainName/common-12			8900175		132.9 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateSRVDomainName/long_names-12		5012017		236.0 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateSRVDomainName/long_labels-12	1534950		757.6 ns/op		0 B/op	0 allocs/op
}

func BenchmarkValidateHostname(b *testing.B) {
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
		in:   testLongValidName,
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				errSink = netutil.ValidateHostname(bc.in)
			}

			require.NoError(b, errSink)
		})
	}

	// goos: darwin
	// goarch: amd64
	// pkg: github.com/AdguardTeam/golibs/netutil
	// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	// BenchmarkValidateHostname/common-12			9037418		134.2 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateHostname/long_names-12		5069252		239.9 ns/op		0 B/op	0 allocs/op
	// BenchmarkValidateHostname/long_labels-12		1581854		765.9 ns/op		0 B/op	0 allocs/op
}
