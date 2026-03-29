package internal

import "reflect"

// IsNil returns true if the supplied value is nil.  It returns false if the supplied value is
// not nil or if the supplied value is of a type that does not support a nil value.
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() { //nolint:exhaustive // only concerned with nilable types
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
