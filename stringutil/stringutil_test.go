package stringutil_test

import (
	"strings"
	"testing"

	"github.com/AdguardTeam/golibs/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestContainsFold(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		inS      string
		inSubstr string
		want     bool
	}{{
		name:     "empty",
		inS:      "",
		inSubstr: "",
		want:     true,
	}, {
		name:     "shorter",
		inS:      "a",
		inSubstr: "abc",
		want:     false,
	}, {
		name:     "same_len_true",
		inS:      "abc",
		inSubstr: "abc",
		want:     true,
	}, {
		name:     "same_len_true_fold",
		inS:      "abc",
		inSubstr: "aBc",
		want:     true,
	}, {
		name:     "same_len_false",
		inS:      "abc",
		inSubstr: "def",
		want:     false,
	}, {
		name:     "longer_true",
		inS:      "abcdedef",
		inSubstr: "def",
		want:     true,
	}, {
		name:     "longer_false",
		inS:      "abcded",
		inSubstr: "ghi",
		want:     false,
	}, {
		name:     "longer_true_fold",
		inS:      "abcdedef",
		inSubstr: "dEf",
		want:     true,
	}, {
		name:     "longer_false_fold",
		inS:      "abcded",
		inSubstr: "gHi",
		want:     false,
	}, {
		name:     "longer_true_cyr_fold",
		inS:      "абвгдедеё",
		inSubstr: "дЕЁ",
		want:     true,
	}, {
		name:     "longer_false_cyr_fold",
		inS:      "абвгдедеё",
		inSubstr: "жЗИ",
		want:     false,
	}, {
		name:     "no_letters_true",
		inS:      "1.2.3.4",
		inSubstr: "2.3.4",
		want:     true,
	}, {
		name:     "no_letters_false",
		inS:      "1.2.3.4",
		inSubstr: "2.3.5",
		want:     false,
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.want {
				assert.True(t, stringutil.ContainsFold(tc.inS, tc.inSubstr))
			} else {
				assert.False(t, stringutil.ContainsFold(tc.inS, tc.inSubstr))
			}
		})
	}
}

func BenchmarkContainsFold(b *testing.B) {
	const s = "aaahBbBhccchDDDeEehFfFhGGGhHhh"
	const substr = "HHH"

	// Compare our implementation of containsFold against a stupid solution
	// of calling strings.ToLower and strings.Contains.
	b.Run("containsfold", func(b *testing.B) {
		var ok bool
		b.ReportAllocs()
		for b.Loop() {
			ok = stringutil.ContainsFold(s, substr)
		}

		assert.True(b, ok)
	})

	b.Run("tolower_contains", func(b *testing.B) {
		var ok bool
		b.ReportAllocs()
		for b.Loop() {
			ok = strings.Contains(strings.ToLower(s), strings.ToLower(substr))
		}

		assert.True(b, ok)
	})

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/stringutil
	//	cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
	//	BenchmarkContainsFold
	//	BenchmarkContainsFold/containsfold
	//	BenchmarkContainsFold/containsfold-16         	18405379	        65.11 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkContainsFold/tolower_contains
	//	BenchmarkContainsFold/tolower_contains-16     	 3056272	       418.4 ns/op	      40 B/op	       2 allocs/op
}
