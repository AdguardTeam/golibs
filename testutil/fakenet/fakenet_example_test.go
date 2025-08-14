package fakenet_test

import (
	"fmt"
	"io"
	"net"
	"slices"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakenet"
)

func Example() {
	var written []byte
	fakeConn := &fakenet.Conn{
		// Use OnClose with a panic to signal that Close is expected to not be
		// called.
		OnClose: func() (err error) { panic(testutil.UnexpectedCall()) },

		// Other methods implemented in the same way as Close.

		// Use OnWrite to record its argument.
		OnWrite: func(b []byte) (n int, err error) {
			written = slices.Clone(b)

			return len(b), nil
		},
	}

	// The function that is expected to call Write.
	testedFunction := func(c net.Conn) (err error) {
		_, err = io.WriteString(c, "test message")
		if err != nil {
			return fmt.Errorf("writing: %w", err)
		}

		return nil
	}

	// A simulation of a successful test.
	gotErr := testedFunction(fakeConn)
	fmt.Printf("written: %v %q\n", gotErr, written)

	// The function that is expected to not call Close.
	failingFunction := func(c net.Conn) (err error) {
		err = c.Close()
		if err != nil {
			return fmt.Errorf("closing: %w", err)
		}

		return nil
	}

	// A simulation of a failing test.
	defer func() {
		fmt.Printf("got panic: %v\n", recover())
	}()

	gotErr = failingFunction(fakeConn)
	fmt.Printf("never printed: %v\n", gotErr)

	// Output:
	// written: <nil> "test message"
	// got panic: unexpected call to github.com/AdguardTeam/golibs/testutil/fakenet.(*Conn).Close()
}
