package netutil_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleHostPort_MarshalText() {
	resp := struct {
		Hosts []netutil.HostPort `json:"hosts"`
	}{
		Hosts: []netutil.HostPort{{
			Host: "example.com",
			Port: 12345,
		}, {
			Host: "example.org",
			Port: 23456,
		}},
	}

	err := json.NewEncoder(os.Stdout).Encode(resp)
	if err != nil {
		panic(err)
	}

	respPtrs := struct {
		HostPtrs []*netutil.HostPort `json:"host_ptrs"`
	}{
		HostPtrs: []*netutil.HostPort{{
			Host: "example.com",
			Port: 12345,
		}, {
			Host: "example.org",
			Port: 23456,
		}},
	}

	err = json.NewEncoder(os.Stdout).Encode(respPtrs)
	if err != nil {
		panic(err)
	}

	// Output:
	//
	// {"hosts":["example.com:12345","example.org:23456"]}
	// {"host_ptrs":["example.com:12345","example.org:23456"]}
}

func ExampleHostPort_String() {
	hp := &netutil.HostPort{
		Host: "example.com",
		Port: 12345,
	}

	fmt.Println(hp)

	hp.Host = "1234::cdef"
	fmt.Println(hp)

	hp.Port = 0
	fmt.Println(hp)

	hp.Host = ""
	fmt.Println(hp)

	// Output:
	//
	// example.com:12345
	// [1234::cdef]:12345
	// [1234::cdef]:0
	// :0
}

func ExampleHostPort_UnmarshalText() {
	resp := &struct {
		Hosts []netutil.HostPort `json:"hosts"`
	}{}

	r := strings.NewReader(`{"hosts":["example.com:12345","example.org:23456"]}`)
	err := json.NewDecoder(r).Decode(resp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", resp.Hosts[0])
	fmt.Printf("%#v\n", resp.Hosts[1])

	respPtrs := &struct {
		HostPtrs []*netutil.HostPort `json:"host_ptrs"`
	}{}

	r = strings.NewReader(`{"host_ptrs":["example.com:12345","example.org:23456"]}`)
	err = json.NewDecoder(r).Decode(respPtrs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", respPtrs.HostPtrs[0])
	fmt.Printf("%#v\n", respPtrs.HostPtrs[1])

	// Output:
	//
	// netutil.HostPort{Host:"example.com", Port:12345}
	// netutil.HostPort{Host:"example.org", Port:23456}
	// &netutil.HostPort{Host:"example.com", Port:12345}
	// &netutil.HostPort{Host:"example.org", Port:23456}
}
