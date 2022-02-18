package testutil_test

import (
	"fmt"
	"strings"

	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/require"
)

func ExamplePanicT() {
	panicChan := make(chan interface{})
	requireConcurrently := func() {
		defer func() { panicChan <- recover() }()

		t := testutil.PanicT{}
		require.True(t, false)
	}

	go requireConcurrently()

	catched := <-panicChan
	catchedStr := fmt.Sprintf("%v", catched)

	// Check against the OS-independent part of the test failure message.
	fmt.Printf("contains a test failure: %t", strings.Contains(catchedStr, "Should be true"))

	// Output:
	//
	// contains a test failure: true
}
