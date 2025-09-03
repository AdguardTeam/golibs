package errors

import "fmt"

// Annotate annotates the error with the message, unless the error is nil.  The
// last verb in format must be a verb compatible with errors, for example "%w".
//
// # In Defers
//
// The primary use case for this function is to simplify code like this:
//
//	func (f *foo) doStuff(s string) (err error) {
//		defer func() {
//			if err != nil {
//				err = fmt.Errorf("bad foo %q: %w", s, err)
//			}
//		}()
//
//		// …
//	}
//
// Instead, write:
//
//	func (f *foo) doStuff(s string) (err error) {
//		defer func() { err = errors.Annotate(err, "bad foo %q: %w", s) }()
//
//		// …
//	}
//
// # At The End Of Functions
//
// Another possible use case is to simplify final checks like this:
//
//	func (f *foo) doStuff(s string) (err error) {
//		// …
//
//		if err != nil {
//			return fmt.Errorf("doing stuff with %s: %w", s, err)
//		}
//
//		return nil
//	}
//
// Instead, you could write:
//
//	func (f *foo) doStuff(s string) (err error) {
//		// …
//
//		return errors.Annotate(err, "doing stuff with %s: %w", s)
//	}
//
// # Warning
//
// This function requires that there be only ONE error named "err" in the
// function and that it is always the one that is returned.  Example (Bad)
// provides an example of the incorrect usage of Annotate.
func Annotate(err error, format string, args ...any) (annotated error) {
	if err == nil {
		return nil
	}

	return fmt.Errorf(format, append(args, err)...)
}

// WithDeferred is a helper function for deferred errors.  For example, to
// preserve errors from the Close method, replace this:
//
//	defer f.Close()
//
// With this:
//
//	defer func() { err = errors.WithDeferred(err, f.Close()) }
//
// If returned is nil and deferred is non-nil, the returned error implements the
// [Deferred] interface.  If both returned and deferred are non-nil, result has
// the underlying type of [Pair].
//
// # Warning
//
// This function requires that there be only ONE error named "err" in the
// function and that it is always the one that is returned.  Example (Bad)
// provides an example of the incorrect usage of WithDeferred.
func WithDeferred(returned, deferred error) (result error) {
	if deferred == nil {
		return returned
	}

	if returned == nil {
		return deferredError{error: deferred}
	}

	return &Pair{
		Returned: returned,
		Deferred: deferredError{error: deferred},
	}
}

// Pair is a pair of errors.  The Returned error is the main error that has been
// returned by a function.  The Deferred error is the error returned by the
// cleanup function, such as [io.Closer.Close].
//
// In pairs returned from [WithDeferred], the Deferred error always implements
// the [Deferred] interface.
type Pair struct {
	Returned error
	Deferred error
}

// type check
var _ error = (*Pair)(nil)

// Error implements the error interface for *Pair.
func (err *Pair) Error() (msg string) {
	return fmt.Sprintf("returned: %q, deferred: %q", err.Returned, Unwrap(err.Deferred))
}

// type check
var _ Wrapper = (*Pair)(nil)

// Unwrap implements the [Wrapper] interface for *Pair.  It returns the
// Returned error.
func (err *Pair) Unwrap() (unwrapped error) {
	return err.Returned
}

// Deferred is the interface for errors that were returned by cleanup functions,
// such as Close.  This is useful in APIs which desire to handle such errors
// differently, for example to log them as warnings.
//
// Method Deferred returns a bool to mirror the behavior of types like
// [net.Error] and allow implementations to decide if the error is a deferred
// one dynamically.  Users of this API must check its return value as well as
// the result [errors.As].
//
//	if errDef := errors.Deferred(nil); errors.As(err, &errDef) && errDef.Deferred() {
//		// …
//	}
//
// See https://dave.cheney.net/2014/12/24/inspecting-errors.
type Deferred interface {
	error
	Deferred() (ok bool)
}

// deferredError is a helper to implement [Deferred].
type deferredError struct {
	error
}

// type check
var _ Deferred = deferredError{}

// Deferred implements the [Deferred] interface for deferredError.
func (err deferredError) Deferred() (ok bool) {
	return true
}

// type check
var _ error = deferredError{}

// Error implements the error interface for deferredError.
func (err deferredError) Error() (msg string) {
	return fmt.Sprintf("deferred: %s", err.error)
}

// type check
var _ Wrapper = deferredError{}

// Unwrap implements the [Wrapper] interface for deferredError.
func (err deferredError) Unwrap() (unwrapped error) {
	return err.error
}

// FromRecovered checks if v, which should be a value returned by recover) is an
// error and, if it isn't, wraps it using [fmt.Errorf].  If v is nil, err is
// nil.
func FromRecovered(v any) (err error) {
	switch v := v.(type) {
	case nil:
		return nil
	case error:
		return v
	default:
		return fmt.Errorf("recovered: %v", v)
	}
}
