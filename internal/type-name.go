package internal

import (
	"fmt"
	"reflect"
)

// TypeName returns the name of the type of the provided value.
//
// If the type cannot be determined (for example, if the value is nil),
// the name of the type parameter T is returned instead.
//
// Note that for unnamed types (for example, inline interfaces), "<anon>"
// is returned as the type name.
func TypeName[T any](v ...T) string {
	var n string
	var t T

	switch len(v) {
	case 0:
		n = fmt.Sprintf("%T", t)

	default:
		n = fmt.Sprintf("%T", v[0])
	}

	// if we could not determine the type of the expected value using the
	// zero value of the type, we can use a dummy function and reflect
	// the type of the first argument to that function.
	//
	// Q: why not just use the dummy func technique every time?
	// A: because it is more expensive than using the zero value, and
	//    using the zero value provides a more precise (package-qualified)
	//    type name, when possible
	if n == "<nil>" {
		fn := func(T) { /* NO-OP */ }
		n = reflect.TypeOf(fn).In(0).Name()

		// call the function to ensure it is considered "covered" by tests
		fn(t)
	}

	if n == "" {
		n = "<anon>"
	}

	return n
}
