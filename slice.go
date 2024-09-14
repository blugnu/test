package test

import (
	"fmt"
	"testing"
)

// provides methods for testing a slice.
type SliceTest[T comparable] struct {
	testable[[]T]
}

// creates a testable for a slice of a type that satisfies comparable. Options
// of the following types are accepted:
//
//	string              // a name for the slice being tested; if not specified, "slice" is used
//
//	Format              // a format verb for formatting the slice being tested; if not
//	                    // specified, FormatDecl is used
//
//	func([]T) string    // a formatting function that returns a string representation
//	                    // of a slice of the type being tested; if not specified, values
//	                    // are formatted using the configured Format verb.
//
// To create a testable for a slice of some type that is not comparable, use test.That()
// and use methods with custom comparison functions, or a testable factory for the specific
// type concerned, if available, such as test.Bytes().
func Slice[T comparable](t *testing.T, got []T, opts ...any) SliceTest[T] {
	t.Helper()

	n := "slice"
	f := FormatDecl
	ffn := *new(func([]T) string)
	checkOptTypes(t, optTypes(n, f, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	getOpt(&ffn, append(opts, func(v []T) string { return fmt.Sprintf(string(f), any(v)) })...)

	return SliceTest[T]{newTestable(t, got, n, ffn)}
}

// fails the test if the slice being tested is not equal to the wanted
// slice.
//
// To be equal, slices must be of the same length and each element must
// be equal to the corresponding element in the other slice.
//
// Options of the following types are accepted:
//
//	Equality                     // determines the method used to compare elements; if not
//	                             // specified, test.ShallowEquality is used and elements are compared
//	                             // using the == operator.  test.DeepEquality may be specified to
//	                             // compare elements using reflect.DeepEqual.
//
//	func(got, wanted T) bool     // a custom comparison function; if specified, the Equality option
//	                             // is ignored and the custom function is used to compare elements.
//
// # Example default comparison (test.ShallowEquality using ==):
//
//	test.Slice(t, result.Bytes(), "result buffer").Equals(expected)
//	test.Slice(t, result.Bytes(), "result buffer").Equals(expected, test.ShallowEquality)
//
// # Example test.DeepEquality:
//
//	test.Slice(t, result, "customers").Equals(expected, test.DeepEquality)
//
// # Example using custom comparison function:
//
//	test.Slice(t, result, "customers").Equals(expected, func(got, wanted Customer) bool { return got.Id == wanted.Id })
func (st SliceTest[T]) Equals(wanted []T, opts ...any) {
	st.Helper()

	fnopts := []any{fmt.Sprintf("%s/equals", st.name), st.ffn}

	eq := ShallowEquality
	cfn := *new(func(T, T) bool)
	checkOptTypes(st.T, optTypes(eq, cfn), opts...)
	getOpt(&eq, opts...)
	fnopts = append(fnopts, eq)

	if getOpt(&cfn, opts...) {
		fnopts = append(fnopts, cfn)
	}

	SlicesEqual(st.T, st.got, wanted, fnopts...)
}

// fails the test if the slice being tested is not empty.
//
// Example:
//
//	test.Slice(t, "result buffer", result.Bytes()).IsEmpty()
func (st SliceTest[T]) IsEmpty() {
	st.Helper()
	st.run(func(t *testing.T) {
		t.Helper()
		if st.got == nil || len(st.got) != 0 {
			t.Errorf("\nwanted: empty %[1]T\ngot   : %#[1]v", st.got)
		}
	})
}

// fails the test if the slice being tested is empty.
//
// Example:
//
//	test.Slice(t, "result buffer", result.Bytes()).IsNotEmpty()
func (st SliceTest[T]) IsNotEmpty() {
	st.Helper()
	st.run(func(t *testing.T) {
		t.Helper()
		if len(st.got) == 0 {
			t.Errorf("\nwanted: non-empty %T\ngot   : empty", st.got)
		}
	})
}

// returns true if the slice contains an element that is considered equal
// to some specified element, using a specified comparison function.
func sliceContains[T comparable](s []T, el T, cmp func(T, T) bool) bool {
	if cmp == nil {
		cmp = func(a, b T) bool { return a == b }
	}
	for _, v := range s {
		if cmp(v, el) {
			return true
		}
	}
	return false
}

// returns true if two slices are of the same length and each element is equal
// to the corresponding element in the other slice, using a specified comparison
// function.
func slicesEqual[T comparable](a, b []T, cmp func(T, T) bool) bool {
	if cmp == nil {
		cmp = func(a, b T) bool { return a == b }
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !cmp(a[i], b[i]) {
			return false
		}
	}
	return true
}

// fails the test if the specified slices are not equal.
func SlicesEqual[T comparable](t *testing.T, got, wanted []T, opts ...any) {
	t.Helper()

	n := "equals"
	eq := ShallowEquality
	f := FormatDefault
	ffn := *new(func([]T) string)
	cfn := *new(func(T, T) bool)
	checkOptTypes(t, optTypes(n, eq, f, cfn, ffn), opts...)
	getOpt(&n, opts...)
	explainMethod := getOpt(&eq, opts...)
	getOpt(&f, opts...)
	getOpt(&ffn, append(opts, func(v []T) string { return fmt.Sprintf(string(f), v) })...)
	if getOpt(&cfn, opts...) {
		explainMethod = true
		eq = customEquality
	}

	// get the desired comparison function
	cmp := compareFunc(eq, cfn)

	if !slicesEqual(got, wanted, cmp) {
		t.Run(n, func(t *testing.T) {
			t.Helper()
			report := fmt.Sprintf("\nwanted: %s\ngot   : %s", ffn(wanted), ffn(got))
			if explainMethod {
				report += fmt.Sprintf("\nmethod: %s", eq)
			}
			t.Error(report)
		})
	}
}
