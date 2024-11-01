// Package urlutil contains types and utilities for dealing with URLs.
package urlutil

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AdguardTeam/golibs/errors"
)

// Known scheme constants.
const (
	SchemeFile  = "file"
	SchemeGRPC  = "grpc"
	SchemeGRPCS = "grpcs"
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
)

// ValidateFileURL returns nil if u is a valid file URL.
//
// TODO(a.garipov):  Make the validations stricter.
func ValidateFileURL(u *url.URL) (err error) {
	if u == nil {
		return fmt.Errorf("bad file url: %w", errors.ErrNoValue)
	}

	defer func() { err = errors.Annotate(err, "bad file url %q: %w", u) }()

	if !strings.EqualFold(u.Scheme, SchemeFile) {
		return fmt.Errorf("scheme: bad value: %q; want %q", u.Scheme, SchemeFile)
	}

	return nil
}

// IsValidGRPCURLScheme returns true if s is a valid gRPC(S) URL scheme.  That
// is, [SchemeGRPC] or [SchemeGRPCS]
func IsValidGRPCURLScheme(s string) (ok bool) {
	return strings.EqualFold(s, SchemeGRPC) || strings.EqualFold(s, SchemeGRPCS)
}

// ValidateGRPCURL returns nil if u is a valid gRPC(S) URL.
//
// TODO(a.garipov):  Make the validations stricter.
func ValidateGRPCURL(u *url.URL) (err error) {
	if u == nil {
		return fmt.Errorf("bad grpc(s) url: %w", errors.ErrNoValue)
	}

	defer func() { err = errors.Annotate(err, "bad grpc(s) url %q: %w", u) }()

	if !IsValidGRPCURLScheme(u.Scheme) {
		return fmt.Errorf(
			"scheme: %w: %q; want %q or %q",
			errors.ErrBadEnumValue,
			u.Scheme,
			SchemeGRPC,
			SchemeGRPCS,
		)
	}

	return nil
}

// IsValidHTTPURLScheme returns true if s is a valid HTTP(S) URL scheme.  That
// is, [SchemeHTTP] or [SchemeHTTPS]
func IsValidHTTPURLScheme(s string) (ok bool) {
	return strings.EqualFold(s, SchemeHTTP) || strings.EqualFold(s, SchemeHTTPS)
}

// ValidateHTTPURL returns nil if u is a valid HTTP(S) URL.
//
// TODO(a.garipov):  Make the validations stricter.
func ValidateHTTPURL(u *url.URL) (err error) {
	if u == nil {
		return fmt.Errorf("bad http(s) url: %w", errors.ErrNoValue)
	}

	defer func() { err = errors.Annotate(err, "bad http(s) url %q: %w", u) }()

	if !IsValidHTTPURLScheme(u.Scheme) {
		return fmt.Errorf(
			"scheme: %w: %q; want %q or %q",
			errors.ErrBadEnumValue,
			u.Scheme,
			SchemeHTTP,
			SchemeHTTPS,
		)
	}

	return nil
}

// redactedUserinfo is the redacted userinfo shared by all redacted URLs.
var redactedUserinfo = url.UserPassword("xxxxx", "xxxxx")

// RedactUserinfo returns u if the URL does not contain any userinfo data.
// Otherwise, it returns a deep clone with both username and password redacted.
// u must not be nil.
func RedactUserinfo(u *url.URL) (redacted *url.URL) {
	if u.User == nil {
		return u
	}

	ru := *u
	ru.User = redactedUserinfo

	return &ru
}

// RedactUserinfoInURLError checks if err is a [*url.Error] and, if it is,
// replaces the underlying URL string with a redacted version of u.  u must not
// be nil.
func RedactUserinfoInURLError(u *url.URL, err error) {
	if err == nil {
		return
	}

	errURL, ok := err.(*url.Error)
	if !ok {
		return
	}

	if u.User == nil {
		return
	}

	errURL.URL = RedactUserinfo(u).String()
}
