package test

import (
	"fmt"
	"testing"
)

// provides methods for testing a value of a comparable type.
type ValueTest[T comparable] struct {
	testable[T]
}

// creates a testable for a value of a type that satisfies comparable.  Options of the
// following types are accepted:
//
//	string              // a name for the value being tested;
//	                    // if not specified "value" is used
//
//	Format              // a format verb for the value being tested;
//	                    // if not specified, FormatDefault is used.
//	                    // This option is ignored if a formatting function is specified
//
//	func(T) string      // a formatting function that returns a string representation
//	                    // of a value of the type being tested; if not specified, values
//	                    // are formatted using the configured Format verb.
//
// If more than one option of any of the above types is specified then only the first
// is applied; additional values of that option type are ignored.
//
// To create a testable for a value of some type that is not comparable, use test.That()
// or a testable factory for the specific type concerned.
func Value[T comparable](t *testing.T, got T, opts ...any) ValueTest[T] {
	t.Helper()

	n := "value"
	f := FormatDefault
	ffn := *new(func(T) string)
	checkOptTypes(t, optTypes(n, f, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	getOpt(&ffn, append(opts, func(v T) string { return fmt.Sprintf(string(f), v) })...)

	return ValueTest[T]{newTestable(t, got, n, f, ffn)}
}

// fails the test if the value being tested equals some specified value. Options
// of the following types are accepted:
//
//	string                     // a name for the test; if not specified, "does not equal" is used
//
//	Equality                   // determines the method used to compare values; if not
//	                           // specified, test.ShallowEquality is used and values are compared
//	                           // using the == operator.  test.DeepEquality may be specified to
//	                           // compare values using reflect.DeepEqual.
//
//	func(got, wanted T) bool   // a custom comparison function; if specified, the Equality option
//	                           // is ignored and the custom function is used to compare values.
//
// The test will fail if the two values are equal.
//
// # Example of shallow equality test (default):
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s).DoesNotEqual("something")
//	}
//
// # Example using a custom comparison function:
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s).DoesNotEqual("something:", "does not begin with", func(got, wanted string) bool {
//			return !strings.HasPrefix(got, wanted)
//		})
//	}
func (v ValueTest[T]) DoesNotEqual(wanted T, opts ...any) {
	v.Helper()

	n := "does not equal"
	eq := ShallowEquality
	cfn := *new(func(T, T) bool)
	checkOptTypes(v.T, optTypes(n, eq, cfn), opts...)

	getOpt(&n, opts...)
	explainMethod := getOpt(&eq, opts...)
	if getOpt(&cfn, opts...) {
		explainMethod = true
		eq = customEquality
	}

	if compareFunc(eq, cfn)(v.got, wanted) {
		v.Run(n, func(t *testing.T) {
			t.Helper()
			report := fmt.Sprintf("\nwanted: not: %s\ngot        : %s", v.ffn(wanted), v.ffn(v.got))
			if explainMethod {
				report += fmt.Sprintf("\nmethod     : %s", eq)
			}
			t.Error(report)
		})
	}
}

// fails the test if the value being tested is not equal to some specified value.
// Options of the following types are accepted:
//
//	string                     // a name for the test; if not specified, "equals" is used
//
//	Equality                   // determines the method used to compare values; if not
//	                           // specified, test.ShallowEquality is used and values are compared
//	                           // using the == operator.  test.DeepEquality may be specified to
//	                           // compare values using reflect.DeepEqual.
//
//	func(got, wanted T) bool   // a custom comparison function; if specified, the Equality option
//	                           // is ignored and the custom function is used to compare values.
//
// The test will fail if the two values are not equal.
//
// # Example of shallow equality test (default):
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s).Equals("something")
//	}
//
// # Example using a custom comparison function:
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s).Equals("something:", "begins with", strings.HasPrefix)
//	}
func (v ValueTest[T]) Equals(wanted T, opts ...any) {
	v.Helper()

	n := "equals"
	eq := ShallowEquality
	cfn := *new(func(T, T) bool)
	checkOptTypes(v.T, optTypes(n, eq, cfn), opts...)

	getOpt(&n, opts...)
	explainMethod := getOpt(&eq, opts...)
	if getOpt(&cfn, opts...) {
		explainMethod = true
		eq = customEquality
	}

	if !compareFunc(eq, cfn)(v.got, wanted) {
		v.Run(n, func(t *testing.T) {
			t.Helper()
			report := fmt.Sprintf("\nwanted: %s\ngot   : %s", v.ffn(wanted), v.ffn(v.got))
			if explainMethod {
				report += fmt.Sprintf("\nmethod: %s", eq)
			}
			t.Error(report)
		})
	}
}

// tests that the value being tested is nil; the test will fail if
// the value is not nil.
//
// If the value being tested is of an underlying type that is not nilable the test will
// fail with ErrNotNilable (written to the test failure report), irrespective of the
// value itself.
func (vt ValueTest[T]) IsNil() {
	vt.Helper()
	vt.run(func(t *testing.T) {
		t.Helper()
		IsNil(t, vt.got)
	})
}

// tests that the value being tested is not nil; the test will fail if
// the value is nil.
//
// If the value being tested is of an underlying type that is not nilable the test will
// fail with ErrNotNilable (written to the test failure report), irrespective of the
// value itself.
//
// NOTE: when an interface value is nil it is not possible to determine
// the expected type of the value.  In these cases the test failure
// report will not be very informative.
//
// For example, if the value is a nil slice, the test failure report
// might be similar to:
//
//	wanted: []uint8
//	got   : nil
//
// Where-as for an error, any (or other interface type) etc, the failure
// report will be similar to:
//
//	wanted: not nil
//	got   : nil
//
// Although it may be sufficient to provide a meaningful name or description
// for the value being tests, if you want or need the type indicated in the
// test, use a NotEqual(t, got, nil) test instead.
func (vt ValueTest[T]) IsNotNil() {
	vt.Helper()
	vt.run(func(t *testing.T) {
		t.Helper()
		IsNotNil(t, vt.got)
	})
}
