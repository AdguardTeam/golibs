package testutil

import (
	"encoding"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertMarshalText checks that the implementation of
// [encoding.TextMarshaler.MarshalText] works in all situations and results in
// the string s.
//
// See https://github.com/dominikh/go-tools/issues/911.
func AssertMarshalText[T encoding.TextMarshaler](t require.TestingT, s string, v *T) (ok bool) {
	if h, isHelper := t.(interface{ Helper() }); isHelper {
		h.Helper()
	}

	// Create a checker value.
	checkerVal := newGenericCodecChecker(v)

	// Get the expected value.
	want, err := json.Marshal(newGenericCodecChecker(&s))
	require.NoErrorf(t, err, "marshaling expected value")

	// Marshal and check against the expected value.
	b, err := json.Marshal(checkerVal)
	require.NoErrorf(t, err, "marshaling checker value")

	return assert.Equal(t, want, b)
}

// codecChecker is a generic structure for checking encoding and decoding of
// types.
type codecChecker[T any] struct {
	PtrMap map[string]*T `json:"ptr_map"`
	Map    map[string]T  `json:"map"`

	PtrValue *T `json:"ptr_value"`
	Value    T  `json:"value"`

	PtrArray [1]*T `json:"ptr_array"`
	Array    [1]T  `json:"array"`

	PtrSlice []*T `json:"ptr_slice"`
	Slice    []T  `json:"slice"`
}

// newGenericCodecChecker returns a codecChecker in which the fields are
// properly initialized with v.
func newGenericCodecChecker[T any](v *T) (c codecChecker[T]) {
	return codecChecker[T]{
		PtrMap: map[string]*T{"1": v},
		Map:    map[string]T{"1": *v},

		PtrValue: v,
		Value:    *v,

		PtrArray: [1]*T{v},
		Array:    [1]T{*v},

		PtrSlice: []*T{v},
		Slice:    []T{*v},
	}
}

// TextUnmarshaler is a constraint for pointer types that implement
// encoding.TextUnmarshaler.
type TextUnmarshaler[T any] interface {
	*T
	encoding.TextUnmarshaler
}

// AssertUnmarshalText checks that the implementation of
// [encoding.TextUnmarshaler.UnmarshalText] works in all situations and results
// in a value deeply equal to want.
func AssertUnmarshalText[T any, U TextUnmarshaler[T]](t require.TestingT, s string, v U) (ok bool) {
	if h, isHelper := t.(interface{ Helper() }); isHelper {
		h.Helper()
	}

	// Create the expected value.
	want := newGenericCodecChecker(v)

	// Create the checker value.
	got := codecChecker[T]{}

	// Marshal the expected data.
	strChecker := newGenericCodecChecker(&s)
	b, err := json.Marshal(strChecker)
	require.NoErrorf(t, err, "marshaling checker value")

	// Unmarshal into the checker value and compare.
	err = json.Unmarshal(b, &got)
	require.NoErrorf(t, err, "unmarshaling value")

	return assert.Equal(t, want, got)
}
