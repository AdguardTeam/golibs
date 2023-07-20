package pprofutil_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AdguardTeam/golibs/pprofutil"
)

func Example() {
	mux := http.NewServeMux()
	pprofutil.RoutePprof(mux)

	srv := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	fmt.Printf("shutdown error: %v\n", err)
	fmt.Printf("server error: %v\n", <-errCh)

	// Output:
	// shutdown error: <nil>
	// server error: http: Server closed
}
