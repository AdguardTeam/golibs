package netutil_test

import (
	"net"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestIsValidMACString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		want assert.BoolAssertionFunc
		name string
		in   string
	}{{
		want: assert.True,
		name: "good_eui_48",
		in:   "00:01:02:03:04:05",
	}, {
		want: assert.True,
		name: "good_eui_48_dot",
		in:   "0000.5e00.5301",
	}, {
		want: assert.True,
		name: "good_eui_48_hyphen",
		in:   "00-00-5e-00-53-01",
	}, {
		want: assert.True,
		name: "good_eui_48_no_sep",
		in:   "00005e005301",
	}, {
		want: assert.True,
		name: "good_eui_64",
		in:   "00:01:02:03:04:05:06:07",
	}, {
		want: assert.True,
		name: "good_eui_64_no_sep",
		in:   "02005e1000000001",
	}, {
		want: assert.True,
		name: "good_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13",
	}, {
		want: assert.True,
		name: "good_infiniband_no_sep",
		in:   "000102030405060708090a0b0c0d0e0f10111213",
	}, {
		want: assert.False,
		name: "bad_empty",
		in:   "",
	}, {
		want: assert.False,
		name: "bad_short",
		in:   "00:01:02:03",
	}, {
		want: assert.False,
		name: "bad_long",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13:14",
	}, {
		want: assert.False,
		name: "bad_long_no_sep",
		in:   "02005e10000000011",
	}, {
		want: assert.False,
		name: "bad_eui_48",
		in:   "00:01:02:03:04:!!",
	}, {
		want: assert.False,
		name: "bad_eui_48_dot",
		in:   "0000.5e00.!!!!",
	}, {
		want: assert.False,
		name: "bad_eui_48_no_sep",
		in:   "0001020304!!",
	}, {
		want: assert.False,
		name: "bad_eui_64",
		in:   "00:01:02:03:04:05:06:!!",
	}, {
		want: assert.False,
		name: "bad_eui_64_no_sep",
		in:   "00010203040506!!",
	}, {
		want: assert.False,
		name: "bad_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:!!",
	}, {
		want: assert.False,
		name: "bad_infiniband_no_sep",
		in:   "000102030405060708090a0b0c0d0e0f101112!!",
	}, {
		want: assert.False,
		name: "bad_eui_48_dot_sep",
		in:   "00.00.00.00.00.00",
	}, {
		want: assert.False,
		name: "bad_eui_64_dot_sep",
		in:   "00.00.00.00.00.00.00.00",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.want(t, netutil.IsValidMACString(tc.in))
		})
	}
}

func BenchmarkIsValidMACString(b *testing.B) {
	benchCases := []struct {
		want require.BoolAssertionFunc
		name string
		in   string
	}{{
		want: require.True,
		name: "good_eui_48",
		in:   "00:01:02:03:04:05",
	}, {
		want: require.True,
		name: "good_eui_48_dot",
		in:   "0000.5e00.5301",
	}, {
		want: require.True,
		name: "good_eui_48_hyphen",
		in:   "00-00-5e-00-53-01",
	}, {
		want: require.True,
		name: "good_eui_48_no_sep",
		in:   "00005e005301",
	}, {
		want: require.True,
		name: "good_eui_64",
		in:   "00:01:02:03:04:05:06:07",
	}, {
		want: require.True,
		name: "good_eui_64_no_sep",
		in:   "02005e1000000001",
	}, {
		want: require.True,
		name: "good_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13",
	}, {
		want: require.True,
		name: "good_infiniband_no_sep",
		in:   "000102030405060708090a0b0c0d0e0f10111213",
	}, {
		want: require.False,
		name: "bad_empty",
		in:   "",
	}, {
		want: require.False,
		name: "bad_short",
		in:   "00:01:02:03",
	}, {
		want: require.False,
		name: "bad_long",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13:14",
	}, {
		want: require.False,
		name: "bad_eui_48",
		in:   "00:01:02:03:04:!!",
	}, {
		want: require.False,
		name: "bad_eui_48_no_sep",
		in:   "0001020304!!",
	}, {
		want: require.False,
		name: "bad_eui_64",
		in:   "00:01:02:03:04:05:06:!!",
	}, {
		want: require.False,
		name: "bad_eui_64_no_sep",
		in:   "00010203040506!!",
	}, {
		want: require.False,
		name: "bad_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:!!",
	}, {
		want: require.False,
		name: "bad_infiniband_no_sep",
		in:   "000102030405060708090a0b0c0d0e0f101112!!",
	}}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			var got bool
			b.ReportAllocs()
			for b.Loop() {
				got = netutil.IsValidMACString(bc.in)
			}

			bc.want(b, got)
		})
	}

	// Most recent results:
	//  goos: linux
	//  goarch: amd64
	//  pkg: github.com/AdguardTeam/golibs/netutil
	//  cpu: AMD Ryzen AI 7 PRO 350 w/ Radeon 860M
	// BenchmarkIsValidMACString/good_eui_48-16                97961534                12.65 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_eui_48_dot-16            100000000               10.32 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_eui_48_hyphen-16         83985573                12.26 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_eui_48_no_sep-16         145311968                8.333 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_eui_64-16                85618686                14.83 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_eui_64_no_sep-16         100000000               10.65 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_infiniband-16            35572928                33.31 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/good_infiniband_no_sep-16             59442024                20.19 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_empty-16                          660256630                1.864 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_short-16                          676261029                1.855 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_long-16                           529534195                2.273 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_eui_48-16                         112847419               10.63 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_eui_48_no_sep-16                  183439570                6.471 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_eui_64-16                         89685344                13.25 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_eui_64_no_sep-16                  160003975                7.552 ns/op           0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_infiniband-16                     37117690                29.74 ns/op            0 B/op          0 allocs/op
	// BenchmarkIsValidMACString/bad_infiniband_no_sep-16              69095757                14.61 ns/op            0 B/op          0 allocs/op
}

func FuzzIsValidMACString(f *testing.F) {
	for _, seed := range []string{
		"",
		" ",
		"00",
		"00:00",
		"00.00",
		"00-00",
		"00:db::68",
		"!!:00:00:00:00:00",
		"00005e005301",
		"02005e1000000001",
		"00.00.00.00.00.00",
		"00.00.00.00.00.00.00.00",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ok := netutil.IsValidMACString(input)
		_, err := net.ParseMAC(input)

		require.Equalf(t, err == nil, ok, "input: %q, error: %v", input, err)
	})
}
