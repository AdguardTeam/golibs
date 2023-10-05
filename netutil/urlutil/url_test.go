package urlutil_test

import (
	"net/url"
	"testing"

	"github.com/AdguardTeam/golibs/netutil/urlutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	testCases := []struct {
		want       *urlutil.URL
		name       string
		in         string
		wantErrMsg string
	}{{
		want:       nil,
		name:       "empty",
		in:         "",
		wantErrMsg: "empty url",
	}, {
		want:       nil,
		name:       "bad_url",
		in:         ":",
		wantErrMsg: `parse ":": missing protocol scheme`,
	}, {
		want: &urlutil.URL{
			URL: url.URL{
				Scheme: "https",
				Host:   "www.example.com",
			},
		},
		name:       "success",
		in:         "https://www.example.com",
		wantErrMsg: "",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := urlutil.Parse(tc.in)
			assert.Equal(t, tc.want, u)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
		})
	}
}

func TestURL_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		want       *urlutil.URL
		name       string
		wantErrMsg string
		in         []byte
	}{{
		want:       &urlutil.URL{},
		name:       "nil",
		wantErrMsg: `empty json value for url`,
		in:         nil,
	}, {
		want:       &urlutil.URL{},
		name:       "empty",
		wantErrMsg: `empty url`,
		in:         []byte(`""`),
	}, {
		want:       &urlutil.URL{},
		name:       "json_null",
		wantErrMsg: ``,
		in:         []byte(`null`),
	}, {
		want:       &urlutil.URL{},
		name:       "bad_url",
		wantErrMsg: `parse ":": missing protocol scheme`,
		in:         []byte(`":"`),
	}, {
		want:       &urlutil.URL{},
		name:       "bad_type",
		wantErrMsg: `json: cannot unmarshal number into Go value of type *urlutil.URL`,
		in:         []byte(`123`),
	}, {
		want: &urlutil.URL{
			URL: url.URL{
				Scheme: "https",
				Host:   "www.example.com",
			},
		},
		name:       "success",
		wantErrMsg: "",
		in:         []byte(`"https://www.example.com"`),
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := &urlutil.URL{}
			err := u.UnmarshalJSON(tc.in)
			assert.Equal(t, tc.want, u)
			testutil.AssertErrorMsg(t, tc.wantErrMsg, err)
		})
	}
}
