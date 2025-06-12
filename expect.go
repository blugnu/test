package test

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/blugnu/test/matchers/matcher"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
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
func Expect[T any](value T, opts ...any) *expectation[T] {
	t := GetT()
	return &expectation[T]{
		t:        t,
		subject:  value,
		name:     opt.Name(opts),
		testName: t.Name(),
	}
}

// Require creates an expectation for the given value which is required
// to pass.  If the expectation is not met, execution continues with the
// *next* test (if any); no further expectations in the current test will
// be evaluated.
//
// This is a convenience function that is equivalent to passing the
// opt.Required() or opt.IsRequired(true) option to a matcher invoked
// using Expect(), i.e. the following are equivalent:
//
//	Expect(value).To(Equal(expected), opt.IsRequired(true))
//	Expect(value).To(Equal(expected), opt.Required())
//	Require(value).To(Equal(expected))
//
// # Supported Options
//
//	string    // a name for the expectation; the name is used in
//	          // the failure message if the expectation fails.
func Require[T any](value T, opts ...any) *expectation[T] {
	t := GetT()
	return &expectation[T]{
		t:        t,
		subject:  value,
		name:     opt.Name(opts),
		testName: t.Name(),
		required: true,
	}
}

// expectation[T] is a type that represents an expectation in a test. It
// holds the TestingT from the test frame in scope at the time the
// expectation was expressed, the name of the expectation (optional) and
// the value to which the expectation applies.
type expectation[T any] struct {
	t        TestingT
	subject  T
	name     string
	testName string

	// required indicates whether the expectation is required to pass; if true
	// and the expectation fails, the current test is failed immediated and no
	// further expectations in the current test are evaluated
	required bool
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
func (e *expectation[T]) err(msg any) {
	e.t.Helper()

	// the fail funcs will set Error by default...
	fail := e.t.Error
	failf := e.t.Errorf

	// ..but if the expectation is marked as required, fail funcs will
	// use Fatal instead
	if e.required {
		fail = e.t.Fatal
		failf = e.t.Fatalf
	}

	handleEmpty := func() {
		e.t.Helper()

		if e.name != "" {
			failf("test failed (%s)", e.name)
			return
		}
		fail("test failed")
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
		fail(msg)

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
		// e.t.Error(rpt)
		fail(rpt)

	default:
		// e.t.Errorf("\ntest failed with: %v", msg)
		failf("\ntest failed with: %v", msg)
	}
}

// fail determines how the test failure should be reported, formats the
// test failure report then fails the test with the report.
func (e *expectation[T]) fail(matcher any, opts ...any) {
	e.t.Helper()

	// check for a custom test failure report function in the
	// options; if none are provided then the matcher is
	// expected to implement a test failure reporter (though it may not;
	// that will be determined later, if needed)
	report := e.getTestFailureReporter(opts...)
	if report == nil {
		report = matcher
	}

	// expectation.required may be preset or may have been specified
	// as an option
	e.required = e.required || opt.IsSet(opts, opt.IsRequired(true))

	switch report := report.(type) {
	case interface{ OnTestFailure(...any) string }:
		e.err([]string{report.OnTestFailure(opts...)})
	case interface{ OnTestFailure(...any) []string }:
		e.err(report.OnTestFailure(opts...))
	case interface{ OnTestFailure(T, ...any) string }:
		e.err([]string{report.OnTestFailure(e.subject, opts...)})
	case interface{ OnTestFailure(T, ...any) []string }:
		e.err(report.OnTestFailure(e.subject, opts...))
	case interface{ OnTestFailure(any, ...any) string }:
		e.err([]string{report.OnTestFailure(e.subject, opts...)})
	case interface{ OnTestFailure(any, ...any) []string }:
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
			e.err(fmt.Sprintf("got %s", gf))
			return
		}

		// if the expected and got values are small, use a one line report
		if len(ef) < 10 && len(gf) < 10 {
			e.err(fmt.Sprintf("expected %s, got %s", ef, gf))
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
func (e *expectation[T]) getExpected(matcher any) any {
	// check for an Expected field if the matcher is a struct or pointer
	// to struct
	// check for an Expected field if the matcher is a struct
	rv := reflect.ValueOf(matcher)
	if rv.Kind() == reflect.Struct || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct) {
		m := reflect.Indirect(rv)
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

// getTestFailureReporter returns the first function in the provided
// options that implements a test failure reporter interface.
//
// The function checks for the following interfaces:
//
// - interface{ OnTestFailure(...any) string }
// - interface{ OnTestFailure(...any) []string }
// - interface{ OnTestFailure(T, ...any) string }
// - interface{ OnTestFailure(T, ...any) []string }
//
// If no matching function is found, nil is returned.
func (e *expectation[T]) getTestFailureReporter(opts ...any) any {
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

// Should applies a matcher to the expectation.  If the matcher
// does not match the value, the test fails.
//
// The matcher must implement the AnyMatcher interface:
//
//	type AnyMatcher interface {
//	    Match(got any, opts ...any) bool
//	}
//
// An AnyMatcher is not type-safe and may fail as an invalid
// test if used incorrectly.  For certain matchers this is
// unavoidable.
//
// Refer to the documentation for individual matchers for any
// specific requirements or limitations and details of any
// options that are supported.
//
// # Supported Options
//
// The function also accepts a variadic list of options. Some options
// are supported directly; all options are also passed to the matcher
// to allow the matcher to apply those it may support.
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
func (e *expectation[T]) Should(match matcher.ForAny, opts ...any) {
	e.t.Helper()

	if match == nil {
		test.Invalid("test.Should: a matcher must be specified")
	}

	if !match.Match(e.subject, opts...) {
		e.fail(match, opts...)
	}
}

// ShouldNot applies a matcher to the expectation.  If the matcher
// matches the value, the test fails.
// The matcher must implement the AnyMatcher interface:
//
//	type AnyMatcher interface {
//	    Match(got any, opts ...any) bool
//	}
//
// An AnyMatcher is not type-safe and may fail as an invalid
// test if used incorrectly.  For certain matchers this is
// unavoidable.
//
// Refer to the documentation for individual matchers for any
// specific requirements or limitations and details of any
// options that are supported.
//
// # Supported Options
//
// The function also accepts a variadic list of options. Some options
// are supported directly; all options are also passed to the matcher
// to allow the matcher to apply those it may support.
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
func (e *expectation[T]) ShouldNot(match matcher.ForAny, opts ...any) {
	e.t.Helper()

	if match == nil {
		test.Invalid("test.ShouldNot: a matcher must be specified")
	}

	opts = append(opts, opt.ToNotMatch(true))

	if match.Match(e.subject, opts...) {
		e.fail(match, opts...)
	}
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
func (e *expectation[T]) To(matcher matcher.ForType[T], opts ...any) {
	e.t.Helper()

	if matcher == nil {
		test.Invalid("test.To: a matcher must be specified")
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
func (e *expectation[T]) ToNot(matcher matcher.ForType[T], opts ...any) {
	e.t.Helper()

	opts = append(opts, opt.ToNotMatch(true))

	if matcher.Match(e.subject, opts...) {
		e.fail(matcher, opts...)
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
// and for error values, an Is() test is equivalent to:
//
//	Expect(errors.Is(got, expected)).To(BeTrue(), opts...)
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
func (e *expectation[T]) Is(expected T, opts ...any) {
	e.t.Helper()

	switch {
	case any(expected) == nil:
		e.IsNil()
		return

	case any(expected) != nil && any(e.subject) == nil:
		e.err(fmt.Sprintf("expected %v, got nil", expected))

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
