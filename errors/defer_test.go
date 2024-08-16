package errors_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

func ExampleAnnotate() {
	f := func(fn string) (err error) {
		defer func() {
			err = errors.Annotate(err, "reading %q: %w", fn)
		}()

		return errors.Error("not found")
	}

	err := f("non-existing")
	fmt.Println("with err    :", err)

	f = func(fn string) (err error) {
		defer func() {
			err = errors.Annotate(err, "reading %q: %w", fn)
		}()

		return nil
	}

	err = f("non-existing")
	fmt.Println("without err :", err)

	// Output:
	//
	// with err    : reading "non-existing": not found
	// without err : <nil>
}

func ExampleAnnotate_end() {
	f := func(fn string) (err error) {
		err = errors.Error("not found")

		// Some operations.

		return errors.Annotate(err, "reading %q: %w", fn)
	}

	err := f("non-existing")
	fmt.Println("with err    :", err)

	f = func(fn string) (err error) {
		// Some operations.

		return errors.Annotate(err, "reading %q: %w", fn)
	}

	err = f("non-existing")
	fmt.Println("without err :", err)

	// Output:
	//
	// with err    : reading "non-existing": not found
	// without err : <nil>
}

func ExampleAnnotate_bad() {
	const errNotFound errors.Error = "not found"

	cond := true
	g := func() (err error) { return nil }
	f := func() error {
		if cond {
			err := g()
			if err != nil {
				return err
			}

			// BAD!  This err is not the same err as the one that is
			// returned from the top level.
			defer func() { err = errors.Annotate(err, "f") }()
		}

		// This error is returned without an annotation.
		return errNotFound
	}

	// Outputs the error without an annotation.
	fmt.Println(f())

	// Output:
	//
	// not found
}

func ExampleDeferred() {
	const (
		errClose    errors.Error = "close fail"
		errNotFound errors.Error = "not found"
	)

	// logErr logs the error according to its severity level.
	logErr := func(err error) {
		if defErr := errors.Deferred(nil); errors.As(err, &defErr) && defErr.Deferred() {
			// Log deferred errors as warnings.
			fmt.Printf("warning: %s\n", errors.Unwrap(defErr))
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
	}

	// Case 1: the function fails, but the cleanup succeeds.
	closeFunc := func(fn string) (err error) { return nil }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return errNotFound
	}

	err := f("non-existing")
	logErr(err)

	// Case 2: the function succeeds, but the cleanup fails.
	closeFunc = func(_ string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return nil
	}

	err = f("non-existing")
	logErr(err)

	// Output:
	//
	// ERROR: not found
	// warning: close fail
}

func ExamplePair() {
	closeFunc := func(_ string) (err error) { return errors.Error("close fail") }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return errors.Error("not found")
	}

	err := f("non-existing")
	fmt.Println("err       :", err)
	fmt.Println("unwrapped :", errors.Unwrap(err))

	fmt.Println()

	errPair := &errors.Pair{}
	if !errors.As(err, &errPair) {
		panic("err is not an *error.Pair")
	}

	defErr := errPair.Deferred.(errors.Deferred)
	fmt.Println("deferred           :", defErr)
	fmt.Println("deferred unwrapped :", errors.Unwrap(defErr))
	fmt.Println("deferred check     :", defErr.Deferred())

	// Output:
	//
	// err       : returned: "not found", deferred: "close fail"
	// unwrapped : not found
	//
	// deferred           : deferred: close fail
	// deferred unwrapped : close fail
	// deferred check     : true
}

func ExampleWithDeferred() {
	const (
		errClose    errors.Error = "close fail"
		errNotFound errors.Error = "not found"
	)

	closeFunc := func(fn string) (err error) { return nil }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return nil
	}

	err := f("non-existing")
	fmt.Println("without errs      :", err)

	closeFunc = func(fn string) (err error) { return nil }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return errNotFound
	}

	err = f("non-existing")
	fmt.Println("with returned err :", err)

	closeFunc = func(fn string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return nil
	}

	err = f("non-existing")
	fmt.Println("with deferred err :", err)

	closeFunc = func(_ string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()

		return errNotFound
	}

	err = f("non-existing")
	fmt.Println("with both errs    :", err)

	// Output:
	//
	// without errs      : <nil>
	// with returned err : not found
	// with deferred err : deferred: close fail
	// with both errs    : returned: "not found", deferred: "close fail"
}

func ExampleWithDeferred_bad() {
	const (
		errClose    errors.Error = "close fail"
		errNotFound errors.Error = "not found"
	)

	cond := true
	closeFunc := func(_ string) (err error) { return errClose }
	g := func() (err error) { return nil }
	f := func(fn string) error {
		if cond {
			err := g()
			if err != nil {
				return err
			}

			// BAD!  This err is not the same err as the one that is
			// returned from the top level.
			defer func() { err = errors.WithDeferred(err, closeFunc(fn)) }()
		}

		// This error is the one that is actually returned.
		return errNotFound
	}

	// Only outputs the returned error and ignores the deferred one.
	fmt.Println(f("non-existing"))

	// Output:
	//
	// not found
}

func ExampleFromRecovered() {
	printRecovered := func() {
		err := errors.FromRecovered(recover())
		if err != nil {
			fmt.Printf("got error: %T(%[1]s)\n", err)
		} else {
			fmt.Println("no errors")
		}
	}

	func() {
		defer printRecovered()

		// No panic.
	}()

	func() {
		defer printRecovered()

		// Panic with error.
		panic(errors.Error("test error"))
	}()

	func() {
		defer printRecovered()

		// Panic with value.
		panic("test error")
	}()

	// Output:
	// no errors
	// got error: errors.Error(test error)
	// got error: *errors.errorString(recovered: test error)
}
