// Package validate contains functions for validating values.
//
// The common argument name can be the name of the JSON or YAML property, the
// name of a function argument, or anything similar.
//
// NOTE:  More specific validations, like those of network addresses or URLs,
// should be put into related utility packages.
//
// TODO(a.garipov):  Add a function that validates for both nilness and
// emptiness.
//
// TODO(a.garipov):  Consider adding validate.KeyValues.
package validate

import (
	"cmp"
	"fmt"

	"github.com/AdguardTeam/golibs/errors"
)

// Interface is the interface for configuration entities that can validate
// themselves.
type Interface interface {
	// Validate returns an error if the entity isn't valid.  Entities should not
	// add a prefix; instead, the callers should add prefixes depending on the
	// use.
	Validate() (err error)
}

// Append validates v and, if it returns an error, appends it to errs and
// returns the result.
func Append(errs []error, name string, v Interface) (res []error) {
	res = errs
	err := v.Validate()
	if err != nil {
		res = append(res, fmt.Errorf("%s: %w", name, err))
	}

	return res
}

// AppendSlice validates values, wraps errors with the name and the index,
// appends them to errs, and returns the result.
func AppendSlice[T Interface](errs []error, name string, values []T) (res []error) {
	res = errs
	for i, v := range values {
		// TODO(a.garipov):  Consider flattening error slices.
		err := v.Validate()
		if err != nil {
			res = append(res, fmt.Errorf("%s: at index %d: %w", name, i, err))
		}
	}

	return res
}

// Slice validates values, wraps errors with the name and the index, and returns
// the result as a single joined error.
func Slice[T Interface](name string, values []T) (err error) {
	return errors.Join(AppendSlice(nil, name, values)...)
}

// Empty returns an error if v is not equal to its zero value.  The underlying
// error of err is [errors.ErrNotEmpty].
func Empty[T comparable](name string, v T) (err error) {
	var zero T
	if v != zero {
		return fmt.Errorf("%s: %w", name, errors.ErrNotEmpty)
	}

	return nil
}

// EmptySlice returns an error if v is neither nil nor empty.  The underlying
// error of err is either [errors.ErrNotEmpty].
//
// TODO(a.garipov):  Find ways of extending to other nilable types with length.
func EmptySlice[T any](name string, v []T) (err error) {
	if len(v) > 0 {
		return fmt.Errorf("%s: %w", name, errors.ErrNotEmpty)
	}

	return nil
}

// Equal returns an error if got is not equal to want.  The underlying error of
// err is [errors.ErrNotEqual].
func Equal[T comparable](name string, got, want T) (err error) {
	if got != want {
		return fmt.Errorf("%s: %w: got %v, want %v", name, errors.ErrNotEqual, got, want)
	}

	return nil
}

// GreaterThan returns an error if a is less than or equal to b.  The underlying
// error of err is [errors.ErrOutOfRange].
//
// NOTE:  NaN is also considered less than anything, since [cmp.Compare] sorts
// it below -Infinity.
func GreaterThan[T cmp.Ordered](name string, a, b T) (err error) {
	// Use cmp.Compare to get consistent results even with NaN, which will be
	// below any other value.
	if cmp.Compare(a, b) <= 0 {
		return fmt.Errorf(
			"%s: %w: must be greater than %v, got %v",
			name,
			errors.ErrOutOfRange,
			b,
			a,
		)
	}

	return nil
}

// InRange returns an error of v is less than min or greater than max.  The
// underlying error of err is [errors.ErrOutOfRange].
//
// NOTE:  NaN is also considered less than anything, since [cmp.Compare] sorts
// it below -Infinity.
func InRange[T cmp.Ordered](name string, v, min, max T) (err error) {
	if err = NoLessThan(name, v, min); err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return err
	} else if err = NoGreaterThan(name, v, max); err != nil {
		// Don't wrap the error, because it's informative enough as is.
		return err
	}

	return nil
}

// LessThan returns an error if a is greater than or equal to b.  The underlying
// error of err is [errors.ErrOutOfRange].
//
// NOTE:  NaN is also considered less than anything, since [cmp.Compare] sorts
// it below -Infinity.
func LessThan[T cmp.Ordered](name string, a, b T) (err error) {
	// Use cmp.Compare to get consistent results even with NaN, which will be
	// below any other value.
	if cmp.Compare(a, b) >= 0 {
		return fmt.Errorf(
			"%s: %w: must be less than %v, got %v",
			name,
			errors.ErrOutOfRange,
			b,
			a,
		)
	}

	return nil
}

