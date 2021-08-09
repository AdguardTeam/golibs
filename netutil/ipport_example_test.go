package netutil_test

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleIPPort_MarshalText() {
	ip4 := net.ParseIP("1.2.3.4")
	ip6 := net.ParseIP("1234::cdef")

	resp := struct {
		IPs []netutil.IPPort `json:"ips"`
	}{
		IPs: []netutil.IPPort{{
			IP:   ip4,
			Port: 12345,
		}, {
			IP:   ip6,
			Port: 23456,
		}},
	}

	err := json.NewEncoder(os.Stdout).Encode(resp)
	if err != nil {
		panic(err)
	}

	respPtrs := struct {
		IPPtrs []*netutil.IPPort `json:"ip_ptrs"`
	}{
		IPPtrs: []*netutil.IPPort{{
			IP:   ip4,
			Port: 12345,
		}, {
			IP:   ip6,
			Port: 23456,
		}},
	}

	err = json.NewEncoder(os.Stdout).Encode(respPtrs)
	if err != nil {
		panic(err)
	}

	// Output:
	//
	// {"ips":["1.2.3.4:12345","[1234::cdef]:23456"]}
	// {"ip_ptrs":["1.2.3.4:12345","[1234::cdef]:23456"]}
}

func ExampleIPPort_String() {
	ip4 := net.ParseIP("1.2.3.4")
	ip6 := net.ParseIP("1234::cdef")

	ipp := &netutil.IPPort{
		IP:   ip4,
		Port: 12345,
	}

	fmt.Println(ipp)

	ipp.IP = ip6
	fmt.Println(ipp)

	ipp.Port = 0
	fmt.Println(ipp)

	ipp.IP = nil
	fmt.Println(ipp)

	// Output:
	//
	// 1.2.3.4:12345
	// [1234::cdef]:12345
	// [1234::cdef]:0
	// :0
}

func ExampleIPPort_TCP() {
	ipp := &netutil.IPPort{
		IP:   net.IP{1, 2, 3, 4},
		Port: 12345,
	}

	fmt.Printf("%#v\n", ipp.TCP())

	// Output:
	//
	// &net.TCPAddr{IP:net.IP{0x1, 0x2, 0x3, 0x4}, Port:12345, Zone:""}
}

func ExampleIPPort_UDP() {
	ipp := &netutil.IPPort{
		IP:   net.IP{1, 2, 3, 4},
		Port: 12345,
	}

	fmt.Printf("%#v\n", ipp.UDP())

	// Output:
	//
	// &net.UDPAddr{IP:net.IP{0x1, 0x2, 0x3, 0x4}, Port:12345, Zone:""}
}

func ExampleIPPort_UnmarshalText() {
	resp := &struct {
		IPs []netutil.IPPort `json:"ips"`
	}{}

	r := strings.NewReader(`{"ips":["1.2.3.4:12345","[1234::cdef]:23456"]}`)
	err := json.NewDecoder(r).Decode(resp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", resp.IPs)

	respPtrs := &struct {
		IPPtrs []*netutil.IPPort `json:"ip_ptrs"`
	}{}

	r = strings.NewReader(`{"ip_ptrs":["1.2.3.4:12345","[1234::cdef]:23456"]}`)
	err = json.NewDecoder(r).Decode(respPtrs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", respPtrs.IPPtrs)

	// Output:
	//
	// [1.2.3.4:12345 [1234::cdef]:23456]
	// [1.2.3.4:12345 [1234::cdef]:23456]
}
