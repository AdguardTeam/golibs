package log_test

import (
	"io"
	"os"

	"github.com/AdguardTeam/golibs/log"
)

func ExampleLevel() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	log.Info("printed")

	log.SetLevel(log.OFF)
	log.Info("not printed")

	// Output:
	//
	// [info] printed
}

func ExampleOnPanic() {
	log.SetFlags(0)
	log.SetLevel(log.INFO)
	log.SetOutput(os.Stdout)

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
	log.SetFlags(0)
	log.SetLevel(log.INFO)
	log.SetOutput(os.Stdout)

	closer := &ErrorCloser{}
	f := func() {
		defer log.OnCloserError(closer, log.ERROR)
	}

	f()

	// Output:
	//
	// [error] github.com/AdguardTeam/golibs/log_test.ExampleOnCloserError.func1(): error occurred in a Close call: EOF
}

func ExamplePanic() {
	log.SetFlags(0)
	log.SetLevel(log.INFO)
	log.SetOutput(os.Stdout)

	defer log.OnPanic("")

	log.Panic("fail")

	// Output:
	//
	// [panic] fail
	// [error] recovered from panic: fail
}

func ExamplePanicf() {
	log.SetFlags(0)
	log.SetLevel(log.INFO)
	log.SetOutput(os.Stdout)

	defer log.OnPanic("")

	log.Panicf("fail, some number: %d", 123)

	// Output:
	//
	// [panic] fail, some number: 123
	// [error] recovered from panic: fail, some number: 123
}
