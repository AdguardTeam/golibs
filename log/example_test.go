package log_test

import (
	"io"
	"os"

	"github.com/AdguardTeam/golibs/log"
)

func ExampleOnPanic() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	f := func() {
		defer log.OnPanic("")

		panic("fail")
	}

	f()

	f = func() {
		defer log.OnPanic("f")

		panic("fail")
	}

	f()

	// Output:
	//
	// [error] recovered from panic: fail
	// [error] f: recovered from panic: fail
}

type ErrorCloser struct{}

func (c *ErrorCloser) Close() error {
	return io.EOF
}

func ExampleOnCloserError() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	closer := &ErrorCloser{}
	f := func() {
		defer log.OnCloserError(closer, log.ERROR)
	}

	f()

	// Output:
	//
	// [error] github.com/AdguardTeam/golibs/log_test.ExampleOnCloserError.func1(): error occurred in a Close call: EOF
}
