package fakeio_test

import (
	"fmt"
	"io"
	"slices"

	"github.com/AdguardTeam/golibs/testutil/fakeio"
)

func Example() {
	var written []byte
	fakeWriter := &fakeio.Writer{
		OnWrite: func(b []byte) (n int, err error) {
			written = slices.Clone(b)

			return len(b), nil
		},
	}

	// The function that is expected to call Write.
	testedFunction := func(w io.Writer) (err error) {
		_, err = io.WriteString(w, "test message")
		if err != nil {
			return fmt.Errorf("writing: %w", err)
		}

		return nil
	}

	// A simulation of a successful test.
	gotErr := testedFunction(fakeWriter)
	fmt.Printf("written: %v %q\n", gotErr, written)

	// Output:
	// written: <nil> "test message"
}
