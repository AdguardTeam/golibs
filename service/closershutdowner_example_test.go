package service_test

import (
	"context"
	"fmt"

	"github.com/AdguardTeam/golibs/service"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
)

func ExampleCloserShutdowner() {
	isClosed := true

	closer := &fakeio.Closer{
		OnClose: func() error {
			isClosed = true

			return nil
		},
	}

	shutdowner := service.NewCloserShutdowner(closer)
	err := shutdowner.Shutdown(context.Background())
	fmt.Printf("isClosed: %t\n", isClosed)
	fmt.Printf("error:    %v\n", err)

	// Output:
	// isClosed: true
	// error:    <nil>
}
