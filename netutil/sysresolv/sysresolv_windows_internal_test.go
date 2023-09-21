//go:build windows

package sysresolv

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanAddrs(t *testing.T) {
	testCases := []struct {
		name        string
		nslookupOut string
		want        []string
		wantErr     bool
	}{{
		name:        "simple",
		nslookupOut: "Default Server:  dns.google\nAddress:  8.8.8.8\n",
		want:        []string{"8.8.8.8"},
	}, {
		name:        "simplev6",
		nslookupOut: "Server:  UnKnown\nAddress:  fec0:0:0::1",
		want:        []string{"fec0:0:0::1"},
	}, {
		name:        "unavailable",
		nslookupOut: "***Default servers are unavailable\nServer:  UnKnown\nAddress:  127.0.0.1",
		want:        []string{"127.0.0.1"},
	}, {
		name:        "invalidip",
		nslookupOut: "Server:  UnKnown\nAddress:  172.168.1.3.local.net\n",
		wantErr:     true,
	}, {
		name:        "garbage",
		nslookupOut: "Server:  UnKnown\nRandom stuff",
		wantErr:     true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewBufferString(tc.nslookupOut)
			s := bufio.NewScanner(r)
			addrs, err := scanAddrs(s)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, addrs)
			}
		})
	}
}
