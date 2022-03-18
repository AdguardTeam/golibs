package netutil_test

import (
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExampleSubnetSet_single() {
	var s netutil.SubnetSet
	_, s, err := net.ParseCIDR("1.2.3.4/16")
	if err != nil {
		panic(err)
	}

	fmt.Println("contains 1.2.3.4:", s.Contains(net.IP{1, 2, 3, 4}))
	fmt.Println("contains 4.3.2.1:", s.Contains(net.IP{4, 3, 2, 1}))

	// Output:
	//
	// contains 1.2.3.4: true
	// contains 4.3.2.1: false
}

func ExampleSubnetSet_slice() {
	nets, err := netutil.ParseSubnets("1.2.3.0/24", "ffff:12ab::/16")
	if err != nil {
		panic(err)
	}

	s := netutil.SliceSubnetSet(nets)

	fmt.Println("contains 1.2.3.4:      ", s.Contains(net.IP{1, 2, 3, 4}))
	fmt.Println("contains 4.3.2.1:      ", s.Contains(net.IP{4, 3, 2, 1}))
	fmt.Println("contains ffff:12ab::10:", s.Contains(net.ParseIP("ffff:12ab::10")))
	fmt.Println("contains 12ab:ffff::10:", s.Contains(net.ParseIP("12ab:ffff::10")))

	fmt.Println()

	s = netutil.SliceSubnetSet{}
	fmt.Println("contains 1.2.3.4:      ", s.Contains(net.IP{1, 2, 3, 4}))
	fmt.Println("contains ffff:12ab::10:", s.Contains(net.ParseIP("ffff:12ab::10")))

	fmt.Println()

	s = netutil.SliceSubnetSet{{
		IP:   make(net.IP, net.IPv4len),
		Mask: make(net.IPMask, net.IPv4len),
	}, {
		IP:   make(net.IP, net.IPv6len),
		Mask: make(net.IPMask, net.IPv6len),
	}}
	fmt.Println("contains 1.2.3.4:      ", s.Contains(net.IP{1, 2, 3, 4}))
	fmt.Println("contains ffff:12ab::10:", s.Contains(net.ParseIP("ffff:12ab::10")))
	fmt.Println("contains <nil>:        ", s.Contains(nil))

	// Output:
	//
	// contains 1.2.3.4:       true
	// contains 4.3.2.1:       false
	// contains ffff:12ab::10: true
	// contains 12ab:ffff::10: false
	//
	// contains 1.2.3.4:       false
	// contains ffff:12ab::10: false
	//
	// contains 1.2.3.4:       true
	// contains ffff:12ab::10: true
	// contains <nil>:         false
}

func ExampleSubnetSet_func() {
	s := netutil.SubnetSetFunc(func(ip net.IP) (ok bool) {
		return len(ip) > 0 && ip[0] == 0xFF
	})

	fmt.Println("contains 255.0.0.1:", s.Contains(net.IP{255, 0, 0, 1}))
	fmt.Println("contains 254.0.0.1:", s.Contains(net.IP{254, 0, 0, 1}))
	fmt.Println("contains ff00:::1: ", s.Contains(net.ParseIP("ff00::1")))
	fmt.Println("contains ff:::1:   ", s.Contains(net.ParseIP("ff::1")))

	// Output:
	//
	// contains 255.0.0.1: true
	// contains 254.0.0.1: false
	// contains ff00:::1:  true
	// contains ff:::1:    false
}
