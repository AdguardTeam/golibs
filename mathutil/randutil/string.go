package randutil

import (
	"math"
	"math/rand/v2"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/runenames"
)

// AppendString appends a randomly-generated string of length l, containing only
// valid Unicode characters, to orig.  rng must not be nil.
func AppendString(orig []byte, rng *rand.Rand, l uint64) (res []byte) {
	res = orig

	rest := l
	for rest != 0 {
		r, rl := randomRune(rng, rest)
		// #nosec G115 -- The first case removes negative numbers, so the
		// conversion is safe
		for rl < 0 || uint64(rl) > rest || runenames.Name(r) == "" {
			r, rl = randomRune(rng, rest)
		}

		res = utf8.AppendRune(res, r)

		// #nosec G115 -- By now rl must not be negative.
		rest -= uint64(rl)
	}

	return res
}

// randomRune returns a random rune r and its length l depending on how much
// space is left in the slice.  rng must not be nil.
func randomRune(rng *rand.Rand, rest uint64) (r rune, l int) {
	maxRune := maxRuneForRest(rest)
	// #nosec G115 -- maxRune shouldn't be larger than [math.MaxInt32].
	r = rune(rng.Uint64N(maxRune + 1))
	l = utf8.RuneLen(r)

	return r, l
}

// Constants for maximum values depending on the rune length.
const (
	maxRune2Bytes = 0x07ff
	maxRune3Bytes = 0xffff
)

// maxRuneForRest returns the maximum rune depending on how much space is left
// in the slice.
func maxRuneForRest(rest uint64) (maxRune uint64) {
	switch rest {
	case 1:
		return unicode.MaxASCII
	case 2:
		return maxRune2Bytes
	case 3:
		return maxRune3Bytes
	default:
		return unicode.MaxRune
	}
}

// AppendStringASCII appends a randomly-generated string of length l, containing
// only printable ASCII characters, to orig.  rng must not be nil.
func AppendStringASCII(orig []byte, rng *rand.Rand, l uint64) (res []byte) {
	res = orig

	const (
		asciiPrintFirst = ' '
		asciiPrintLast  = '~'
		asciiPrintRange = asciiPrintLast - asciiPrintFirst + 1
	)

	for range l {
		// #nosec G115 -- The value shouldn't be larger than [math.MaxUint8].
		b := byte(rng.Uint64N(asciiPrintRange) + asciiPrintFirst)
		res = append(res, b)
	}

	return res
}

// AppendStringAlphabet appends a randomly-generated string with runeLen runes,
// containing only characters from Alphabet ab, to orig.
func AppendStringAlphabet(orig []byte, rng *rand.Rand, runeLen uint64, ab Alphabet) (res []byte) {
	res = orig

	for range runeLen {
		idx := rng.Uint64N(uint64(len(ab)))

		// Rewind the index if it doesn't point to a valid rune start.
		for !utf8.RuneStart(ab[idx]) && idx != math.MaxUint64 {
			idx--
		}

		r, _ := utf8.DecodeRuneInString(ab[idx:])
		res = utf8.AppendRune(res, r)
	}

	return res
}

// String returns a randomly-generated string of length l, containing only valid
// Unicode characters.  rng must not be nil.
func String(rng *rand.Rand, l uint64) (s string) {
	b := make([]byte, 0, l)
	b = AppendString(b, rng, l)

	return string(b)
}

// StringASCII returns a randomly-generated string of length l, containing only
// printable ASCII characters.  rng must not be nil.
func StringASCII(rng *rand.Rand, l uint64) (s string) {
	b := make([]byte, 0, l)
	b = AppendStringASCII(b, rng, l)

	return string(b)
}

// StringAlphabet returns a randomly-generated string with at runeLen runes,
// containing only characters from Alphabet ab, to orig.
func StringAlphabet(rng *rand.Rand, l uint64, ab Alphabet) (s string) {
	b := make([]byte, 0, l)
	b = AppendStringAlphabet(b, rng, l, ab)

	return string(b)
}
