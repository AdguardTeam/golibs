package urlutil_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil/urlutil"
)

func ExampleValidateHTTPURL() {
	fmt.Println(urlutil.ValidateHTTPURL(nil))

	fmt.Println(urlutil.ValidateHTTPURL(&url.URL{
		Scheme: urlutil.SchemeGRPC,
		Host:   "host.example",
	}))

	fmt.Println(urlutil.ValidateHTTPURL(&url.URL{
		Scheme: urlutil.SchemeHTTP,
		Host:   "host.example",
	}))

	fmt.Println(urlutil.ValidateHTTPURL(&url.URL{
		Scheme: "HTTP",
		Host:   "HOST.EXAMPLE",
	}))

	// Output:
	// bad http(s) url: no value
	// bad http(s) url "grpc://host.example": scheme: bad enum value: "grpc"; want "http" or "https"
	// <nil>
	// <nil>
}

func ExampleRedactUserinfoInURLError() {
	c := &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
				return nil, errors.Error("test error")
			},
		},
	}

	u := &url.URL{
		Scheme: urlutil.SchemeHTTP,
		Host:   "does-not-exist.example",
		User:   url.UserPassword("secretUser", "secretPassword"),
	}

	_, err := c.Get(u.String())
	urlutil.RedactUserinfoInURLError(u, err)

	fmt.Printf("got error: %s", err)

	// Output:
	// got error: Get "http://xxxxx:xxxxx@does-not-exist.example": test error
}
