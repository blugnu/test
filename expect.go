package test

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

// Matcher[T] is a generic interface that defines a method for
// matching a value of type T. This is the interface that all
// matchers must implement to be accepted by the To() and
// ToNot() methods.
type Matcher[T any] interface {
	Match(T, ...any) bool
}

// Expect creates an expectation for the given value.  The value
// may be of any type.
//
// # Supported Options
//
//	string    // a name for the expectation; the name is used in
//	          // the failure message if the expectation fails.
func Expect[T any](value T, opts ...any) expectation[T] {
	t := GetT()
	return expectation[T]{
		t:        t,
		subject:  value,
		name:     opt.Name(opts),
		testName: t.Name(),
	}
}

// expectation[T] is a type that represents an expectation in a test. It
// holds the test runner in scope at the time the expectation was
// expressed, the name of the expectation (optional) and the value to
// which the expectation applies.
type expectation[T any] struct {
	t        TestingT
	subject  T
	name     string
	testName string
}

// err fails a test with an optional message.  If specified, the
// message may be given as a string or []string.  Only the first msg
// is used; any additional msg args are ignored.
//
// If no msg is supplied, the test fails with no message.  If the
// expectation has a name, the test fails with the message "<name> failed".
//
// If the first msg is a string it is used as the message. If the
// expectation has a name, it is prepended to the string.
//
// If the first msg is a []string, it is used as the message with each
// string in the slice on a new line. If the expectation has a name,
// it is prepended to the first string in the slice.
func (e expectation[T]) err(msg any) {
	e.t.Helper()

	handleEmpty := func() {
		e.t.Helper()

		if e.name != "" {
			e.t.Errorf("test failed (%s)", e.name)
			return
		}
		e.t.Error("test failed")
	}

	switch msg := msg.(type) {
	case nil:
		handleEmpty()

	case string:
		if len(msg) == 0 {
			handleEmpty()
			return
		}
		if e.name != "" {
			msg = e.name + ": " + msg
		}
		e.t.Error(msg)

	case []string:
		if len(msg) == 0 {
			handleEmpty()
			return
		}

		var rpt string
		var indent string
		if e.name != "" {
			rpt = "\n" + e.name + ":"
			indent = "  "
		}
		for _, s := range msg {
			rpt += "\n" + indent + s
		}
		e.t.Error(rpt)

	default:
		e.err(fmt.Sprintf("\ntest failed with: %v", msg))
	}
}

// errf fails a test with an optional message.  If specified, the
// message may be given as a string containing formatting verbs and
// a variadic number of arguments.
//
// If a message string is given without formatting verbs, it is used
// as the message and any additional args are ignored.
//
// If the expectation is named, the name is prepended to the message.
func (e expectation[T]) errf(s string, args ...any) {
	e.t.Helper()
	if e.name != "" {
		e.t.Errorf("\n"+e.name+": "+s, args...)
	} else {
		e.t.Errorf("\n"+s, args...)
	}
}

// fail determines how the test failure should be reported, formats the
// test failure report then fails the test with the report.
func (e expectation[T]) fail(matcher Matcher[T], opts ...any) {
	// e.t = GetT()
	e.t.Helper()

	report := e.getTestFailureReporter(opts...)
	if report == nil {
		report = matcher
	}

	switch report := report.(type) {
	case interface{ OnTestFailure(...any) string }:
		e.err([]string{report.OnTestFailure(opts...)})
	case interface{ OnTestFailure(...any) []string }:
		e.err(report.OnTestFailure(opts...))
	case interface{ OnTestFailure(T, ...any) string }:
		e.err([]string{report.OnTestFailure(e.subject, opts...)})
	case interface{ OnTestFailure(T, ...any) []string }:
		e.err(report.OnTestFailure(e.subject, opts...))
	default:
		exp := e.getExpected(matcher)

		var ef, gf string
		if f, ok := report.(interface{ FormatValue(any, ...any) string }); ok {
			ef = f.FormatValue(exp, opts...)
			gf = f.FormatValue(e.subject, opts...)
		} else {
			ef = fmt.Sprintf("%v", exp)
			gf = fmt.Sprintf("%v", e.subject)
		}

		if exp == nil {
			e.errf("got %s", gf)
			return
		}

		// if the expected and got values are small, use a one line report
		if len(ef) < 10 && len(gf) < 10 {
			e.errf("expected %s, got %s", ef, gf)
			return
		}

		// otherwise, use a multi-line report
		e.err([]string{
			fmt.Sprintf("expected: %s", ef),
			fmt.Sprintf("got     : %s", gf),
		})
	}
}

