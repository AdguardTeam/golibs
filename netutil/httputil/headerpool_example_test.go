package httputil_test

import (
	"fmt"
	"net/http"

	"github.com/AdguardTeam/golibs/netutil/httputil"
)

func ExampleHeaderPool() {
	original := http.Header{
		"Accept":       {"*/*"},
		"Content-Type": {"application/json"},
		"Host":         {"example.com"},
		"Nil-Header":   nil,
		"User-Agent":   {"MyHTTPClient/1.0"},
	}

	p := httputil.NewHeaderPool()

	cloneEntry := p.Get(original)
	defer p.Put(cloneEntry)

	clone := cloneEntry.Header()

	fmt.Printf("before mutation:\noriginal: %#v\nclone: %#v\n", original, clone)

	clone.Set("New-Header", "value")

	fmt.Printf("after mutation:\noriginal: %#v\nclone: %#v\n", original, clone)

	// Output:
	// before mutation:
	// original: http.Header{"Accept":[]string{"*/*"}, "Content-Type":[]string{"application/json"}, "Host":[]string{"example.com"}, "Nil-Header":[]string(nil), "User-Agent":[]string{"MyHTTPClient/1.0"}}
	// clone: http.Header{"Accept":[]string{"*/*"}, "Content-Type":[]string{"application/json"}, "Host":[]string{"example.com"}, "Nil-Header":[]string(nil), "User-Agent":[]string{"MyHTTPClient/1.0"}}
	// after mutation:
	// original: http.Header{"Accept":[]string{"*/*"}, "Content-Type":[]string{"application/json"}, "Host":[]string{"example.com"}, "Nil-Header":[]string(nil), "User-Agent":[]string{"MyHTTPClient/1.0"}}
	// clone: http.Header{"Accept":[]string{"*/*"}, "Content-Type":[]string{"application/json"}, "Host":[]string{"example.com"}, "New-Header":[]string{"value"}, "Nil-Header":[]string(nil), "User-Agent":[]string{"MyHTTPClient/1.0"}}
}
