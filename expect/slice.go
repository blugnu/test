package expect

import (
	"slices"

	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
)

type SliceTest[T any] struct {
	got   []T
	eopts []any
	mopts []any
}

// Slice returns an expectation for the provided slice. The expectation provides
// assertions specific to slices. Any options provided are used for all assertions
// made using the returned expectation.
//
// # Supported Options
//
//	string                      // a name for the expectation subject; used in
//	opt.Name(string)            // failure reports.
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
//
//	opt.Required()              // if this option is provided, any assertion
//	                            // that fails using the returned expectation
//	                            // will halt further execution of the test.
//	opt.IsRequired(true)        //
func Slice[T any](got []T, opts ...any) SliceTest[T] {
	st := SliceTest[T]{got: got}
	st.eopts, st.mopts = internal.SeparateOptions(opts)
	return st
}

// Contains asserts that the slice under test contains all of the specified
// values. The items may appear in any order in the slice under test and do
// not need to be contiguous.
//
// If any of the specified values are duplicates they must appear in the
// slice under test at least that many times.
func (st SliceTest[T]) Contains(values ...T) {
	test.T().Helper()
	test.Expect(st.got, st.eopts...).To(test.ContainItems(values), st.mopts...)
}

// ContainsAny asserts that the slice under test contains at least one of the
// specified values.
func (st SliceTest[T]) ContainsAny(values ...T) {
	test.T().Helper()
	opts := append(slices.Clone(st.mopts), opt.AnyItem)
	test.Expect(st.got, st.eopts...).To(test.ContainItems(values), opts...)
}

// ContainsSlice asserts that the slice under test contains the specified
// sub-slice in the same contiguous order.
func (st SliceTest[T]) ContainsSlice(values ...T) {
	test.T().Helper()
	test.Expect(st.got, st.eopts...).To(test.ContainSlice(values), st.mopts...)
}

// Equals asserts that the slice under test is equal to the specified values.
// The slices must be of the same length and contain the same values in the
// same order.
func (st SliceTest[T]) Equals(other []T) {
	test.T().Helper()
	test.Expect(st.got, st.eopts...).To(test.EqualSlice(other), st.mopts...)
}

// HasLength asserts that the slice under test has the specified length.
func (st SliceTest[T]) HasLength(length int) {
	test.T().Helper()
	test.Expect(len(st.got), st.eopts...).To(test.Equal(length), st.mopts...)
}

// IsEmpty asserts that the slice under test is empty.
func (st SliceTest[T]) IsEmpty() {
	test.T().Helper()
	test.Expect(st.got, st.eopts...).Should(test.BeEmpty(), st.mopts...)
}

// IsNotEmpty asserts that the slice under test is not empty.
func (st SliceTest[T]) IsNotEmpty() {
	test.T().Helper()
	test.Expect(st.got, st.eopts...).ShouldNot(test.BeEmpty(), st.mopts...)
}