// getExpected returns the expected value from the matcher.  If the
// matcher is a struct with an Expected or expected field, that value
// is returned.  If the matcher implements an Expected() method, that
// value is returned.  Otherwise, nil is returned.
func (e expectation[T]) getExpected(matcher Matcher[T]) any {
	// check for an Expected field if the matcher is a struct
	if reflect.ValueOf(matcher).Kind() == reflect.Struct {
		m := reflect.Indirect(reflect.ValueOf(matcher))
		if fld := m.FieldByName("Expected"); fld.IsValid() {
			return fld.Interface()
		}
	}

	// use a suitable Expected() method if implemented by the matcher
	switch m := matcher.(type) {
	case interface{ Expected() T }:
		return m.Expected()
	case interface{ Expected() any }:
		return m.Expected()
	}

	// unable to find an expected value, return nil
	return nil
}

// TestFailure returns the first TestFailure[T] function found in the provided options.
// It checks for the following types:
//
// - interface{ TestFailure(...any) string }
// - interface{ TestFailure(...any) []string }
// - interface{ TestFailure(T, ...any) string }
// - interface{ TestFailure(T, ...any) []string }
//
// If no matching function is found, nil is returned.
func (e expectation[T]) getTestFailureReporter(opts ...any) any {
	for _, opt := range opts {
		switch opt := opt.(type) {
		case
			interface{ OnTestFailure(...any) string },
			interface{ OnTestFailure(...any) []string },
			interface{ OnTestFailure(T, ...any) string },
			interface{ OnTestFailure(T, ...any) []string }:
			return opt
		default:
			continue
		}
	}
	return nil
}

// To applies a matcher to the expectation.  If the matcher
// does not match the value, the test fails.  The matcher may
// be a value of any type that implements the Matcher[T]
// interface.
//
// In addition to a matcher, the function also accepts a variadic list
// of options. While some options are applied directly by the function; all
// options are also passed to the matcher.
//
// # Supported Options
//
//	opt.FailureReport(func)      // a function that provides a custom
//	                             // test failure report if the test fails.
//	                             //
//	                             // the func must be of the form:
//	                             //
//	                             //    func(...any) []string
//
//	opt.OnFailure(string)        // a simple string to output as the
//	                             // failure report if the test fails.
func (e expectation[T]) To(matcher Matcher[T], opts ...any) {
	e.t.Helper()

	if matcher == nil {
		invalidTest("test.To: a matcher must be specified")
		return
	}

	if !matcher.Match(e.subject, opts...) {
		e.fail(matcher, opts...)
	}
}

// ToNot applies a matcher to the expectation.  If the matcher
// matches the value, the test fails.
//
// In addition to a matcher, the function also accepts a variadic list
// of options. While some options are applied directly by the function; all
// options are also passed to the matcher.
//
// # Supported Options
//
//	opt.FailureReport(func)      // a function that provides a custom
//	                             // test failure report if the test fails.
//	                             //
//	                             // the func must be of the form:
//	                             //
//	                             //    func(...any) []string
//
//	opt.OnFailure(string)        // a simple string to output as the
//	                             // failure report if the test fails.
func (e expectation[T]) ToNot(matcher Matcher[T], opts ...any) {
	e.t.Helper()

	if matcher.Match(e.subject, opts...) {
		e.fail(matcher, append(opts, opt.ToNotMatch(true))...)
	}
}

