package pprofutil_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/AdguardTeam/golibs/pprofutil"
)

// must is a helper for tests.
func must[T any](v T, err error) (res T) {
	if err != nil {
		panic(err)
	}

	return v
}

func Example() {
	mux := http.NewServeMux()
	pprofutil.RoutePprof(mux)

	srv := httptest.NewServer(mux)
	defer srv.Close()

	u := must(url.Parse(srv.URL)).JoinPath(pprofutil.BasePath)
	req := must(http.NewRequest(http.MethodGet, u.String(), nil))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	resp := must(http.DefaultClient.Do(req))
	respBody := must(io.ReadAll(resp.Body))

	lines := bytes.Split(respBody, []byte("\n"))
	for i := range 15 {
		fmt.Printf("%s\n", lines[i])
	}

	fmt.Println("…")

	// Output:
	// <html>
	// <head>
	// <title>/debug/pprof/</title>
	// <style>
	// .profile-name{
	// 	display:inline-block;
	// 	width:6rem;
	// }
	// </style>
	// </head>
	// <body>
	// /debug/pprof/
	// <br>
	// <p>Set debug=1 as a query parameter to export in legacy text format</p>
	// <br>
	// …
}
