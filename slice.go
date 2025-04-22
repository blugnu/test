package test

import (
	"fmt"
	"strings"
)

type AnySliceMatcher[T any] struct {
	Expecting[[]T]
	match func([]T, []T, func(T, T) bool) bool
	// condition sliceMatchCondition
}

func (m AnySliceMatcher[T]) Match(actual []T) bool {
	return m.match(actual, m.Expecting.value, nil)
}

func (m AnySliceMatcher[T]) Format(a any) string {
	s := a.([]T)

	result := "\n"
	for i, s := range s {
		result = result + fmt.Sprintf("%3d: %v\n", i, s)
	}
	return result[:len(result)-1]
}

type SliceMatcher[T comparable] struct {
	Expecting[[]T]
	match func([]T, []T, func(T, T) bool) bool
}

func (m SliceMatcher[T]) Match(actual []T, opts ...any) bool {
	if hasOpt(opts, ExactOrder(false)) {
		panic("not implemented: SliceMatcher ExactOrder(false)")
	}
	return m.match(actual, m.Expecting.value, nil)
}

func (m SliceMatcher[T]) Format(a any) string {
	s := a.([]T)

	result := "\n"
	for i, s := range s {
		result = result + fmt.Sprintf("%3d: %v\n", i, s)
	}
	return result[:len(result)-1]
}

func ExpectEmpty[T any](actual []T, opts ...any) {
	GetT().Helper()
	Expect(len(actual), opts...).To(Equal(0))
}

func ContainStrings[T ~string](e []T) SliceMatcher[T] {
	return SliceMatcher[T]{Expecting[[]T]{e}, func(act, exp []T, _ func(T, T) bool) bool {
		if len(act) == 0 && len(exp) == 0 {
			return true
		}
		if len(act) < len(exp) {
			return false
		}
		if len(exp) == 0 {
			return false
		}
		for i, a := range act {
			if strings.Contains(string(a), string(exp[0])) {
				if i+len(exp) > len(act) {
					return false
				}
				for j, b := range exp {
					if !strings.Contains(string(act[i+j]), string(b)) {
						return false
					}
				}
				return true
			}
		}
		return false
	}}
}

func Contain[T comparable](e T) SliceMatcher[T] {
	return SliceMatcher[T]{Expecting[[]T]{[]T{e}}, func(act, exp []T, f func(T, T) bool) bool {
		if len(act) == 0 {
			return false
		}
		for _, t := range act {
			if f(t, exp[0]) {
				return true
			}
		}
		return false
	}}
}

func ContainSlice[T comparable](e []T) SliceMatcher[T] {
	return SliceMatcher[T]{Expecting[[]T]{e}, sliceContainsS[T]}
}

// EqualSlice compares the actual slice with an expected slice and fails
// the test if they are not equal.
//
// By default, the order of elements in each slice is significant.  That is,
// the nth each slice must be equal. If the order of elements is not
// significant, use the ExactOrder option to specify that the order of elements
// is not significant, e.g.:
//
//	got := []int{1, 2, 3}
//	expected := []int{3, 2, 1}
//	Expect(got).To(EqualSlice(expected)) // will fail
//	Expect(got).To(EqualSlice(expected, ExactOrder(false))) // will pass
func EqualSlice[T comparable](e []T, opts ...any) SliceMatcher[T] {
	fn := slicesEqual[T]
	if hasOpt(opts, ExactOrder(false)) {
		fn = slicesEquivalentNoDiff[T]
	}
	return SliceMatcher[T]{Expecting[[]T]{e}, fn}
}

// // provides methods for testing a slice.
// type SliceTest[T comparable] struct {
// 	testable[[]T]
// }

// // creates a testable for a slice of a type that satisfies comparable. Options
// // of the following types are accepted:
// //
// //	string              // a name for the slice being tested; if not specified, "slice" is used
// //
// //	Format              // a format verb for formatting the slice being tested; if not
// //	                    // specified, FormatDecl is used
// //
// //	func([]T) string    // a formatting function that returns a string representation
// //	                    // of a slice of the type being tested; if not specified, values
// //	                    // are formatted using the configured Format verb.
// //
// // To create a testable for a slice of some type that is not comparable, use test.That()
// // and use methods with custom comparison functions, or a testable factory for the specific
// // type concerned, if available, such as test.Bytes().
// func Slice[T comparable](t *testing.T, got []T, opts ...any) SliceTest[T] {
// 	t.Helper()

// 	n := "slice"
// 	f := FormatDecl
// 	ffn := *new(func([]T) string)
// 	checkOptTypes(t, optTypes(n, f, ffn), opts...)
// 	getOpt(&n, opts...)
// 	getOpt(&f, opts...)
// 	getOpt(&ffn, append(opts, func(v []T) string { return fmt.Sprintf(string(f), any(v)) })...)

