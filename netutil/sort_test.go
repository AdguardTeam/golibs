package netutil_test

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/AdguardTeam/golibs/netutil"
	"github.com/stretchr/testify/assert"
)

func TestPreferIPv4(t *testing.T) {
	testCases := []struct {
		name  string
		addrs []netip.Addr
		want  []netip.Addr
	}{{
		name: "v4_preferred",
		addrs: []netip.Addr{
			testIPv6Addr.Next(),
			testIPv6Addr,
			{},
			testIPv4Addr.Next(),
			testIPv4Addr,
		},
		want: []netip.Addr{
			testIPv4Addr,
			testIPv4Addr.Next(),
			testIPv6Addr,
			testIPv6Addr.Next(),
			{},
		},
	}, {
		name: "shuffled_v4_preferred",
		addrs: []netip.Addr{
			{},
			testIPv4Addr,
			testIPv6Addr.Next(),
			testIPv6Addr,
			testIPv4Addr.Next(),
		},
		want: []netip.Addr{
			testIPv4Addr,
			testIPv4Addr.Next(),
			testIPv6Addr,
			testIPv6Addr.Next(),
			{},
		},
	}, {
		name:  "empty",
		addrs: []netip.Addr{},
		want:  []netip.Addr{},
	}, {
		name:  "single",
		addrs: []netip.Addr{testIPv4Addr},
		want:  []netip.Addr{testIPv4Addr},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ips := slices.Clone(tc.addrs)
			slices.SortFunc(ips, netutil.PreferIPv4)
			assert.Equal(t, tc.want, ips)
		})
	}
}

func TestPreferIPv6(t *testing.T) {
	testCases := []struct {
		name  string
		addrs []netip.Addr
		want  []netip.Addr
	}{{
		name: "v6_preferred",
		addrs: []netip.Addr{
			testIPv4Addr,
			testIPv4Addr.Next(),
			{},
			testIPv6Addr,
			testIPv6Addr.Next(),
		},
		want: []netip.Addr{
			testIPv6Addr,
			testIPv6Addr.Next(),
			testIPv4Addr,
			testIPv4Addr.Next(),
			{},
		},
	}, {
		name: "shuffled_v6_preferred",
		addrs: []netip.Addr{
			{},
			testIPv4Addr,
			testIPv6Addr.Next(),
			testIPv6Addr,
			testIPv4Addr.Next(),
		},
		want: []netip.Addr{
			testIPv6Addr,
			testIPv6Addr.Next(),
			testIPv4Addr,
			testIPv4Addr.Next(),
			{},
		},
	}, {
		name: "start_with_ipv4",
		addrs: []netip.Addr{
			testIPv4Addr,
			testIPv6Addr,
			testIPv4Addr.Next(),
			testIPv6Addr.Next(),
		},
		want: []netip.Addr{
			testIPv6Addr,
			testIPv6Addr.Next(),
			testIPv4Addr,
			testIPv4Addr.Next(),
		},
	}, {
		name: "start_with_ipv6",
		addrs: []netip.Addr{
			testIPv6Addr,
			testIPv4Addr,
			testIPv6Addr.Next(),
			testIPv4Addr.Next(),
		},
		want: []netip.Addr{
			testIPv6Addr,
			testIPv6Addr.Next(),
			testIPv4Addr,
			testIPv4Addr.Next(),
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ips := slices.Clone(tc.addrs)
			slices.SortFunc(ips, netutil.PreferIPv6)
			assert.Equal(t, tc.want, ips)
		})
	}
}
