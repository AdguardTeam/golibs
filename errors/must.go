package errors

// Check is a simple error-checking helper that panics if err is not nil.  It
// must only be used within main/entrypoint functions or in simple scripts.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Must is a helper that wraps a call to a function returning (T, error) and
// panics if the error is non-nil.  It must only be used within main/entrypoint
// functions, in simple scripts, and in variable initializations such as:
//
//	var testAddr = errors.Must(parseAddr("addr_value"))
//
// If an appropriate function already exists, for example [netip.MustParseAddr],
// Must should not be used.
func Must[T any](v T, err error) (res T) {
	Check(err)

	return v
}
