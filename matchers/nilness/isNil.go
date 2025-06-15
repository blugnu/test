package nilness

import "reflect"

// isNil returns two boolean values.  The first indicates whether
// the supplied value is nil.  The second indicates whether the
// supplied value is of a type that supports a nil value.
//
// i.e. if this function returns (false, false), then the value
// is not nil and does not support a nil value.  If it returns
// (false, true), then the value type supports a nil value but
// the value itself is not nil.
func isNil(v any) (bool, bool) {
	const Nilable = true
	const NotNilable = false

	if v == nil {
		return true, Nilable
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return reflect.ValueOf(v).IsNil(), Nilable
	}

	return false, NotNilable
}
