package osutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/osutil"
)

func ExampleRootDirFS() {
	fsys := osutil.RootDirFS()

	d, err := fsys.Open(".")
	if err != nil {
		fmt.Printf("opening: %v\n", err)

		return
	}

	fi, err := d.Stat()
	if err != nil {
		fmt.Printf("getting info: %v\n", err)

		return
	}

	fmt.Printf("is dir: %t\n", fi.IsDir())

	err = d.Close()
	if err != nil {
		fmt.Printf("closing: %v\n", err)

		return
	}

	// Output:
	// is dir: true
}