// 	return SliceTest[T]{newTestable(t, got, n, ffn)}
// }

// // fails the test if the slice being tested is not equal to the wanted
// // slice.
// //
// // To be equal, slices must be of the same length and each element must
// // be equal to the corresponding element in the other slice.
// //
// // Options of the following types are accepted:
// //
// //	Equality                     // determines the method used to compare elements; if not
// //	                             // specified, test.ShallowEquality is used and elements are compared
// //	                             // using the == operator.  test.DeepEquality may be specified to
// //	                             // compare elements using reflect.DeepEqual.
// //
// //	func(got, wanted T) bool     // a custom comparison function; if specified, the Equality option
// //	                             // is ignored and the custom function is used to compare elements.
// //
// // # Example default comparison (test.ShallowEquality using ==):
// //
// //	test.Slice(t, result.Bytes(), "result buffer").Equals(expected)
// //	test.Slice(t, result.Bytes(), "result buffer").Equals(expected, test.ShallowEquality)
// //
// // # Example test.DeepEquality:
// //
// //	test.Slice(t, result, "customers").Equals(expected, test.DeepEquality)
// //
// // # Example using custom comparison function:
// //
// //	test.Slice(t, result, "customers").Equals(expected, func(got, wanted Customer) bool { return got.Id == wanted.Id })
// func (st SliceTest[T]) Equals(wanted []T, opts ...any) {
// 	st.Helper()

// 	fnopts := []any{fmt.Sprintf("%s/equals", st.name), st.ffn}

// 	eq := ShallowEquality
// 	cfn := *new(func(T, T) bool)
// 	checkOptTypes(st.T, optTypes(eq, cfn), opts...)
// 	getOpt(&eq, opts...)
// 	fnopts = append(fnopts, eq)

// 	if getOpt(&cfn, opts...) {
// 		fnopts = append(fnopts, cfn)
// 	}

// 	SlicesEqual(st.T, st.got, wanted, fnopts...)
// }

// // fails the test if the slice being tested is not equivalent to the wanted
// // slice.
// //
// // To be equivalent, slices must be of the same length and each element occuring
// // in the slice being tested must occur in the wanted slice, and vice versa.
// // The order of elements is not considered.
// //
// // Options of the following types are accepted:
// //
// //	Equality                     // determines the method used to compare elements; if not
// //	                             // specified, test.ShallowEquality is used and elements are compared
// //	                             // using the == operator.  test.DeepEquality may be specified to
// //	                             // compare elements using reflect.DeepEqual.
// //
// //	func(got, wanted T) bool     // a custom comparison function; if specified, the Equality option
// //	                             // is ignored and the custom function is used to compare elements.
// func (st SliceTest[T]) EquivalentTo(wanted []T, opts ...any) {
// 	st.Helper()

// 	fnopts := []any{fmt.Sprintf("%s/equals", st.name), st.ffn}

// 	eq := ShallowEquality
// 	cfn := *new(func(T, T) bool)
// 	checkOptTypes(st.T, optTypes(eq, cfn), opts...)
// 	getOpt(&eq, opts...)
// 	fnopts = append(fnopts, eq)

// 	if getOpt(&cfn, opts...) {
// 		fnopts = append(fnopts, cfn)
// 	}

// 	SlicesEquivalent(st.T, st.got, wanted, fnopts...)
// }

// // fails the test if the slice being tested is not of the expected length.
// func (st SliceTest[T]) HasLength(wanted int) {
// 	st.Helper()
// 	st.run(func(t *testing.T) {
// 		t.Helper()
// 		if len(st.got) != wanted {
// 			t.Errorf("\nwanted: length %d\ngot   : length %d", wanted, len(st.got))
// 		}
// 	})
// }

// // fails the test if the slice being tested is not empty.
// //
// // Example:
// //
// //	test.Slice(t, "result buffer", result.Bytes()).IsEmpty()
// func (st SliceTest[T]) IsEmpty() {
// 	st.Helper()
// 	st.run(func(t *testing.T) {
// 		t.Helper()
// 		if st.got == nil || len(st.got) != 0 {
// 			t.Errorf("\nwanted: empty %[1]T\ngot   : %#[1]v", st.got)
// 		}
// 	})
// }

// // fails the test if the slice being tested is empty.
// //
// // Example:
// //
// //	test.Slice(t, "result buffer", result.Bytes()).IsNotEmpty()
// func (st SliceTest[T]) IsNotEmpty() {
// 	st.Helper()
// 	st.run(func(t *testing.T) {
// 		t.Helper()
// 		if len(st.got) == 0 {
// 			t.Errorf("\nwanted: non-empty %T\ngot   : empty", st.got)
// 		}
// 	})
// }

