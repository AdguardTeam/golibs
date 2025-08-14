package testutil

import (
	"bytes"
	"context"
	"fmt"
	"runtime"

	"github.com/AdguardTeam/golibs/internal/reflectutil"
)

// UnexpectedCall is a helper for creating errors about unexpected calls in fake
// implementations of interfaces.
//
// See examples to learn how to use this function correctly.
//
// The error message is composed using the following formatting verbs:
//   - The first [context.Context] value is printed as "ctx".
//   - Nilable values that are nil are printed with the "%#v" verb.
//   - Strings are printed with the "%q" verb.
//   - The rest are printed using the "%+v" verb.
func UnexpectedCall(args ...any) (err error) {
	const N = 1
	var pcs [N]uintptr
	numCallers := runtime.Callers(3, pcs[:])
	if numCallers != N {
		panic(fmt.Errorf("runtime.Callers did not fill pcs: got %d items", numCallers))
	}

	fs := runtime.CallersFrames(pcs[:])
	f, _ := fs.Next()
	if f == (runtime.Frame{}) {
		panic(fmt.Errorf("no iteration in runtime.Frames"))
	}

	argsBuf := &bytes.Buffer{}
	_ = argsBuf.WriteByte('(')
	for i, arg := range args {
		printErr := printArg(argsBuf, arg, i)
		if printErr != nil {
			// Should not happen, as the formatting verb must be correct for the
			// type.
			panic(fmt.Errorf("arg at index %d: %w", i, printErr))
		}

		if i < len(args)-1 {
			_, _ = argsBuf.WriteString(", ")
		}
	}

	_ = argsBuf.WriteByte(')')

	return fmt.Errorf("unexpected call to %s%s", f.Function, argsBuf)
}

// printArg prints arg to buf using special formatting.
func printArg(buf *bytes.Buffer, arg any, idx int) (err error) {
	if _, ok := arg.(context.Context); ok && idx == 0 {
		_, err = buf.WriteString("ctx")
		if err != nil {
			return fmt.Errorf("printing ctx: %w", err)
		}

		return nil
	}

	verb := reflectutil.FormatVerb(arg)

	_, err = fmt.Fprintf(buf, verb, arg)
	if err != nil {
		return fmt.Errorf("printing value: %w", err)
	}

	return nil
}
