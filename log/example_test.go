package log_test

import (
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
