package randutil

// Alphabet is an alias for strings containing alphabets for random-string
// generation.  Alphabets should not contain duplicates.  It is recommended that
// alphabets be sorted in the Unicode order.
type Alphabet = string

// Some common alphabets
const (
	AlphabetNumbers       Alphabet = "0123456789"
	AlphabetLowercase     Alphabet = "abcdefghijklmnopqrstuvwxyz"
	AlphabetUppercase     Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetBase64URLSafe Alphabet = "-" + AlphabetNumbers + AlphabetUppercase + "_" + AlphabetLowercase
)
