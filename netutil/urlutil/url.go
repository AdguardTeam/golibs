// Package urlutil contains types and utilities for dealing with URLs.
package urlutil

import (
	"bytes"
	"encoding"
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/AdguardTeam/golibs/errors"
)

// ErrEmpty is returned from [Parse] and [URL.UnmarshalText] when the input is
// empty.
const ErrEmpty errors.Error = "empty url"

// URL is a wrapper around url.URL that can marshal and unmarshal itself from
// text form more easily.
type URL struct {
	url.URL
}

// Parse is a wrapper around [url.Parse] that returns *URL.  Unlike url.Parse,
// it does not consider empty string a valid URL and returns [ErrEmpty]
func Parse(rawURL string) (u *URL, err error) {
	if rawURL == "" {
		return nil, ErrEmpty
	}

	uu, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &URL{
		URL: *uu,
	}, nil
}

// type check
var _ encoding.TextMarshaler = (*URL)(nil)

// MarshalText implements the [encoding.TextMarshaler] interface for *URL.
func (u *URL) MarshalText() (b []byte, err error) {
	return u.MarshalBinary()
}

// type check
var _ encoding.TextUnmarshaler = (*URL)(nil)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface for *URL.
func (u *URL) UnmarshalText(b []byte) (err error) {
	if len(b) == 0 {
		return errors.Error("empty url")
	}

	return u.UnmarshalBinary(b)
}

// type check
var _ json.Unmarshaler = (*URL)(nil)

// UnmarshalJSON implements the [json.Unmarshaler] interface for *URL.
func (u *URL) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	l := len(b)
	if l == 0 {
		return errors.Error("empty json value for url")
	}

	if b[0] != '"' || b[l-1] != '"' {
		// Try to create a type error with the Value field set.  If unable, just
		// use the more general description of the value.
		//
		// TODO(a.garipov): Use a better API for this if encoding/json ever
		// supports one.
		var s string
		err = json.Unmarshal(b, &s)
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			typeErr.Type = reflect.TypeOf(u)

			return typeErr
		}

		return &json.UnmarshalTypeError{
			Value: "non-string",
			Type:  reflect.TypeOf(u),
		}
	}

	return u.UnmarshalText(b[1 : l-1])
}
