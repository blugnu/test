package test

import (
	"fmt"
	"reflect"
)

// ExpectType tests that a value is of a required type.  If the test passes,
// the value is returned as that type, with true. If the test fails the zero
// value of the specified type is returned, with false.
func ExpectType[T any](got any) (T, bool) {
	GetT().Helper()

	z := *new(T)
	gotType := reflect.TypeOf(got)
	expectedType := reflect.TypeOf(z)

	Expect(gotType).To(Equal(expectedType), CustomOneLineReportFunc(func() string {
		return fmt.Sprintf("expected type %s, got %s", expectedType, gotType)
	}))

	if gotType == expectedType {
		return got.(T), true
	}

	return z, false
}
