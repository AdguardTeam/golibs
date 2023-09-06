package urlutil_test

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/AdguardTeam/golibs/netutil/urlutil"
)

// check is an error-checking helper for examples.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExampleURL() {
	type jsonStruct struct {
		Stdlib *url.URL
		Util   *urlutil.URL
	}

	const rawURL = "https://host.example:1234/path?query=1#fragment"

	stdlibURL, err := url.Parse(rawURL)
	check(err)

	utilURL, err := urlutil.Parse(rawURL)
	check(err)

	v := &jsonStruct{
		Stdlib: stdlibURL,
		Util:   utilURL,
	}

	data, err := json.MarshalIndent(v, "", "  ")
	check(err)

	fmt.Printf("%s\n", data)

	v = &jsonStruct{}
	data = []byte(`{"Util":"` + rawURL + `"}`)
	err = json.Unmarshal(data, v)
	check(err)

	fmt.Printf("%q\n", v.Util)

	// Output:
	// {
	//   "Stdlib": {
	//     "Scheme": "https",
	//     "Opaque": "",
	//     "User": null,
	//     "Host": "host.example:1234",
	//     "Path": "/path",
	//     "RawPath": "",
	//     "OmitHost": false,
	//     "ForceQuery": false,
	//     "RawQuery": "query=1",
	//     "Fragment": "fragment",
	//     "RawFragment": ""
	//   },
	//   "Util": "https://host.example:1234/path?query=1#fragment"
	// }
	// "https://host.example:1234/path?query=1#fragment"
}
