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
		name: "good_eui_64",
		in:   "00:01:02:03:04:05:06:07",
	}, {
		want: assert.True,
		name: "good_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13",
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
		name: "bad_eui_48",
		in:   "00:01:02:03:04:!!",
	}, {
		want: assert.False,
		name: "bad_eui_48_dot",
		in:   "0000.5e00.!!!!",
	}, {
		want: assert.False,
		name: "bad_eui_64",
		in:   "00:01:02:03:04:05:06:!!",
	}, {
		want: assert.False,
		name: "bad_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:!!",
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
		name: "good_eui_64",
		in:   "00:01:02:03:04:05:06:07",
	}, {
		want: require.True,
		name: "good_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:13",
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
		name: "bad_eui_64",
		in:   "00:01:02:03:04:05:06:!!",
	}, {
		want: require.False,
		name: "bad_infiniband",
		in:   "00:01:02:03:04:05:06:07:08:09:0a:0b:0c:0d:0e:0f:10:11:12:!!",
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
	//	goos: darwin
	//	goarch: arm64
	//	pkg: github.com/AdguardTeam/golibs/netutil
	//	cpu: Apple M1 Pro
	//	BenchmarkIsValidMACString/good_eui_48-8         	64798680	        18.55 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/good_eui_48_dot-8     	67330148	        18.06 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/good_eui_48_hyphen-8  	65400295	        18.23 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/good_eui_64-8         	52725235	        22.59 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/good_infiniband-8     	24440022	        49.04 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_empty-8           	587863916	         2.044 ns/op       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_short-8           	583646232	         2.041 ns/op       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_long-8            	494821214	         2.413 ns/op       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_eui_48-8          	67780974	        17.69 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_eui_64-8          	54463244	        22.02 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkIsValidMACString/bad_infiniband-8      	24736444	        48.30 ns/op	       0 B/op	       0 allocs/op
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
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ok := netutil.IsValidMACString(input)
		_, err := net.ParseMAC(input)

		require.Equalf(t, err == nil, ok, "input: %q, error: %v", input, err)
	})
}