// Is tests the value of the expectation against some expected
// value.
//
// The function behaves differently depending on the values and
// types of the expected and actual values:
//
//   - If both values are nil, the test passes;
//
//   - If either value is nil and the other is not, the test fails;
//
//   - If both values implement the error interface, the test passes
//     if the error being tested satisfies errors.Is(expected);
//
//   - Otherwise, the values are compared using reflect.DeepEqual
//     or a comparison function supplied in the options;
//
// i.e. for non-nil, non-error values, an Is() test is equivalent to:
//
//	Expect(got).To(DeepEqual(expected), opts...)
//
// # Supported Options
//
//	func(a, b any) bool         // a function to compare the values
//	                            // (overriding the use of reflect.DeepEqual)
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) Is(expected T, opts ...any) {
	e.t.Helper()

	switch {
	case any(expected) == nil:
		e.IsNil()
		return

	case any(expected) != nil && any(e.subject) == nil:
		e.errf("expected %v, got nil", expected)

	default:
		experr, _ := any(expected).(error)
		goterr, _ := any(e.subject).(error)
		if experr != nil && goterr != nil {
			ExpectTrue(
				errors.Is(goterr, experr),
				opt.FailureReport(func(...any) []string {
					return []string{
						fmt.Sprintf("expected error: %v", experr),
						fmt.Sprintf("got           : %v", goterr),
					}
				}),
			)
			return
		}
		Expect(e.subject, e.name).To(DeepEqual(expected), opts...)
	}
}

// IsNil checks that the value of the expectation is nil.  If the
// value is not nil, the test fails.  If the value is nil, the test
// passes.
//
// If the value being tested does not support a nil value the test
// will fail and produce a report similar to:
//
//	test.IsNil: values of type '<type>' are not nilable
//
// # Supported Options
//
//	opt.QuotedStrings(bool)     // determines whether any non-nil string
//	                            // values are quoted in any test failure
//	                            // report.  The default is false (string
//	                            // values are quoted).
//	                            //
//	                            // If the value is not a string type this
//	                            // option has no effect.
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsNil(opts ...any) {
	e.t.Helper()

	switch result := isNil(e.subject).(type) {
	case bool:
		if result {
			return
		}

		if report, ok := opt.Get[opt.FailureReport](opts); ok {
			e.err(report.OnTestFailure(true))
			return
		}

		switch got := any(e.subject).(type) {
		case error:
			e.errf("expected nil, got error: %v", got)
		default:
			if reflect.ValueOf(e.subject).Kind() == reflect.Pointer {
				v := reflect.Indirect(reflect.ValueOf(e.subject))
				if v.Kind() == reflect.String {
					e.errf("expected nil, got &(%s)", opt.ValueAsString(v.String(), opts...))
					return
				}
				e.errf("expected nil, got &(%#v)", v)
				return
			}
			e.errf("expected nil, got %#v", got)
		}

	case error:
		// since it is expected that the value being tested is nil, the value
		// being of a non-nilable type is an invalid test
		invalidTest(
			fmt.Sprintf("test.IsNil: values of type '%T' are not nilable", e.subject),
		)
	}
}

// IsNotNil checks that a specified value is not nil.  If the value
// is not nil, the test fails.  If the value is nil, the test passes.
//
// NOTE: If the value being tested does not support a nil value the
// test will pass.  This is to allow for testing values that may be
// nilable or non-nilable.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsNotNil(opts ...any) {
	e.t.Helper()

	switch result := isNil(e.subject).(type) {
	case bool:
		if !result {
			return
		}

		if report, ok := opt.Get[opt.FailureReport](opts); ok {
			e.err(report.OnTestFailure(true))
			return
		}

		e.errf("expected not nil, got nil")

	case error:
		// since it is expected that the value being tested is not nil, the
		// fact that the value is not nilable is not a failure.
		//
		// It may be a pointless test, but we cannot assume that.
		//
		// It may also be a test where the value being tested is an 'any'
		// value which may hold either a nilable or non-nilable value
		// but where in any event it is expected that the value is not nil.
		//
		// This is different to the IsNil() test where the value being tested
		// is expected to be nil and therefore must also be of nilable type.
		return
	}
}

// isNil returns true if the supplied value is nil or false if
// the value is not nil and is of a type that could be nil.
//
// If the supplied value is not nil and is of a type that does not
// support a nil value, the function will return ErrNotNilable.
func isNil(v any) any {
	if v == nil {
		return true
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		if v := reflect.ValueOf(v); v != reflect.Zero(reflect.TypeOf(v)) && !v.IsNil() {
			return false
		}
		return true
	}
	return ErrNotNilable
}
