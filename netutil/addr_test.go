package netutil_test

import (
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/stringutil"
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
		wantErrAs  interface{}
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
	assert.Equal(t, "[1234::5678%lo]:12345", netutil.JoinHostPort("1234::5678%lo", 12345))
}

func TestSplitHostPort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         string
		wantHost   string
		wantErrMsg string
		wantPort   int
	}{{
		name:       "success_ipv4",
		in:         "1.2.3.4:12345",
		wantHost:   "1.2.3.4",
		wantErrMsg: "",
		wantPort:   12345,
	}, {
		name:       "success_ipv6",
		in:         "[1234::5678]:12345",
		wantHost:   "1234::5678",
		wantErrMsg: "",
		wantPort:   12345,
	}, {
		name:       "success_ipv6_zone",
		in:         "[1234::5678%lo]:12345",
		wantHost:   "1234::5678%lo",
		wantErrMsg: "",
		wantPort:   12345,
	}, {
		name:       "success_host",
		in:         "example.com:12345",
		wantHost:   "example.com",
		wantErrMsg: "",
		wantPort:   12345,
	}, {
		name:       "bad_port",
		in:         "example.com:!!!",
		wantHost:   "",
		wantErrMsg: "parsing port: strconv.Atoi: parsing \"!!!\": invalid syntax",
		wantPort:   0,
	}, {
		name:       "bad_syntax",
		in:         "[1234::5678:12345",
		wantHost:   "",
		wantErrMsg: "address [1234::5678:12345: missing ']' in address",
		wantPort:   0,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			host, port, err := netutil.SplitHostPort(tc.in)
			assert.Equal(t, tc.wantHost, host)
			assert.Equal(t, tc.wantPort, port)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
		})
	}
}

func repeatStr(b *strings.Builder, s string, n int) {
	for i := 0; i < n; i++ {
		stringutil.WriteToBuilder(b, s)
	}
}

func TestValidateDomainName(t *testing.T) {
	t.Parallel()

	b := &strings.Builder{}
	repeatStr(b, "a", 255)
	longDomainName := b.String()

	b.Reset()
	repeatStr(b, "a", 64)
	longLabel := b.String()

	_, _ = b.WriteString(".com")
	longLabelDomainName := b.String()

	testCases := []struct {
		name       string
		in         string
		wantErrAs  interface{}
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
