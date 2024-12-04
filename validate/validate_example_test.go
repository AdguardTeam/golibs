package validate_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
)

// Value is a simple value that returns err in [*Value.Validate].
type Value struct {
	err error
}

// type check
var _ validate.Interface = (*Value)(nil)

// Validate implements the [validate.Interface] interface for *Value.
func (v *Value) Validate() (err error) {
	return v.err
}

func ExampleSlice() {
	values := []*Value{
		0: {
			err: nil,
		},
		1: {
			err: errors.Error("test error 1"),
		},
		2: {
			err: errors.Error("test error 2"),
		},
	}

	fmt.Println(validate.Slice("values", values))

	// Output:
	// values: at index 1: test error 1
	// values: at index 2: test error 2
}

func ExampleInRange() {
	fmt.Println(validate.InRange("foo", 0, 0, 100))
	fmt.Println(validate.InRange("foo", 100, 0, 100))
	fmt.Println(validate.InRange("foo", 101, 0, 100))
	fmt.Println(validate.InRange("foo", -1, 0, 100))

	// Output:
	// <nil>
	// <nil>
	// foo: out of range: must be no greater than 100, got 101
	// foo: out of range: must be no less than 0, got -1
}

func ExampleNotNegative() {
	fmt.Println(validate.NotNegative("foo", 1))
	fmt.Println(validate.NotNegative("foo", 0))
	fmt.Println(validate.NotNegative("foo", -1))

	// Output:
	// <nil>
	// <nil>
	// foo: negative value: -1
}

func ExampleNotNil() {
	v := 1
	fmt.Println(validate.NotNil("foo", &v))
	fmt.Println(validate.NotNil("foo", (*int)(nil)))

	// Output:
	// <nil>
	// foo: no value
}

func ExamplePositive() {
	fmt.Println(validate.Positive("foo", 1))
	fmt.Println(validate.Positive("foo", 0))
	fmt.Println(validate.Positive("foo", -1))

	// Output:
	// <nil>
	// foo: not positive: 0
	// foo: not positive: -1
}
