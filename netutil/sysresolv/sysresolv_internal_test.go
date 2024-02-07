package sysresolv

import (
	"fmt"
	"net/netip"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemResolvers_Refresh(t *testing.T) {
	t.Parallel()

	t.Run("expected_error", func(t *testing.T) {
		t.Parallel()

		sr, err := NewSystemResolvers(defaultHostGenFunc, 53)
		require.NoError(t, err)

		assert.NoError(t, sr.Refresh())
	})

	t.Run("unexpected_error", func(t *testing.T) {
		t.Parallel()

		_, err := NewSystemResolvers(func() (host string) {
			return "127.0.0.1::123"
		}, 53)

		assert.Error(t, err)
	})
}

func TestSystemResolvers_Parse(t *testing.T) {
	t.Parallel()

	sr, err := NewSystemResolvers(defaultHostGenFunc, 53)
	require.NoError(t, err)

	testCases := []struct {
		want    netip.AddrPort
		wantErr error
		name    string
		address string
	}{{
		want:    netip.MustParseAddrPort("127.0.0.1:53"),
		wantErr: nil,
		name:    "valid_ipv4",
		address: "127.0.0.1",
	}, {
		want:    netip.MustParseAddrPort("[::1]:53"),
		wantErr: nil,
		name:    "valid_ipv6_port",
		address: "[::1]:53",
	}, {
		want:    netip.MustParseAddrPort("[::1%lo0]:53"),
		wantErr: nil,
		name:    "valid_ipv6_zone_port",
		address: "[::1%lo0]:53",
	}, {
		want:    netip.AddrPort{},
		wantErr: errBadAddrPassed,
		name:    "invalid_split_host",
		address: "127.0.0.1::123",
	}, {
		want:    netip.AddrPort{},
		wantErr: errBadAddrPassed,
		name:    "invalid_ipv6_zone_port",
		address: "[:::1%lo0]:53",
	}, {
		want:    netip.AddrPort{},
		wantErr: errBadAddrPassed,
		name:    "invalid_parse_ip",
		address: "not-ip",
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, parseErr := sr.parse(tc.address)
			assert.Equal(t, tc.want, got)

			assert.ErrorIs(t, parseErr, tc.wantErr)
		})
	}
}

func TestCompareAddrPorts(t *testing.T) {
	t.Parallel()

	var (
		addr1v4 = netip.MustParseAddrPort("1.2.3.4:1")
		addr2v4 = netip.MustParseAddrPort("4.3.2.1:1")
		addr3v4 = netip.MustParseAddrPort("1.2.3.4:2")

		addr1v6 = netip.MustParseAddrPort("[::1]:1")
		addr2v6 = netip.MustParseAddrPort("[::2]:1")
		addr3v6 = netip.MustParseAddrPort("[::1]:2")
	)

	testCases := []struct {
		name  string
		addrs []netip.AddrPort
		want  []netip.AddrPort
	}{{
		name:  "ipv4",
		addrs: []netip.AddrPort{addr3v4, addr2v4, addr1v4},
		want:  []netip.AddrPort{addr1v4, addr3v4, addr2v4},
	}, {
		name:  "ipv6",
		addrs: []netip.AddrPort{addr3v6, addr2v6, addr1v6},
		want:  []netip.AddrPort{addr1v6, addr3v6, addr2v6},
	}, {
		name:  "mixed",
		addrs: []netip.AddrPort{addr3v4, addr3v6, addr2v4, addr2v6, addr1v4, addr1v6},
		want:  []netip.AddrPort{addr1v4, addr3v4, addr2v4, addr1v6, addr3v6, addr2v6},
	}}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			slices.SortFunc(tc.addrs, compareAddrPorts)
			assert.Equal(t, tc.want, tc.addrs)
		})
	}
}

// strSink is a typed sink for sinking string values from benchmarks.
var strSink string

func BenchmarkHostGenFunc_default(b *testing.B) {
	b.Run("builder", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			strSink = defaultHostGenFunc()
		}
	})

	b.Run("sprintf", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			strSink = hostGenFuncSprintf()
		}
	})

	// The most recent results:
	//
	// goos: darwin
	// goarch: amd64
	// pkg: github.com/AdguardTeam/golibs/netutil/sysresolv
	// cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
	// BenchmarkHostGenFunc_default/builder-12	7903302		136.4 ns/op		32 B/op		1 allocs/op
	// BenchmarkHostGenFunc_default/sprintf-12	5793406		205.5 ns/op		40 B/op		2 allocs/op
}

// hostGenFuncSprintf is a [HostGenFunc] implementation that uses a plain
// [fmt.Sprintf] call.  It exists for benchmarking purposes only.
func hostGenFuncSprintf() (hostname string) {
	return fmt.Sprintf("host%d.test", time.Now().UnixNano())
}
