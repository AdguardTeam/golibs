package netutil_test

import (
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneMAC(t *testing.T) {
	t.Parallel()

	assert.Equal(t, net.HardwareAddr(nil), netutil.CloneMAC(nil))
	assert.Equal(t, net.HardwareAddr{}, netutil.CloneMAC(net.HardwareAddr{}))

	mac := net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}
	clone := netutil.CloneMAC(mac)
	assert.Equal(t, mac, clone)

	require.Len(t, clone, len(mac))

	assert.NotSame(t, &mac[0], &clone[0])
}

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
		wantErrMsg: `bad mac address "": address is empty`,
		wantErrAs:  new(errors.Error),
		in:         nil,
	}, {
		name:       "error_empty",
		wantErrMsg: `bad mac address "": address is empty`,
		wantErrAs:  new(errors.Error),
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
		wantPort   int
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
		in:         "example.com:12345",
		wantErrMsg: "",
		wantHost:   "example.com",
		wantPort:   12345,
	}, {
		name:       "bad_port",
		in:         "example.com:!!!",
		wantErrMsg: "parsing port: strconv.ParseUint: parsing \"!!!\": invalid syntax",
		wantHost:   "",
		wantPort:   0,
	}, {
		name:       "port_too_big",
		in:         "example.com:99999",
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
		in:         "example.com",
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
		wantErrMsg: `bad domain name "xn---.com": idna: invalid label "-"`,
	}, {
		name:       "success_one",
		in:         "e",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "empty",
		in:         "",
		wantErrAs:  new(errors.Error),
		wantErrMsg: `bad domain name "": address is empty`,
	}, {
		name:      "bad_symbol",
		in:        "!!!",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad domain name "!!!": ` +
			`bad domain name label "!!!": bad domain name label rune '!'`,
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
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad domain name "example..com": ` +
			`bad domain name label "": label is empty`,
	}, {
		name:      "bad_label_first_symbol",
		in:        "example.-aa.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad domain name "example.-aa.com": ` +
			`bad domain name label "-aa": bad domain name label rune '-'`,
	}, {
		name:      "bad_label_last_symbol",
		in:        "example-.aa.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad domain name "example-.aa.com": ` +
			`bad domain name label "example-": bad domain name label rune '-'`,
	}, {
		name:      "bad_label_symbol",
		in:        "example.a!!!.com",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad domain name "example.a!!!.com": ` +
			`bad domain name label "a!!!": bad domain name label rune '!'`,
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
		in:         "_http.example.com",
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
		in:         "_u",
		wantErrAs:  nil,
		wantErrMsg: "",
	}, {
		name:       "empty",
		in:         "",
		wantErrAs:  new(errors.Error),
		wantErrMsg: `bad service domain name "": address is empty`,
	}, {
		name:      "bad_symbol",
		in:        "_!",
		wantErrAs: new(*netutil.RuneError),
		wantErrMsg: `bad service domain name "_!": ` +
			`bad service name label "_!": bad service name label rune '!'`,
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
		wantErrAs: new(errors.Error),
		wantErrMsg: `bad service domain name "example..com": ` +
			`bad domain name label "": label is empty`,
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
		wantErrAs: new(*netutil.AddrError),
		wantErrMsg: `bad service domain name "example._.com": ` +
			`bad service name label "_": label is empty`,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.name == "bad_label_empty" {
				assert.True(t, true)
			}

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
		wantErrAs:  new(*netutil.AddrError),
		wantErrMsg: `bad service name label "": label is empty`,
		in:         "",
		name:       "empty",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := netutil.ValidateServiceNameLabel(tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)

			assert.ErrorAs(t, err, new(*netutil.AddrError))
			if tc.wantErrAs != nil {
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}
