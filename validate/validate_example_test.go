package validate_test

import (
	"fmt"
	"math"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/validate"
)

// value is a simple value that returns err in [*value.Validate].
type value struct {
	err error
}

// type check
var _ validate.Interface = (*value)(nil)

// Validate implements the [validate.Interface] interface for *Value.
func (v *value) Validate() (err error) {
	return v.err
}

func ExampleAppend() {
	var errs []error

	var (
		badValue = &value{
			err: errors.Error("test error"),
		}
		goodValue = &value{}
	)

	errs = validate.Append(errs, "first_value", goodValue)
	errs = validate.Append(errs, "second_value", badValue)

	fmt.Println(errors.Join(errs...))

	// Output:
	// second_value: test error
}

func ExampleSlice() {
	values := []*value{
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

func ExampleNotEmpty() {
	fmt.Println(validate.NotEmpty("foo", "value"))
	fmt.Println(validate.NotEmpty("foo", ""))

	type Bar struct {
		Field int
	}

	fmt.Println(validate.NotEmpty("bar", Bar{Field: 1}))
	fmt.Println(validate.NotEmpty("bar", Bar{}))

	// Output:
	// <nil>
	// foo: empty value
	// <nil>
	// bar: empty value
}

func ExampleNotEmptySlice() {
	fmt.Println(validate.NotEmptySlice("foo", []int{1}))
	fmt.Println(validate.NotEmptySlice("foo", []int(nil)))
	fmt.Println(validate.NotEmptySlice("foo", []int{}))

	// Output:
	// <nil>
	// foo: no value
	// foo: empty value
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

func Example_withNaN() {
	nan := math.NaN()

	fmt.Println(validate.InRange("foo", nan, 0, 1))
	fmt.Println(validate.NotNegative("foo", nan))
	fmt.Println(validate.Positive("foo", nan))

	// Output:
	// foo: out of range: must be no less than 0, got NaN
	// foo: negative value: NaN
	// foo: not positive: NaN
}
