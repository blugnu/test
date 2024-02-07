package test

import (
	"fmt"
	"testing"
)

// provides methods for testing a bool value.
type BoolTest struct {
	testable[bool]
}

// returns a new BoolTest with methods for testing a specified
// bool value with a specified name, using a specified testing.T
//
// Additional options of the following types may be specified:
//
//   - string : a name for the value being tested; if not specified "got" is used
//
//   - Format : a format verb for the value being tested; if not specified
//     FormatDefault is used.
//
// If more than one option of any of the above types is specified then only the first
// is applied; additional values of that option type are ignored.
func Bool(t *testing.T, got bool, opts ...any) BoolTest {
	t.Helper()

	n := "got"
	f := FormatDefault
	checkOptTypes(t, optTypes(n, f), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)

	return BoolTest{newTestable(t, got, n, f)}
}

// fails the test if the bool being tested does not equal the specified
// value.
//
// Example:
//
//	test.Bool(t, "is called", isCalled).Equals(tc.callsFn)
func (bt BoolTest) Equals(wanted bool) {
	bt.Helper()
	n := fmt.Sprintf("%s/equals", bt.name)
	if wanted {
		IsTrue(bt.T, bt.got, n)
	} else {
		IsFalse(bt.T, bt.got, n)
	}
}

// fails the test if the bool being tested is not false.
//
// Example:
//
//	test.Bool(t, "is called", isCalled).IsFalse()
func (bt BoolTest) IsFalse() {
	bt.Helper()
	IsFalse(bt.T, bt.got, fmt.Sprintf("%s/is false", bt.name))
}

// fails the test if the bool being tested is not true.
//
// Example:
//
//	test.Bool(t, "is called", isCalled).IsTrue()
func (bt BoolTest) IsTrue() {
	bt.Helper()
	IsTrue(bt.T, bt.got, fmt.Sprintf("%s/is_true", bt.name))
}

// fails the test if a specified bool is not false.  An optional
// name/message may be specified; if no name is specified then
// a name of "is false" is used.
// Example:
//
//	test.IsFalse(t, got)
//
// This function is a short-hand convenience for:
//
//	test.Bool(t, name, got).IsFalse()
//
// The following two tests are exactly equivalent:
//
//	test.IsFalse(t, got)
//	test.Bool(t, "is false", got).IsFalse()
func IsFalse(t *testing.T, got bool, name ...string) {
	t.Helper()

	if len(name) == 0 {
		name = []string{"is false"}
	}

	t.Run(name[0], func(t *testing.T) {
		t.Helper()
		if got {
			t.Errorf("\nwanted: false\ngot   : %v", got)
		}
	})
}

// fails the test if a specified bool is not true.  An optional
// name/message may be specified.  If no name is specified then
// a name of "is true" is used.
//
// Example:
//
//	test.IsTrue(t, got)
//
// This function is a short-hand convenience for:
//
//	test.Bool(t, name, got).IsTrue()
//
// The following two tests are exactly equivalent:
//
//	test.IsTrue(t, got)
//	test.Bool(t, "is true", got).IsTrue()
func IsTrue(t *testing.T, got bool, name ...string) {
	t.Helper()

	if len(name) == 0 {
		name = []string{"is true"}
	}

	t.Run(name[0], func(t *testing.T) {
		t.Helper()
		if !got {
			t.Errorf("\nwanted: true\ngot   : %v", got)
		}
	})
}
