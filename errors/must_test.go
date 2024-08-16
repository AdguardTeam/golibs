package errors_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

func ExampleCheck() {
	defer func() {
		fmt.Println("recovered:", errors.FromRecovered(recover()))
	}()

	runFoo := func() (err error) {
		return errors.Error("test error")
	}

	errors.Check(runFoo())

	// Output:
	// recovered: test error
}

func ExampleMust() {
	newFooBad := func() (foo string, err error) {
		return "", errors.Error("test error")
	}

	newFooGood := func() (foo string, err error) {
		return "foo", nil
	}

	func() {
		fmt.Println("good:", errors.Must(newFooGood()))
	}()

	func() {
		defer func() {
			fmt.Println("bad:", errors.FromRecovered(recover()))
		}()

		fmt.Println(errors.Must(newFooBad()))
	}()

	// Output:
	// good: foo
	// bad: test error
}
