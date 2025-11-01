package randutil_test

import (
	"strconv"
	"testing"
	"testing/quick"
	"unicode/utf8"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/runenames"
)

const (
	// testQuickCheckCount defines how many quick checks to make.
	testQuickCheckCount = 100_000

	// testMaxLen is the maximum length of a string in tests.
	testMaxLen = 1024
)

func TestASCIIString(t *testing.T) {
	t.Parallel()

	isValid := func(l uint64) (ok bool) {
		// Prevent excessive memory usage.
		if l > testMaxLen {
			return true
		}

		s := randutil.StringASCII(testRNG, l)

		return assertASCII(t, []byte(s))
	}

	err := quick.Check(isValid, &quick.Config{
		MaxCount: testQuickCheckCount,
	})
	assert.NoError(t, err)
}

// assertASCII is a helper that asserts that data contains only printable ASCII
// characters.
func assertASCII(tb testing.TB, data []byte) (ok bool) {
	tb.Helper()

	for _, b := range data {
		if b < ' ' || b > '~' {
			tb.Errorf("got bad value %q with seed %v", data, testSeed)

			return false
		}
	}

	return true
}

func TestAppendASCIIString(t *testing.T) {
	t.Parallel()

	const max = testMaxLen
	data := make([]byte, 0, max)

	isValid := func(l uint64) (ok bool) {
		defer func() { data = data[:0] }()

		// Prevent excessive memory usage.
		if l > testMaxLen {
			return true
		}

		data = randutil.AppendStringASCII(data, testRNG, l)

		return assertASCII(t, data)
	}

	err := quick.Check(isValid, &quick.Config{
		MaxCount: testQuickCheckCount,
	})
	assert.NoError(t, err)
}

func TestAppendString(t *testing.T) {
	t.Parallel()

	const max = testMaxLen
	data := make([]byte, 0, max)

	isValid := func(l uint64) (ok bool) {
		defer func() { data = data[:0] }()

		// Prevent excessive memory usage.
		if l > testMaxLen {
			return true
		}

		data = randutil.AppendString(data, testRNG, l)

		return assertUnicode(t, string(data))
	}

	err := quick.Check(isValid, &quick.Config{
		MaxCount: testQuickCheckCount,
	})
	assert.NoError(t, err)
}

// assertUnicode is a helper that asserts that s contains only valid Unicode
// runes.
func assertUnicode(tb testing.TB, s string) (ok bool) {
	tb.Helper()

	for _, r := range s {
		if runenames.Name(r) == "" {
			tb.Errorf("got bad value %q with seed %v", s, testSeed)

			return false
		}
	}

	return true
}

func TestAppendString_small(t *testing.T) {
	t.Parallel()

	for l := uint64(1); l <= 3; l++ {
		t.Run(strconv.FormatUint(l, 10), func(t *testing.T) {
			t.Parallel()

			data := make([]byte, 0, l)

			isValid := func() (ok bool) {
				defer func() { data = data[:0] }()

				data = randutil.AppendString(data, testRNG, l)

				return utf8.Valid(data)
			}

			err := quick.Check(isValid, &quick.Config{
				MaxCount: testQuickCheckCount,
			})
			assert.NoError(t, err)
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	const max = testMaxLen
	data := make([]byte, 0, max)

	isValid := func(l uint64) (ok bool) {
		defer func() { data = data[:0] }()

		// Prevent excessive memory usage.
		if l > testMaxLen {
			return true
		}

		s := randutil.String(testRNG, l)

		return assertUnicode(t, s)
	}

	err := quick.Check(isValid, &quick.Config{
		MaxCount: testQuickCheckCount,
	})
	assert.NoError(t, err)
}

func BenchmarkAppendString(b *testing.B) {
	data := make([]byte, 0, testMaxLen)

	b.ReportAllocs()
	for b.Loop() {
		data = randutil.AppendString(data[:0], testRNG, testMaxLen)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/mathutil/randutil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkAppendString-24    	   57822	     18256 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkAppendStringASCII(b *testing.B) {
	data := make([]byte, 0, testMaxLen)

	b.ReportAllocs()
	for b.Loop() {
		data = randutil.AppendStringASCII(data[:0], testRNG, testMaxLen)
	}

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/mathutil/randutil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkAppendStringASCII-24    	  261384	      4527 ns/op	       0 B/op	       0 allocs/op
}

func BenchmarkAppendStringAlphabet(b *testing.B) {
	data := make([]byte, 0, testMaxLen)

	require.True(b, b.Run("bytes=1", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			data = randutil.AppendStringAlphabet(
				data[:0],
				testRNG,
				testMaxLen,
				randutil.AlphabetBase64URLSafe,
			)
		}

		require.True(b, utf8.Valid(data))
	}))

	require.True(b, b.Run("bytes=2", func(b *testing.B) {
		const ab = "ΎΏΐΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩΪΫάέήίΰαβγδεζηθικλμνξοπρςστυφχψωϊϋόύώ"
		b.ReportAllocs()
		for b.Loop() {
			data = randutil.AppendStringAlphabet(data[:0], testRNG, testMaxLen, ab)
		}

		require.True(b, utf8.Valid(data))
	}))

	// Most recent results:
	//	goos: linux
	//	goarch: amd64
	//	pkg: github.com/AdguardTeam/golibs/mathutil/randutil
	//	cpu: AMD Ryzen AI 9 HX PRO 370 w/ Radeon 890M
	//	BenchmarkAppendStringAlphabet/bytes=1-24         	  216681	      5682 ns/op	       0 B/op	       0 allocs/op
	//	BenchmarkAppendStringAlphabet/bytes=2-24         	  103119	     11637 ns/op	       0 B/op	       0 allocs/op
}
