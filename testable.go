package test

import (
	"fmt"
	"testing"
)

// captures a generic value to be tested (the 'got' value) together with
// a name/description and a formatting function for producing a string
// representation of that value (or the corresponding desired or 'wanted'
// value, where required) in a test failure report.
//
// This type is not intended to be used directly. Instead, use one of the
// types in this package which composes this type and provides methods
// for testing values of type T, instantiated by that testable type.
type testable[T any] struct {
	*testing.T
	got  T
	ffn  func(T) string
	name string
}

// returns a new test for a specified value, using a specified testing.T
func newTestable[T any](t *testing.T, got T, opt ...any) testable[T] {
	t.Helper()

	n := "got"
	f := FormatDefault
	ffn := *new(func(T) string)

	checkOptTypes(t, optTypes(n, f, ffn), opt...)
	getOpt(&n, opt...)
	getOpt(&f, opt...)
	getOpt(&ffn, append(opt, func(v T) string { return fmt.Sprintf(string(f), v) })...)

	return testable[T]{T: t, name: n, ffn: ffn, got: got}
}

func (t testable[T]) errorf(tt *testing.T, format string, args ...any) {
	tt.Helper()
	tt.Errorf("\n"+format, args...)
}

func (t testable[T]) fail(tt *testing.T, wanted T, notes ...string) {
	tt.Helper()
	report := fmt.Sprintf("wanted: %s\ngot   : %s", t.ffn(wanted), t.ffn(t.got))
	for n := range notes {
		report += "\n" + notes[n]
	}
	t.errorf(tt, report)
}

func (t testable[T]) format(v T) string {
	return t.ffn(v)
}

func (t testable[T]) run(f func(*testing.T)) {
	t.Helper()
	t.T.Run(t.name, f)
}

func (t testable[T]) Run(name string, f func(*testing.T)) {
	t.Helper()
	t.T.Run(fmt.Sprintf("%s/%s", t.name, name), func(tt *testing.T) {
		tt.Helper()
		f(tt)
	})
}
