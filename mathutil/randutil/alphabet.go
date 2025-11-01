package randutil

// Alphabet is an alias for strings containing alphabets for random-string
// generation.  Alphabets must contain valid UTF-8.  Alphabets should not
// contain duplicates.  It is recommended that alphabets be sorted in the
// Unicode order.
type Alphabet = string

const (
	// AlphabetNumbers contains only decimal digits.
	AlphabetNumbers Alphabet = "0123456789"

	// AlphabetLowercase contains only lowercase Latin letters.
	AlphabetLowercase Alphabet = "abcdefghijklmnopqrstuvwxyz"

	// AlphabetUppercase contains only uppercase Latin letters.
	AlphabetUppercase Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// AlphabetBase64URLSafe contains characters used in URL-safe Base64
	// encodings.
	AlphabetBase64URLSafe Alphabet = "-" + AlphabetNumbers + AlphabetUppercase + "_" + AlphabetLowercase
)
