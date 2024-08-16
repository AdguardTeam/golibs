package errors_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

func ExampleError() {
	const errNotFound errors.Error = "not found"

	f := func(fn string) (err error) {
		return fmt.Errorf("opening %q: %w", fn, errNotFound)
	}

	err := f("non-existing")
	fmt.Println("err       :", err)
	fmt.Println("unwrapped :", errors.Unwrap(err))

	// Output:
	//
	// err       : opening "non-existing": not found
	// unwrapped : not found
}
