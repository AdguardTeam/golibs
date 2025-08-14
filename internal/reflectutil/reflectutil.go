// Package reflectutil provides a safe and limited interface to some of the
// features of package reflect.
//
// TODO(a.garipov):  Consider exporting.
package reflectutil

import "reflect"

// FormatVerb returns a more fitting formatting verb depending on the type and
// the value.  The current formatting exceptions are:
//   - Nilable values that are nil are printed with the "%#v" verb.
//   - Strings are printed with the "%q" verb.
func FormatVerb(val any) (verb string) {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice:
		if v.IsNil() {
			return "%#v"
		}
	case reflect.String:
		return "%q"
	}

	return "%+v"
}

// IsInterface returns true if T is an interface type.
func IsInterface[T any]() (ok bool) {
	return reflect.TypeFor[*T]().Elem().Kind() == reflect.Interface
}
