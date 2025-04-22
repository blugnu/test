package test

import (
	"fmt"
	"reflect"
	"testing"
)

// this file contains the implementation of the test helper functions
// (test.Equal, test.IsNil etc)
//
// the implementation of the test.Helper test factory may be found in
// the file helper.go

// verifies that a specified value is nil, failing the test if it is not.
//
// A failure report for this test produces output similar to:
//
//	wanted: nil
//	got   : <%v representation of value>
//
// This test has special handling for errors.  If the supplied value is
// a non-nil error then the test failure report will be:
//
//	unexpected error: <%T of value>: <%v of value>
//
// An optional name may be supplied to identify or describe the value
// being tested.
//
// Example:
//
//	err = doSomething()
//	test.IsNil(t, err, "returns nil")
func IsNil(t *testing.T, value any, name ...string) {
	t.Helper()

	if len(name) == 0 {
		name = []string{"is nil"}
	}

	t.Run(name[0], func(t *testing.T) {
		t.Helper()
		switch result := isNil(value).(type) {
		case bool:
			if result {
				return
			}
			if got, ok := value.(error); ok && got != nil {
				t.Errorf("\nunexpected error: %[1]T: %[1]s", got)
				return
			}
			t.Errorf("\nwanted: nil (%T)\ngot   : not nil", value)
		case error:
			t.Errorf("\ntest.IsNil: invalid test: values of type '%T' are not nilable", value)
		}
	})
}

// verifies that a specified value is not nil, failing the test if it is.
//
// If the value being tested does not support a nil value the test will
// fail and produce a report similar to:
//
//	test.IsNotNil: invalid test: values of type '<type>' are not nilable
//
// Otherwise, a failure report for this test produces output similar to:
//
//	wanted: []uint8
//	got   : nil
//
// However, it is not possible to report the concrete type if the value being
// tested is a nil interface.  In these cases, the test failure report will be:
//
//	wanted: not nil
//	got   : nil
//
// This may not be very helpful.  If you want or need the test failure report
// to reliably identify the desired type, use a NotEqual(t, got, nil) test
// instead.
//
// An optional name may be supplied to identify or describe the value
// being tested.
//
// Example:
//
//	test.IsNotNil(t, err)
func IsNotNil(t *testing.T, value any, name ...string) {
	t.Helper()

	if len(name) == 0 {
		name = []string{"is not nil"}
	}

	t.Run(name[0], func(t *testing.T) {
		t.Helper()

		switch result := isNil(value).(type) {
		case bool:
			if !result {
				return
			}
			wanted := fmt.Sprintf("not nil (%T)", value)
			if wanted == "not nil (<nil>)" {
				wanted = "not nil"
			}
			t.Errorf("\nwanted: %s\ngot   : nil", wanted)
		case error:
			t.Errorf("\ninvalid test: values of type '%T' are not nilable", value)
		}
	})
}

// // returns a function that compares two values for compareFunc.  The
// // returned function determined by an optional Equality value and
// // the Equality value that was applied.
// //
// // test.ShallowEquality is the default, returning a function that
// // uses == comparison.
// //
// // test.DeepEquality returns a function that uses reflect.DeepEqual
// func compareFunc[T comparable](opt Equality, fn func(T, T) bool) func(T, T) bool {
// 	switch opt {
// 	case ShallowEquality:
// 		return func(a, b T) bool { return a == b }
// 	case DeepEquality:
// 		return func(a, b T) bool { return reflect.DeepEqual(a, b) }
// 	default:
// 		return fn
// 	}
// }

// isNil returns true if the supplied value is nil or false if
// the value is not nil and is of a type that could be nil.
//
// If the supplied value is not nil and is of a type that does not
// support a nil value, the function will return ErrNotNilable.
func isNil(v any) any {
	if v == nil {
		return true
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		if v := reflect.ValueOf(v); v != reflect.Zero(reflect.TypeOf(v)) && !v.IsNil() {
			return false
		}
		return true
	}
	return ErrNotNilable
}
