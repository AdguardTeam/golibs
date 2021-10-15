package errors_test

import (
	"fmt"
	"os"

	"github.com/AdguardTeam/golibs/errors"
)

// Standard Library Examples

func ExampleAs() {
	if _, err := os.Open("non-existing"); err != nil {

		var pathError *os.PathError

		if errors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}

	}

	// Output:
	//
	// Failed at path: non-existing
}

func ExampleIs() {
	if _, err := os.Open("non-existing"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	//
	// file does not exist
}

func ExampleNew() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}

	// Output:
	//
	// emit macho dwarf: elf header corrupted
}

// Extension Examples

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
	close := func(fn string) (err error) { return nil }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

		return errNotFound
	}

	err := f("non-existing")
	logErr(err)

	// Case 2: the function succeeds, but the cleanup fails.
	close = func(fn string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

		return nil
	}

	err = f("non-existing")
	logErr(err)

	// Output:
	//
	// ERROR: not found
	// warning: close fail
}

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

func ExampleList() {
	const (
		err1 errors.Error = "stage 1"
		err2 errors.Error = "stage 2"
	)

	err := errors.List("fail")
	fmt.Printf("msg only     : %q %v\n", err, errors.Unwrap(err))

	err = errors.List("fail", err1)
	fmt.Printf("msg and err  : %q %q\n", err, errors.Unwrap(err))

	err = errors.List("fail", err1, err2)
	fmt.Printf("msg and errs : %q %q\n", err, errors.Unwrap(err))

	// Output:
	//
	// msg only     : "fail" <nil>
	// msg and err  : "fail: stage 1" "stage 1"
	// msg and errs : "fail: 2 errors: \"stage 1\", \"stage 2\"" "stage 1"
}

func ExamplePair() {
	close := func(fn string) (err error) { return errors.Error("close fail") }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

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

	close := func(fn string) (err error) { return nil }
	f := func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

		return nil
	}

	err := f("non-existing")
	fmt.Println("without errs      :", err)

	close = func(fn string) (err error) { return nil }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

		return errNotFound
	}

	err = f("non-existing")
	fmt.Println("with returned err :", err)

	close = func(fn string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

		return nil
	}

	err = f("non-existing")
	fmt.Println("with deferred err :", err)

	close = func(fn string) (err error) { return errClose }
	f = func(fn string) (err error) {
		defer func() { err = errors.WithDeferred(err, close(fn)) }()

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
	close := func(fn string) (err error) { return errClose }
	g := func() (err error) { return nil }
	f := func(fn string) error {
		if cond {
			err := g()
			if err != nil {
				return err
			}

			// BAD!  This err is not the same err as the one that is
			// returned from the top level.
			defer func() { err = errors.WithDeferred(err, close(fn)) }()
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