// returns true if the slice contains an element that is considered equal
// to some specified element, using a specified comparison function.
func sliceContainsE[T comparable](s []T, el T, cmp func(T, T) bool) bool {
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

// returns true if the slice contains an element that is considered equal
// to some specified element, using a specified comparison function.
func sliceContainsS[T comparable](s []T, other []T, cmp func(T, T) bool) bool {
	// special case: if both slices are empty they are considered equal
	if len(s) == 0 && len(other) == 0 {
		return true
	}
	// special case: the slice s has fewer elements than the subslice then it
	// cannot possibly "contain" that subslice
	if len(s) < len(other) {
		return false
	}

	if cmp == nil {
		cmp = func(a, b T) bool { return a == b }
	}

	var first int = -1
	for i, is := range s {
		if cmp(is, other[0]) {
			first = i
		}
	}
	if first == -1 {
		return false
	}

	for i, o := range other {
		if first+i >= len(s) || !cmp(s[first+i], o) {
			return false
		}
	}
	return true
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

func slicesEquivalentNoDiff[T comparable](a, b []T, cmp func(T, T) bool) bool {
	_, ok := slicesEquivalent(a, b, cmp, false)
	return ok
}

// returns true if two slices are of the same length and each element in one slice
// also occurs in the other.  The order of elements is not considered.  Equality of
// items is determined using a specified comparison function.
func slicesEquivalent[T comparable](a, b []T, cmp func(T, T) bool, wantDiff bool) (map[T][2]int, bool) {
	if !wantDiff && len(a) != len(b) {
		return nil, false
	}
	if cmp == nil {
		cmp = func(a, b T) bool { return a == b }
	}

	c := make([]T, len(b))
	for i, v := range b {
		c[i] = v
	}

	diff := make(map[T][2]int, len(a))
	for _, v := range a {
		refs := diff[v]
		refs[0]++
		for i, w := range c {
			if cmp(v, w) {
				refs[1]++
				// remove the item from b
				c = append(c[:i], c[i+1:]...)
				break
			}
		}
		diff[v] = refs
	}
	for _, v := range c {
		refs := diff[v]
		refs[1]++
		diff[v] = refs
	}
	return diff, len(c) == 0
}

// // fails the test if the specified slices are not equal.
// func SlicesEqual[T comparable](t *testing.T, got, wanted []T, opts ...any) {
// 	t.Helper()

// 	n := "equals"
// 	eq := ShallowEquality
// 	f := FormatDefault
// 	ffn := *new(func([]T) string)
// 	cfn := *new(func(T, T) bool)
// 	checkOptTypes(t, optTypes(n, eq, f, cfn, ffn), opts...)
// 	getOpt(&n, opts...)
// 	explainMethod := getOpt(&eq, opts...)
// 	getOpt(&f, opts...)
// 	getOpt(&ffn, append(opts, func(v []T) string { return fmt.Sprintf(string(f), v) })...)
// 	if getOpt(&cfn, opts...) {
// 		explainMethod = true
// 		eq = customEquality
// 	}

// 	// get the desired comparison function
// 	cmp := compareFunc(eq, cfn)

// 	if !slicesEqual(got, wanted, cmp) {
// 		t.Run(n, func(t *testing.T) {
// 			t.Helper()
// 			report := fmt.Sprintf("\nwanted: %s\ngot   : %s", ffn(wanted), ffn(got))
// 			if explainMethod {
// 				report += fmt.Sprintf("\nmethod: %s", eq)
// 			}
// 			t.Error(report)
// 		})
// 	}
// }

// // fails the test if the specified slices are not equal.
// func SlicesEquivalent[T comparable](t *testing.T, got, wanted []T, opts ...any) {
// 	t.Helper()

// 	n := "equivalent"
// 	eq := ShallowEquality
// 	f := FormatDefault
// 	ffn := *new(func([]T) string)
// 	cfn := *new(func(T, T) bool)
// 	checkOptTypes(t, optTypes(n, eq, f, cfn, ffn), opts...)
// 	getOpt(&n, opts...)
// 	explainMethod := getOpt(&eq, opts...)
// 	getOpt(&f, opts...)
// 	getOpt(&ffn, append(opts, func(v []T) string { return fmt.Sprintf(string(f), v) })...)
// 	if getOpt(&cfn, opts...) {
// 		explainMethod = true
// 		eq = customEquality
// 	}

// 	// get the desired comparison function
// 	cmp := compareFunc(eq, cfn)

// 	// TODO: diff report...
// 	if _, ok := slicesEquivalent(got, wanted, cmp, false); !ok {
// 		t.Run(n, func(t *testing.T) {
// 			t.Helper()
// 			report := fmt.Sprintf("\nwanted: %s\ngot   : %s", ffn(wanted), ffn(got))
// 			if explainMethod {
// 				report += fmt.Sprintf("\nmethod: %s", eq)
// 			}
// 			t.Error(report)
// 		})
// 	}
// }
