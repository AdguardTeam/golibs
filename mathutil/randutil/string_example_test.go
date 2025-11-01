package randutil_test

import (
	"fmt"
	"math/rand/v2"

	"github.com/AdguardTeam/golibs/mathutil/randutil"
)

func ExampleString() {
	rng := rand.New(rand.NewChaCha8([32]byte{}))

	for range 5 {
		fmt.Printf("%+q\n", randutil.String(rng, 16))
	}

	// Output:
	// "\U000208b5\U00102b06\U000f1b5c\u458fo"
	// "\U000118d0\U000281b8\u1db7\U000280e1\x14"
	// "\U0010d7e9\U0001e869\u6a9d\ue194\u05e6"
	// "\ue91b\U0001f3a4\U000f720f\U000fe986{"
	// "\ud5e9\U00023c56\U000fedab\U000215f2="
}

func ExampleStringASCII() {
	rng := rand.New(rand.NewChaCha8([32]byte{}))

	for range 5 {
		fmt.Printf("%q\n", randutil.StringASCII(rng, 16))
	}

	// Output:
	// "`G+KzB8t&I_XU!GE"
	// "5fX&@QEB.OeYp4_ "
	// "Z\\M=--'~I*H>=YA6"
	// "Y[:\"LU?h$GR:3i%f"
	// "UF4L*LXaOC\\O;Gvx"
}

func ExampleStringAlphabet() {
	rng := rand.New(rand.NewChaCha8([32]byte{}))
	const ab = "1234"

	for range 5 {
		fmt.Printf("%q\n", randutil.StringAlphabet(rng, 16, ab))
	}

	// Output:
	// "2341214114341141"
	// "1224442324341133"
	// "4331211134343343"
	// "2424241144341434"
	// "3312423221343313"
}

func ExampleStringAlphabet_emoji() {
	rng := rand.New(rand.NewChaCha8([32]byte{}))
	const ab = "🌑🌒🌓🌔🌕🌖🌗🌘"

	for range 5 {
		fmt.Printf("%q\n", randutil.StringAlphabet(rng, 16, ab))
	}

	// Output:
	// "🌗🌗🌗🌒🌔🌔🌓🌘🌔🌑🌔🌑🌗🌗🌔🌓"
	// "🌗🌖🌗🌖🌘🌕🌒🌓🌕🌒🌗🌘🌗🌕🌔🌗"
	// "🌑🌔🌘🌑🌓🌖🌔🌑🌘🌖🌓🌘🌔🌑🌘🌔"
	// "🌑🌕🌖🌖🌑🌘🌔🌘🌕🌗🌒🌕🌓🌗🌘🌕"
	// "🌑🌘🌓🌔🌖🌓🌒🌓🌗🌑🌔🌕🌑🌕🌓🌑"
}
