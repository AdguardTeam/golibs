package testutil_test

import (
	"fmt"
	"strings"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func ExamplePanicT() {
	sigChan := make(chan struct{})
	catchf := func(substr, msg string) {
		caught := fmt.Sprintf("%v", recover())

		// Check against the OS-independent part of the test failure message.
		fmt.Printf("%s: %t\n", msg, strings.Contains(caught, substr))

		sigChan <- struct{}{}
	}

	t := testutil.PanicT{}

	go func() {
		defer catchf("Should be true", "contains meaningful message")

		require.True(t, false)
	}()
	<-sigChan

	go func() {
		defer catchf("test failed", "contains a test failure")

		t.FailNow()
	}()
	<-sigChan

	// Output:
	//
	// contains meaningful message: true
	// contains a test failure: true
}
