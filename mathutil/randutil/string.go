package randutil

import (
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
		// #nosec G115 -- rl cannot be negative.
		for runenames.Name(r) == "" || rl == -1 || uint64(rl) > rest {
			r, rl = randomRune(rng, rest)
		}

		res = utf8.AppendRune(res, r)

		// #nosec G115 -- rl cannot be negative.
		rest -= uint64(rl)
	}

	return res
}

// randomRune returns a random rune r and its length l depending on how much
// space is left in the slice.  rng must not be nil.
func randomRune(rng *rand.Rand, rest uint64) (r rune, l int) {
	maxRune := maxRuneForRest(rest)
	r = rune(rng.Uint64N(maxRune))
	l = utf8.RuneLen(r)

	return r, l
}

// maxRuneForRest returns the maximum rune depending on how much space is left
// in the slice.
func maxRuneForRest(rest uint64) (maxRune uint64) {
	switch rest {
	case 1:
		return unicode.MaxASCII
	case 2:
		return 0x07ff
	case 3:
		return 0xffff
	default:
		return unicode.MaxRune
	}
}

// AppendStringASCII appends a randomly-generated string of length l, containing
// only printable ASCII characters, to orig.  rng must not be nil.
func AppendStringASCII(orig []byte, rng *rand.Rand, l uint64) (res []byte) {
	res = orig

	const (
		asciiPrintStart = ' '
		asciiPrintLast  = '~'
		asciiPrintRange = asciiPrintLast - asciiPrintStart + 1
	)

	for range l {
		b := byte(rng.Uint64N(asciiPrintRange) + asciiPrintStart)
		res = append(res, b)
	}

	return res
}

// AppendStringAlphabet appends a randomly-generated string of length l,
// containing only characters from Alphabet ab.
func AppendStringAlphabet(orig []byte, rng *rand.Rand, l uint64, ab Alphabet) (res []byte) {
	res = orig

	for range l {
		idx := rng.Uint64N(uint64(len(ab)))
		res = append(res, ab[idx])
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

// StringAlphabet returns a randomly-generated string of length l, containing
// only characters from Alphabet ab.
func StringAlphabet(rng *rand.Rand, l uint64, ab Alphabet) (s string) {
	b := make([]byte, 0, l)
	b = AppendStringAlphabet(b, rng, l, ab)

	return string(b)
}
