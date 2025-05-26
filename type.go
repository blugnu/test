package test

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

// ExpectType tests that a value is of an expected type.  If the test passes,
// the value is returned as that type, with true. If the test fails the zero
// value of the specified type is returned, with false.
func ExpectType[T any](got any, opts ...any) (T, bool) {
	GetT().Helper()

	z := *new(T)
	gotType := reflect.TypeOf(got)
	expectedType := reflect.TypeOf(z)

	if fmt.Sprintf("%s", expectedType) == "%!s(<nil>)" { //nolint:gosimple // .String() on a nil is not a great idea!
		invalidTest("ExpectType: cannot be used to test for interfaces")
		return z, false
	}

	Expect(gotType, opts...).To(Equal(expectedType), opt.FailureReport(func(...any) []string {
		return []string{
			fmt.Sprintf("expected: %s", expectedType),
			fmt.Sprintf("got     : %s", gotType),
		}
	}))

	if gotType == expectedType {
		return got.(T), true
	}

	return z, false
}
