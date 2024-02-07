package test

import (
	"fmt"
	"reflect"
	"testing"
)

// provides methods for testing values of any type.
type AnyTest[T any] struct {
	testable[T]
}

// creates a new test for a value of any type.  Options of the following
// types are accepted:
//
//	string                   // a name for the test;
//	                         // if not specified, "got" will be used
//
//	Format                   // a format verb for formatting the value being tested;
//	                         // if not specified, FormatDefault will be used
//
//	func(T) string           // a function for formatting the value being tested;
//	                         // if not specified, fmt.Sprintf() will be used
//
// Example:
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s, "result").DeepEquals("something")
//	}
func That[T any](t *testing.T, got T, opts ...any) AnyTest[T] {
	t.Helper()

	n := "got"
	f := FormatDefault
	ffn := *new(func(T) string)
	checkOptTypes(t, optTypes(n, f, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	getOpt(&ffn, append(opts, func(v T) string { return fmt.Sprintf(string(f), v) })...)

	return AnyTest[T]{newTestable(t, got, n, ffn)}
}

// compares the value being tested with a wanted value of the same type using
// reflect.DeepEqual().
//
// The test fails if the two values are not equal.
//
// If you need to compare values using something other than reflect.DeepEquals(),
// use the Equals() method and provide a comparison function.
func (at AnyTest[T]) DeepEquals(wanted T) {
	at.Helper()

	if !reflect.DeepEqual(at.got, wanted) {
		at.Run("deep equals", func(t *testing.T) {
			t.Helper()
			t.Errorf("\nwanted: %s\ngot   : %s", at.ffn(wanted), at.ffn(at.got))
		})
	}
}

// compares the value being tested with a wanted value of the same type.  The following
// options are accepted:
//
//	string                   // a name for the test;
//	                         // if not specified, "equals" will be used
//
//	func(got, wanted T) bool // a comparison function that returns true if the two values are considered equal;
//	                         // if not specified, reflect.DeepEqual() will be used.
//
// The test fails if the two values do not satisfy the comparison function.
//
// Example:
//
//	func TestSomething(t *testing.T) {
//		// ACT
//		s := DoSomething()
//
//		// ASSERT
//		test.That(t, s, "result").Equals("something:", "has prefix", func (got, wanted string) bool {
//			return strings.HasPrefix(got, wanted)
//		})
//	}
func (at AnyTest[T]) Equals(wanted T, opts ...any) {
	at.Helper()

	cmp := func(got, wanted T) bool { return reflect.DeepEqual(got, wanted) }
	cfn := "reflect.DeepEqual"
	n := "equals"
	checkOptTypes(at.T, optTypes(n, cmp), opts...)
	getOpt(&n, opts...)
	if getOpt(&cmp, opts...) {
		cfn = "<comparison func>"
	}

	if !cmp(at.got, wanted) {
		at.Run(n, func(t *testing.T) {
			t.Helper()
			t.Errorf("\nwanted: %s\ngot   : %s\nmethod: %s", at.ffn(wanted), at.ffn(at.got), cfn)
		})
	}
}

// fails the test if the value being tested is not nil.
//
// If the value being tested is of an underlying type that is not nilable the test will
// fail with ErrNotNilable (written to the test failure report), irrespective of the
// value itself.
func (at AnyTest[T]) IsNil() {
	at.Helper()
	at.run(func(t *testing.T) {
		t.Helper()
		IsNil(t, at.got)
	})
}

// fails the test if the value being tested is nil.
//
// If the value being tested is of an underlying type that is not nilable the test will
// fail with ErrNotNilable (written to the test failure report), irrespective of the
// value itself.
func (at AnyTest[T]) IsNotNil() {
	at.Helper()
	at.run(func(t *testing.T) {
		t.Helper()
		IsNotNil(t, at.got)
	})
}
