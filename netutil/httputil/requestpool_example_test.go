package httputil_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/httphdr"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/AdguardTeam/golibs/testutil/fakeio"
	"github.com/AdguardTeam/golibs/testutil/fakenet/fakehttp"
	"github.com/AdguardTeam/golibs/version"
)

func ExampleRequestPool() {
	// The original response origBody.
	origBody := &struct {
		io.Closer
		io.Reader
	}{
		Closer: &fakeio.Closer{
			OnClose: func() (err error) { return nil },
		},
		Reader: &fakeio.Reader{
			OnRead: func(b []byte) (n int, err error) { panic(testutil.UnexpectedCall(b)) },
		},
	}

	// The original response.
	origResp := &http.Response{
		Body: origBody,
	}

	// The original http.RoundTripper wrapped by the middleware.
	origRoundTripper := &fakehttp.RoundTripper{
		OnRoundTrip: func(_ *http.Request) (resp *http.Response, err error) {
			return origResp, nil
		},
	}

	reqPool := httputil.NewRequestPool()

	// The wrapped http.RoundTripper.
	rt := &fakehttp.RoundTripper{
		OnRoundTrip: func(req *http.Request) (resp *http.Response, err error) {
			clone, hpe := reqPool.GetWithHeader(req)
			hpe.Set(httphdr.UserAgent, "My-User-Agent/1.0")

			resp, err = origRoundTripper.RoundTrip(clone)

			reqPool.PutOnClose(resp, clone, hpe)

			return resp, err
		},
	}

	req := errors.Must(http.NewRequest(http.MethodGet, "", nil))
	allocs := testing.AllocsPerRun(10_000, func() {
		resp, err := rt.RoundTrip(req)
		errors.Check(err)

		err = resp.Body.Close()
		errors.Check(err)

		// Reset the body for the next iteration.
		origResp.Body = origBody
	})

	if version.RaceEnabled {
		// NOTE:  The race detector skews the actual number of allocations.  Run
		// this test without --race to see the actual number.
		_, err := fmt.Fprintf(os.Stderr, "overriding allocs due to --race=1, actual allocs: %.0f\n", allocs)
		errors.Check(err)

		allocs = 0
	}

	fmt.Println("Allocations:", allocs)

	// Output:
	// Allocations: 0
}
