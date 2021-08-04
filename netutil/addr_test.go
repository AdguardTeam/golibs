package netutil_test

import (
	"net"
	"net/url"
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/stringutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneIP(t *testing.T) {
	assert.Equal(t, net.IP(nil), netutil.CloneIP(nil))
	assert.Equal(t, net.IP{}, netutil.CloneIP(net.IP{}))

	ip := net.IP{1, 2, 3, 4}
	clone := netutil.CloneIP(ip)
	assert.Equal(t, ip, clone)

	require.Len(t, clone, len(ip))

	assert.NotSame(t, &ip[0], &clone[0])
}

func TestCloneIPs(t *testing.T) {
	assert.Equal(t, []net.IP(nil), netutil.CloneIPs(nil))
	assert.Equal(t, []net.IP{}, netutil.CloneIPs([]net.IP{}))

	ips := []net.IP{{1, 2, 3, 4}}
	clone := netutil.CloneIPs(ips)
	assert.Equal(t, ips, clone)

	require.Len(t, clone, len(ips))
	require.Len(t, clone[0], len(ips[0]))

	assert.NotSame(t, &ips[0], &clone[0])
	assert.NotSame(t, &ips[0][0], &clone[0][0])
}

func TestCloneMAC(t *testing.T) {
	assert.Equal(t, net.HardwareAddr(nil), netutil.CloneMAC(nil))
	assert.Equal(t, net.HardwareAddr{}, netutil.CloneMAC(net.HardwareAddr{}))

	mac := net.HardwareAddr{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC}
	clone := netutil.CloneMAC(mac)
	assert.Equal(t, mac, clone)

	require.Len(t, clone, len(mac))

	assert.NotSame(t, &mac[0], &clone[0])
}

func TestCloneURL(t *testing.T) {
	assert.Equal(t, (*url.URL)(nil), netutil.CloneURL(nil))
	assert.Equal(t, &url.URL{}, netutil.CloneURL(&url.URL{}))

	u, err := url.Parse("https://example.com/path?q=1&q=2#frag")
	require.NoError(t, err)

	clone := netutil.CloneURL(u)
	assert.Equal(t, u, clone)
	assert.NotSame(t, u, clone)
}

func TestIPPortFromAddr(t *testing.T) {
	ip := net.IP{1, 2, 3, 4}

	testCases := []struct {
		name     string
		in       net.Addr
		wantIP   net.IP
		wantPort int
	}{{
		name:     "nil",
		in:       nil,
		wantIP:   nil,
		wantPort: 0,
	}, {
		name:     "tcp",
		in:       &net.TCPAddr{IP: ip, Port: 12345},
		wantIP:   ip,
		wantPort: 12345,
	}, {
		name:     "udp",
		in:       &net.UDPAddr{IP: ip, Port: 12345},
		wantIP:   ip,
		wantPort: 12345,
	}, {
		name:     "custom",
		in:       struct{ net.Addr }{},
		wantIP:   nil,
		wantPort: 0,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotIP, gotPort := netutil.IPPortFromAddr(tc.in)
			assert.Equal(t, tc.wantIP, gotIP)
			assert.Equal(t, tc.wantPort, gotPort)
		})
	}
}

func TestValidateMAC(t *testing.T) {
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
		wantErrMsg: `bad mac address "": mac address is empty`,
		wantErrAs:  new(*netutil.EmptyError),
		in:         nil,
	}, {
		name:       "error_empty",
		wantErrMsg: `bad mac address "": mac address is empty`,
		wantErrAs:  new(*netutil.EmptyError),
		in:         net.HardwareAddr{},
	}, {
		name: "error_bad",
		wantErrMsg: `bad mac address "00:01:02:03": ` +
			`bad mac address length 4, allowed: [6 8 20]`,
		wantErrAs: new(*netutil.BadLengthError),
		in:        net.HardwareAddr{0x00, 0x01, 0x02, 0x03},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := netutil.ValidateMAC(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestJoinHostPort(t *testing.T) {
	assert.Equal(t, ":0", netutil.JoinHostPort("", 0))
	assert.Equal(t, "host:12345", netutil.JoinHostPort("host", 12345))
	assert.Equal(t, "1.2.3.4:12345", netutil.JoinHostPort("1.2.3.4", 12345))
	assert.Equal(t, "[1234::5678]:12345", netutil.JoinHostPort("1234::5678", 12345))
	assert.Equal(t, "[1234::5678%lo]:12345", netutil.JoinHostPort("1234::5678%lo", 12345))
}

func TestSplitHostPort(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			host, port, err := netutil.SplitHostPort(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantHost, host)
				assert.Equal(t, tc.wantPort, port)
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
			}
		})
	}
}

