// Package stringutil contains utilities for dealing with strings.
package stringutil

import (
	"strings"
)

// CloneSliceOrEmpty returns the copy of strs or empty strings slice if strs is
// a nil slice.
func CloneSliceOrEmpty(strs []string) (clone []string) {
	return append([]string{}, strs...)
}

// CloneSlice returns the exact copy of strs.
func CloneSlice(strs []string) (clone []string) {
	if strs == nil {
		return nil
	}

	return CloneSliceOrEmpty(strs)
}

// Coalesce returns the first non-empty string.  It is named after the function
// COALESCE in SQL except that since strings in Go are non-nullable, it uses an
// empty string as a NULL value.  If strs or all it's elements are empty, it
// returns an empty string.
func Coalesce(strs ...string) (res string) {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}

	return ""
}

// FilterOut returns a copy of strs with all strings for which f returned true
// removed.
func FilterOut(strs []string, f func(s string) (ok bool)) (filtered []string) {
	for _, s := range strs {
		if !f(s) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// InSlice checks if strs contains str.
func InSlice(strs []string, str string) (ok bool) {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}

// WriteToBuilder is a convenient wrapper for strings.(*Builder).WriteString
// that deals with multiple strings and ignores errors, since they are
// guaranteed to be nil.
//
// b must not be nil.
func WriteToBuilder(b *strings.Builder, strs ...string) {
	for _, s := range strs {
		_, _ = b.WriteString(s)
	}
}
