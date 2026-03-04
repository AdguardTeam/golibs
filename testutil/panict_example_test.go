package testutil_test

import (
	"fmt"
	"strings"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func ExamplePanicT() {
	catchf := func(substr, msg string) {
		caught := fmt.Sprintf("%v", recover())

		// Check against the OS-independent part of the test failure
		// message.
		fmt.Printf("%s: %t\n", msg, strings.Contains(caught, substr))
	}

	t := testutil.PanicT{}

	func() {
		defer catchf("Should be true", "contains meaningful message")

		require.True(t, false)
	}()

	func() {
		defer catchf("test failed", "contains a test failure")

		t.FailNow()
	}()

	// Output:
	//
	// contains meaningful message: true
	// contains a test failure: true
}

func ExampleNewPanicT() {
	catchf := func(substr, msg string) {
		caught := fmt.Sprintf("%v", recover())

		// Check against the OS-independent part of the test failure
		// message.
		fmt.Printf("%s: %t\n", msg, strings.Contains(caught, substr))
	}

	const testName = "ExampleNewPanicT"
	calledHelper := false

	tb := newTestTB()
	tb.onName = func() (name string) { return testName }
	tb.onHelper = func() { calledHelper = true }

	pt := testutil.NewPanicT(tb)

	func() {
		defer catchf("Should be true", "contains meaningful message")

		require.True(pt, false)
	}()

	func() {
		defer catchf("test failed", "contains a test failure")

		pt.FailNow()
	}()

	func() {
		name := pt.Name()

		defer catchf(name, "contains a test name")

		// Panic here to trigger the test helper and catch the test name
		// in the panic message.
		require.True(pt, false)
	}()

	pt.Helper()
	fmt.Printf("contains a test helper: %v\n", calledHelper)

	// Output:
	//
	// contains meaningful message: true
	// contains a test failure: true
	// contains a test name: true
	// contains a test helper: true
}
