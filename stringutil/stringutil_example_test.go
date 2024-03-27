package stringutil_test

import (
	"fmt"
	"strings"

	"github.com/AdguardTeam/golibs/stringutil"
)

func ExampleCloneSliceOrEmpty() {
	var a, b []string

	b = stringutil.CloneSliceOrEmpty(a)
	fmt.Printf("b == nil is %t, len(b) is %d\n", b == nil, len(b))

	a = []string{}
	b = stringutil.CloneSliceOrEmpty(a)
	fmt.Printf("b == nil is %t, len(b) is %d\n", b == nil, len(b))

	a = []string{"a", "b", "c"}
	b = stringutil.CloneSliceOrEmpty(a)
	fmt.Printf("b is %v\n", b)
	fmt.Printf("&a[0] == &b[0] is %t\n", &a[0] == &b[0])

	// Output:
	//
	// b == nil is false, len(b) is 0
	// b == nil is false, len(b) is 0
	// b is [a b c]
	// &a[0] == &b[0] is false
}

func ExampleCoalesce() {
	fmt.Printf("%q\n", stringutil.Coalesce())
	fmt.Printf("%q\n", stringutil.Coalesce("", "a"))
	fmt.Printf("%q\n", stringutil.Coalesce("a", ""))
	fmt.Printf("%q\n", stringutil.Coalesce("a", "b"))

	// Output:
	//
	// ""
	// "a"
	// "a"
	// "a"
}

func ExampleContainsFold() {
	if stringutil.ContainsFold("abc", "b") {
		fmt.Println("works with the same case")
	}

	if stringutil.ContainsFold("abc", "B") {
		fmt.Println("works with a different case")
	}

	// Output:
	//
	// works with the same case
	// works with a different case
}

func ExampleFilterOut() {
	strs := []string{
		"some text",
		"",
		"# comments",
	}

	// Remove all empty and comment lines.
	filtered := stringutil.FilterOut(strs, func(s string) (ok bool) {
		return len(s) == 0 || s[0] == '#'
	})

	fmt.Printf("%q\n", filtered)

	// Output:
	//
	// ["some text"]
}

func ExampleSplitTrimmed() {
	s := ""
	fmt.Printf("%q is split into %q\n", s, stringutil.SplitTrimmed(s, ","))

	s = "a, b  ,  , c"
	fmt.Printf("%q is split into %q\n", s, stringutil.SplitTrimmed(s, ","))

	// Output:
	//
	// "" is split into []
	// "a, b  ,  , c" is split into ["a" "b" "c"]
}

func ExampleWriteToBuilder() {
	b := &strings.Builder{}

	stringutil.WriteToBuilder(
		b,
		"a",
		"b",
		"c",
	)

	fmt.Println(b)

	// Output:
	//
	// abc
}
