package log_test

import (
	"os"

	"github.com/AdguardTeam/golibs/log"
)

func ExampleLogPanic() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	f := func() {
		defer log.LogPanic("")

		panic("fail")
	}

	f()

	f = func() {
		defer log.LogPanic("f")

		panic("fail")
	}

	f()

	// Output:
	//
	// [error] recovered from panic: fail
	// [error] f: recovered from panic: fail
}
