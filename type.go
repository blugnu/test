package test

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

// ExpectType tests that a value is of an expected type.  If the test passes,
// the value is returned as that type, with true. If the test fails the zero
// value of the specified type is returned, with false.
func ExpectType[T any](got any, opts ...any) (T, bool) {
	GetT().Helper()

	z := *new(T)
	gotType := reflect.TypeOf(got)
	expectedType := reflect.TypeOf(z)

	if fmt.Sprintf("%s", expectedType) == "%!s(<nil>)" {
		test.Invalid("ExpectType: cannot be used to test for interfaces")
		return z, false
	}

	Expect(gotType, opts...).To(Equal(expectedType), opt.FailureReport(func(...any) []string {
		return []string{
			fmt.Sprintf("expected: %s", expectedType),
			fmt.Sprintf("got     : %s", gotType),
		}
	}))

	if gotType == expectedType {
		got, ok := got.(T)
		return got, ok
	}

	return z, false
}

func RequireType[T any](got any, opts ...any) T {
	GetT().Helper()

	z, _ := ExpectType[T](got, append(opts, opt.Required())...)
	return z
}
