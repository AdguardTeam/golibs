package mathutil_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/mathutil"
)

func ExampleBoolToNumber() {
	fmt.Println(mathutil.BoolToNumber[int](true))
	fmt.Println(mathutil.BoolToNumber[int](false))

	type flag float64
	fmt.Println(mathutil.BoolToNumber[flag](true))
	fmt.Println(mathutil.BoolToNumber[flag](false))

	// Output:
	// 1
	// 0
	// 1
	// 0
}

func ExampleMax() {
	fmt.Println(mathutil.Max(1, 2))
	fmt.Println(mathutil.Max(2, 1))

	// Output:
	// 2
	// 2
}

func ExampleMin() {
	fmt.Println(mathutil.Min(1, 2))
	fmt.Println(mathutil.Min(2, 1))

	// Output:
	// 1
	// 1
}
