package testutil_test

import (
	"context"
	"fmt"
	"net"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/testutil"
)

func runAndPrintPanic(f func()) {
	defer func() { fmt.Println(recover()) }()

	f()
}

type fakeWithArgs struct {
	onCall func(first string, args ...any)
}

func (f *fakeWithArgs) Call(first string, args ...any) { f.onCall(first, args...) }

func ExampleUnexpectedCall_args() {
	f := &fakeWithArgs{
		onCall: func(first string, args ...any) {
			panic(testutil.UnexpectedCall(first, args))
		},
	}

	runAndPrintPanic(func() { f.Call("") })
	runAndPrintPanic(func() { f.Call("", 1) })
	runAndPrintPanic(func() { f.Call("", 1, 2, 3) })

	// Output:
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithArgs).Call("", []interface {}(nil))
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithArgs).Call("", [1])
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithArgs).Call("", [1 2 3])
}

type fakeWithNoArgs struct {
	onCall func()
}

func (f *fakeWithNoArgs) Call() { f.onCall() }

func ExampleUnexpectedCall_empty() {
	f := &fakeWithNoArgs{
		onCall: func() {
			panic(testutil.UnexpectedCall())
		},
	}

	runAndPrintPanic(func() { f.Call() })

	// Output:
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithNoArgs).Call()
}

type fakeWithError struct {
	onCall func(ctx context.Context, err error)
}

func (f *fakeWithError) Call(ctx context.Context, err error) { f.onCall(ctx, err) }

func ExampleUnexpectedCall_interface() {
	f := &fakeWithError{
		onCall: func(ctx context.Context, err error) {
			panic(testutil.UnexpectedCall(ctx, err))
		},
	}

	ctx := context.Background()

	var errPtr *net.OpError
	var nilPtrIface error = errPtr

	runAndPrintPanic(func() { f.Call(ctx, nil) })
	runAndPrintPanic(func() { f.Call(ctx, nilPtrIface) })
	runAndPrintPanic(func() { f.Call(ctx, errors.Error("test")) })

	// Output:
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithError).Call(ctx, <nil>)
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithError).Call(ctx, (*net.OpError)(nil))
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithError).Call(ctx, "test")
}

type fakeWithMap struct {
	onCall func(ctx context.Context, m map[string]string)
}

func (f *fakeWithMap) Call(ctx context.Context, m map[string]string) { f.onCall(ctx, m) }

func ExampleUnexpectedCall_map() {
	f := &fakeWithMap{
		onCall: func(ctx context.Context, m map[string]string) {
			panic(testutil.UnexpectedCall(ctx, m))
		},
	}

	ctx := context.Background()

	runAndPrintPanic(func() { f.Call(ctx, nil) })
	runAndPrintPanic(func() { f.Call(ctx, map[string]string{}) })
	runAndPrintPanic(func() { f.Call(ctx, map[string]string{"1": "a"}) })

	// Output:
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithMap).Call(ctx, map[string]string(nil))
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithMap).Call(ctx, map[])
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithMap).Call(ctx, map[1:a])
}

type fakeWithSlice struct {
	onCall func(ctx context.Context, b []byte)
}

func (f *fakeWithSlice) Call(ctx context.Context, b []byte) { f.onCall(ctx, b) }

func ExampleUnexpectedCall_slice() {
	f := &fakeWithSlice{
		onCall: func(ctx context.Context, b []byte) {
			panic(testutil.UnexpectedCall(ctx, b))
		},
	}

	ctx := context.Background()

	runAndPrintPanic(func() { f.Call(ctx, nil) })
	runAndPrintPanic(func() { f.Call(ctx, []byte{}) })
	runAndPrintPanic(func() { f.Call(ctx, []byte{0}) })

	// Output:
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithSlice).Call(ctx, []byte(nil))
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithSlice).Call(ctx, [])
	// unexpected call to github.com/AdguardTeam/golibs/testutil_test.(*fakeWithSlice).Call(ctx, [0])
}
