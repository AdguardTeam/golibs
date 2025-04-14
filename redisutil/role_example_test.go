package redisutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/redisutil"
)

func ExampleNewRole() {
	r, err := redisutil.NewRole("master")
	fmt.Printf("%q %v\n", r, err)

	r, err = redisutil.NewRole("bad role")
	fmt.Printf("%q %v\n", r, err)

	// Output:
	// "master" <nil>
	// "" bad enum value: "bad role"
}
