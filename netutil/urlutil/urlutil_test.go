package urlutil_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/netutil/urlutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

// Common constants for tests.
const (
	testHostname = "hostname.example"
	testPassword = "pass"
	testUsername = "user"
)

// newTestURL is a helper function that returns a URL with dummy values.
func newTestURL(info *url.Userinfo) (u *url.URL) {
	return &url.URL{
		Scheme:   urlutil.SchemeHTTP,
		User:     info,
		Host:     testHostname,
		Path:     "/a/b/c/",
		RawQuery: "d=e",
		Fragment: "f",
	}
}

func TestRedactUserinfo(t *testing.T) {
	urlRedacted := newTestURL(url.UserPassword("xxxxx", "xxxxx"))

	testCases := []struct {
		in   *url.URL
		want *url.URL
		name string
	}{{
		in:   newTestURL(url.UserPassword(testUsername, testPassword)),
		want: urlRedacted,
		name: "with_auth",
	}, {
		in:   newTestURL(nil),
		want: newTestURL(nil),
		name: "without_auth",
	}, {
		in:   newTestURL(url.UserPassword(testUsername, "")),
		want: urlRedacted,
		name: "only_user",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, urlutil.RedactUserinfo(tc.in))
		})
	}
}

func TestRedactUserinfoInURLError(t *testing.T) {
	testURL := &url.URL{
		Scheme: urlutil.SchemeHTTP,
		User:   url.UserPassword(testUsername, testPassword),
		Host:   testHostname,
	}

	const errTest errors.Error = "test error"

	testURLError := &url.Error{
		Op:  http.MethodGet,
		URL: fmt.Sprintf("http://%s:%s@%s", testUsername, testPassword, testHostname),
		Err: errTest,
	}

	urlWithoutAuth := &url.URL{
		Scheme: urlutil.SchemeHTTP,
		Host:   testHostname,
	}

	errWithoutAuth := &url.Error{
		Op:  http.MethodGet,
		URL: urlWithoutAuth.String(),
		Err: errTest,
	}

	testCases := []struct {
		url        *url.URL
		in         error
		wantErrMsg string
		name       string
	}{{
		url:        testURL,
		in:         nil,
		wantErrMsg: "",
		name:       "nil_error",
	}, {
		url:        testURL,
		in:         testURLError,
		wantErrMsg: `GET "http://xxxxx:xxxxx@` + testHostname + `": test error`,
		name:       "redacted",
	}, {
		url:        urlWithoutAuth,
		in:         errWithoutAuth,
		wantErrMsg: `GET "http://` + testHostname + `": test error`,
		name:       "without_auth",
	}, {
		url:        testURL,
		in:         fmt.Errorf("not url error"),
		wantErrMsg: "not url error",
		name:       "not_url_error",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlutil.RedactUserinfoInURLError(tc.url, tc.in)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, tc.in)
		})
	}
}
