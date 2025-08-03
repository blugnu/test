package test

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/matchers/typecheck"
	"github.com/blugnu/test/opt"
)

// ExpectType tests that a value is of an expected type.  If the test passes,
// the value is returned as that type, with true. If the test fails the zero
// value of the specified type is returned, with false.
//
// If the value being of the expected type is essential to the test, consider
// using the [RequireType] function instead, which will return the value or
// fail the test immediately, avoiding the need to check the ok value.
//
// If a test does not use the returned value, consider using the [BeOfType]
// matcher instead, to avoid lint warnings about unused return values.
func ExpectType[T any](got any, opts ...any) (T, bool) {
	GetT().Helper()

	result, ok := got.(T)
	if ok {
		return result, true
	}

	var (
		zero         T
		expectedType = fmt.Sprintf("%T", zero)
	)

	// if we could not determine the type of the expected value using the
	// zero value of the type, we can use a dummy function and reflect
	// the type of the first argument to that function.
	//
	// Q: why not just use the dummy func technique every time?
	// A: because it is more expensive than using the zero value, and
	//    using the zero value provides a more precise (package-qualified)
	//    type name, when successful
	if expectedType == "<nil>" {
		fn := func(T) { /* NO-OP */ }
		fn(zero) // ensures that the dummy function is 'covered' by tests

		expectedType = reflect.TypeOf(fn).In(0).Name()
	}

	gotType := fmt.Sprintf("%T", got)

	opts = append(opts,
		opt.FailureReport(func(...any) []string {
			return []string{
				"expected type: " + expectedType,
				"got          : " + gotType,
			}
		}),
	)

	Expect(gotType, opts...).To(Equal(expectedType), opts...)

	return zero, false
}

// RequireType tests that a value is of an expected type. If the test passes,
// the value is returned as that type otherwise the test fails immediately
// without evaluating any further expectations.
//
// If a test does not use the returned value, consider using the [BeOfType]
// matcher instead, to avoid lint warnings about unused return values.
func RequireType[T any](got any, opts ...any) T {
	GetT().Helper()

	z, _ := ExpectType[T](got, append(opts, opt.Required())...)

	return z
}

// BeOfType returns a matcher that checks if the value is of the expected type.
//
// If a test wishes to perform further tests on the value that rely on having a
// value of the expected type, consider using the [ExpectType] or [RequireType]
// functions instead, which will return the value as the expected
// type, or fail the test immediately.
func BeOfType[T any]() typecheck.Matcher[T] {
	return typecheck.Matcher[T]{}
}
