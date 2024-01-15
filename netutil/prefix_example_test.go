package netutil_test

import (
	"fmt"
	"net/netip"

	"github.com/AdguardTeam/golibs/netutil"
)

func ExamplePrefix_UnmarshalText() {
	addr := netip.MustParseAddr("1.2.3.4")

	stdPref := netip.PrefixFrom(addr, 16)
	encPref, err := stdPref.MarshalText()
	if err != nil {
		panic(err)
	}

	p := netutil.Prefix{}
	err = p.UnmarshalText(encPref)
	fmt.Println(p, err)

	err = p.UnmarshalText([]byte("1.2.3.4"))
	fmt.Println(p, err)

	// Output:
	//
	// 1.2.3.4/16 <nil>
	// 1.2.3.4/32 <nil>
}

func ExampleUnembedPrefixes() {
	addr := netip.MustParseAddr("1.2.3.4")
	stdPref := netip.PrefixFrom(addr, 16)

	p := netutil.Prefix{
		Prefix: stdPref,
	}

	fmt.Println(netutil.UnembedPrefixes([]netutil.Prefix{p}))

	// Output:
	//
	// [1.2.3.4/16]
}