// Nil returns an error if v is not nil.  The underlying error of err is
// [errors.ErrUnexpectedValue].
//
// For checking against emptiness (comparing with the zero value), prefer
// [Empty].
//
// TODO(a.garipov):  Find ways of extending to other nilable types.
func Nil[T any](name string, v *T) (err error) {
	if v != nil {
		return fmt.Errorf("%s: %w", name, errors.ErrUnexpectedValue)
	}

	return nil
}

// NoGreaterThan returns an error if v is greater than max.  The underlying
// error of err is [errors.ErrOutOfRange].
//
// NOTE:  NaN is also considered less than anything, since [cmp.Compare] sorts
// it below -Infinity.
func NoGreaterThan[T cmp.Ordered](name string, v, max T) (err error) {
	// Use cmp.Compare to get consistent results even with NaN, which will be
	// below any other value.
	if cmp.Compare(v, max) > 0 {
		return fmt.Errorf(
			"%s: %w: must be no greater than %v, got %v",
			name,
			errors.ErrOutOfRange,
			max,
			v,
		)
	}

	return nil
}

// NoLessThan returns an error if v is less than min.  The underlying error of
// err is [errors.ErrOutOfRange].
//
// NOTE:  NaN is also considered less than anything, since [cmp.Compare] sorts
// it below -Infinity.
func NoLessThan[T cmp.Ordered](name string, v, min T) (err error) {
	// Use cmp.Compare to get consistent results even with NaN, which will be
	// below any other value.
	if cmp.Compare(v, min) < 0 {
		return fmt.Errorf(
			"%s: %w: must be no less than %v, got %v",
			name,
			errors.ErrOutOfRange,
			min,
			v,
		)
	}

	return nil
}

// NotEmpty returns an error if v is its zero value.  The underlying error of
// err is [errors.ErrEmpty].
//
// For pointers, prefer [NotNil].
func NotEmpty[T comparable](name string, v T) (err error) {
	var zero T
	if v == zero {
		return fmt.Errorf("%s: %w", name, errors.ErrEmptyValue)
	}

	return nil
}

// NotEmptySlice returns an error if v is nil or empty.  The underlying error of
// err is either [errors.ErrNoValue] or [errors.ErrEmpty] correspondingly.
//
// TODO(a.garipov):  Find ways of extending to other nilable types with length.
func NotEmptySlice[T any](name string, v []T) (err error) {
	if v == nil {
		return fmt.Errorf("%s: %w", name, errors.ErrNoValue)
	} else if len(v) == 0 {
		return fmt.Errorf("%s: %w", name, errors.ErrEmptyValue)
	}

	return nil
}

// NotNegative returns an error if v is less than the zero value of type T.  The
// underlying error of err is [errors.ErrNegative].
//
// NOTE:  NaN is also considered negative, since [cmp.Compare] sorts it below
// -Infinity.
func NotNegative[T cmp.Ordered](name string, v T) (err error) {
	var zero T

	// Use cmp.Compare to get consistent results even with NaN, which will be
	// below any other value.
	if cmp.Compare(v, zero) < 0 {
		return fmt.Errorf("%s: %w: %v", name, errors.ErrNegative, v)
	}

	return nil
}

// NotNil returns an error if v is nil.  The underlying error of err is
// [errors.ErrNoValue].
//
// For checking against emptiness (comparing with the zero value), prefer
// [NotEmpty].
//
// TODO(a.garipov):  Find ways of extending to other nilable types.
func NotNil[T any](name string, v *T) (err error) {
	if v == nil {
		return fmt.Errorf("%s: %w", name, errors.ErrNoValue)
	}

	return nil
}

// NotNilInterface returns an error if v is nil.  The underlying error of err is
// [errors.ErrNoValue].
//
// For checking against emptiness (comparing with the zero value), prefer
// [NotEmpty].
//
// NOTE:  This function returns an error only if v is a nil interface value.
// This means that if v is an interface value with a type and a nil pointer, err
// is nil.
//
// TODO(a.garipov):  Find ways of merging with [NotNil].
func NotNilInterface(name string, v any) (err error) {
	if v == nil {
		return fmt.Errorf("%s: %w", name, errors.ErrNoValue)
	}

	return nil
}

// Positive returns an error if v is less than or equal to the zero value of
// type T.  The underlying error of err is [errors.ErrNotPositive].
//
// NOTE:  NaN is also considered negative, since [cmp.Compare] sorts it below
// -Infinity.
func Positive[T cmp.Ordered](name string, v T) (err error) {
	var zero T
	if cmp.Compare(v, zero) <= 0 {
		return fmt.Errorf("%s: %w: %v", name, errors.ErrNotPositive, v)
	}

	return nil
}