func repeatStr(b *strings.Builder, s string, n int) {
	for i := 0; i < n; i++ {
		stringutil.WriteToBuilder(b, s)
	}
}

func TestValidateDomainName(t *testing.T) {
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
		wantErrAs:  new(*netutil.EmptyError),
		wantErrMsg: `bad domain name "": domain name is empty`,
	}, {
		name:       "bad_symbol",
		in:         "!!!",
		wantErrAs:  new(*netutil.BadRuneError),
		wantErrMsg: `bad domain name "!!!": bad domain name label "!!!": bad domain name label rune '!'`,
	}, {
		name:       "bad_length",
		in:         longDomainName,
		wantErrAs:  new(*netutil.TooLongError),
		wantErrMsg: `bad domain name "` + longDomainName + `": domain name is too long, max: 253`,
	}, {
		name:      "bad_label_length",
		in:        longLabelDomainName,
		wantErrAs: new(*netutil.TooLongError),
		wantErrMsg: `bad domain name "` + longLabelDomainName + `": ` +
			`bad domain name label "` + longLabel + `": ` +
			`domain name label is too long, max: 63`,
	}, {
		name:      "bad_label_empty",
		in:        "example..com",
		wantErrAs: new(*netutil.EmptyError),
		wantErrMsg: `bad domain name "example..com": ` +
			`bad domain name label "": domain name label is empty`,
	}, {
		name:      "bad_label_first_symbol",
		in:        "example.-aa.com",
		wantErrAs: new(*netutil.BadRuneError),
		wantErrMsg: `bad domain name "example.-aa.com": ` +
			`bad domain name label "-aa": bad domain name label rune '-'`,
	}, {
		name:      "bad_label_last_symbol",
		in:        "example-.aa.com",
		wantErrAs: new(*netutil.BadRuneError),
		wantErrMsg: `bad domain name "example-.aa.com": ` +
			`bad domain name label "example-": bad domain name label rune '-'`,
	}, {
		name:      "bad_label_symbol",
		in:        "example.a!!!.com",
		wantErrAs: new(*netutil.BadRuneError),
		wantErrMsg: `bad domain name "example.a!!!.com": ` +
			`bad domain name label "a!!!": bad domain name label rune '!'`,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := netutil.ValidateDomainName(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.ErrorAs(t, err, new(*netutil.BadDomainError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}

func TestUnreverseAddr(t *testing.T) {
	const (
		ipv4Good          = `1.0.0.127.in-addr.arpa`
		ipv4GoodUppercase = `1.0.0.127.In-Addr.Arpa`

		ipv4Missing = `.0.0.127.in-addr.arpa`
		ipv4Char    = `1.0.z.127.in-addr.arpa`
	)

	const (
		ipv6Zeroes = `0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0`
		ipv6Suffix = ipv6Zeroes + `.ip6.arpa`

		ipv6Good          = `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
		ipv6GoodUppercase = `4.3.2.1.D.C.B.A.0.0.0.0.0.0.0.0.` + ipv6Suffix

		ipv6CharHi  = `4.3.2.1.d.c.b.a.0.z.0.0.0.0.0.0.` + ipv6Suffix
		ipv6CharLo  = `4.3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6Suffix
		ipv6Dots    = `4.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6Zeroes + `..ip6.arpa`
		ipv6Len     = `3.2.1.d.c.b.a.z.0.0.0.0.0.0.0.` + ipv6Suffix
		ipv6Many    = `4.3.2.1.dbc.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
		ipv6Missing = `.3.2.1.d.c.b.a.0.0.0.0.0.0.0.0.` + ipv6Suffix
		ipv6Space   = `4.3.2.1.d.c.b.a. .0.0.0.0.0.0.0.` + ipv6Suffix
	)

	testCases := []struct {
		name       string
		in         string
		wantErrMsg string
		wantErrAs  interface{}
		want       net.IP
	}{{
		name:       "good_ipv4",
		in:         ipv4Good,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.IP{127, 0, 0, 1},
	}, {
		name:       "good_ipv4_fqdn",
		in:         ipv4Good + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.IP{127, 0, 0, 1},
	}, {
		name:       "good_ipv4_case",
		in:         ipv4GoodUppercase,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.IP{127, 0, 0, 1},
	}, {
		name: "bad_ipv4_missing",
		in:   ipv4Missing,
		wantErrMsg: `bad arpa domain name "` + ipv4Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.EmptyError),
		want:      nil,
	}, {
		name: "bad_ipv4_char",
		in:   ipv4Char,
		wantErrMsg: `bad arpa domain name "` + ipv4Char + `": ` +
			`bad ip address "1.0.z.127"`,
		wantErrAs: new(*netutil.BadIPError),
		want:      nil,
	}, {
		name:       "good_ipv6",
		in:         ipv6Good,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.ParseIP("::abcd:1234"),
	}, {
		name:       "good_ipv6_fqdn",
		in:         ipv6Good + ".",
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.ParseIP("::abcd:1234"),
	}, {
		name:       "good_ipv6_case",
		in:         ipv6GoodUppercase,
		wantErrMsg: "",
		wantErrAs:  nil,
		want:       net.ParseIP("::abcd:1234"),
	}, {
		name: "bad_ipv6_many",
		in:   ipv6Many,
		wantErrMsg: `bad arpa domain name "` + ipv6Many + `": ` +
			`not a full reversed ip address`,
		wantErrAs: new(*netutil.BadDomainError),
		want:      nil,
	}, {
		name: "bad_ipv6_missing",
		in:   ipv6Missing,
		wantErrMsg: `bad arpa domain name "` + ipv6Missing + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.EmptyError),
		want:      nil,
	}, {
		name: "bad_ipv6_char_lo",
		in:   ipv6CharLo,
		wantErrMsg: `bad arpa domain name "` + ipv6CharLo + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.BadRuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_char_hi",
		in:   ipv6CharHi,
		wantErrMsg: `bad arpa domain name "` + ipv6CharHi + `": ` +
			`bad arpa domain name rune 'z'`,
		wantErrAs: new(*netutil.BadRuneError),
		want:      nil,
	}, {
		name: "bad_ipv6_dots",
		in:   ipv6Dots,
		wantErrMsg: `bad arpa domain name "` + ipv6Dots + `": ` +
			`bad domain name label "": domain name label is empty`,
		wantErrAs: new(*netutil.EmptyError),
		want:      nil,
	}, {
		name: "bad_ipv6_len",
		in:   ipv6Len,
		wantErrMsg: `bad arpa domain name "` + ipv6Len + `": ` +
			`bad arpa domain name length 70, allowed: [72]`,
		wantErrAs: new(*netutil.BadLengthError),
		want:      nil,
	}, {
		name: "bad_ipv6_space",
		in:   ipv6Space,
		wantErrMsg: `bad arpa domain name "` + ipv6Space + `": ` +
			`bad domain name label " ": bad domain name label rune ' '`,
		wantErrAs: new(*netutil.BadRuneError),
		want:      nil,
	}, {
		name:       "not_a_reversed_ip",
		in:         "1.2.3.4",
		wantErrMsg: `bad arpa domain name "1.2.3.4": not a full reversed ip address`,
		wantErrAs:  new(errors.Error),
		want:       nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, err := netutil.IPFromReversedAddr(tc.in)
			if tc.wantErrMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.To16(), ip.To16())
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.ErrorAs(t, err, new(*netutil.BadDomainError))
				assert.ErrorAs(t, err, tc.wantErrAs)
			}
		})
	}
}
