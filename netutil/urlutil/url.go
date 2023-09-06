// Package urlutil contains types and utilities for dealing with URLs.
package urlutil

import (
	"encoding"
	"net/url"
)

// URL is a wrapper around url.URL that can marshal and unmarshal itself from
// text form more easily.
type URL struct {
	url.URL
}

// Parse is a wrapper around [url.Parse] that returns *URL.
func Parse(rawURL string) (u *URL, err error) {
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
	return u.UnmarshalBinary(b)
}
