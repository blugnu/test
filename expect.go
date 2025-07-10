package test

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/blugnu/test/matchers/matcher"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

// AnyMatcher is the interface implemented by matchers that can test
// any type of value.  It is used to apply matchers to expectations
// that are not type-specific.
//
// It is preferable to use the Matcher[T] interface for type-safe
// expectations; AnyMatcher is provided for situations where the
// compatible types for a matcher cannot be enforced at compile-time.
//
// When implementing an AnyMatcher, it is important to ensure that
// the matcher fails a test if it is not used correctly, i.e. if the
// matcher is not compatible with the type of the value being tested.
//
// An AnyMatcher must be used with the Expect().Should() matching
// function; they may also be used with Expect(got).To() where the got
// value is of type `any`, though this is not recommended.
type AnyMatcher interface {
	Match(got any, opts ...any) bool
}

// Matcher[T] is the interface implemented by matchers that can test
// a value of type T.  It is used to apply matchers to expectations
// that are type-specific and type-safe.
//
// Note that not all type-safe matchers implement a generic interface;
// a matcher that implements Match(got X, opts ...any) bool, where X is
// a formal, literal type (i.e. not generic) is also a type-safe matcher.
//
// Matcher[T] describes the general form of a type-safe matcher.
//
// Generic matchers are able to leverage the type system to ensure
// that the matcher is used correctly with a variety of types, i.e. where
// the type of the Expect() value satisfies the constraints of the matcher
// type, T.  The equals.Matcher[T comparable] uses this approach, for
// example, to ensure that the value being tested is comparable
// with the expected value (since the matcher uses the == operator for
// equality testing).
type Matcher[T any] interface {
	Match(got T, opts ...any) bool
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
		required: opt.IsSet(opts, opt.IsRequired(true)),
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

	// required indicates whether the expectation is required to pass.
	// If true and the expectation is not met, the test will fail immediately
	// and no further expectations in the current test will be evaluated.
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

	errorFn := e.t.Error
	if e.required {
		errorFn = e.t.Fatal
	}

	msg = e.errMsg(msg)
	switch msg := msg.(type) {
	case string:
		if e.name != "" {
			msg = e.name + ": " + msg
		}
		errorFn(msg)

	case []string:
		var rpt string
		var indent string

		if e.name != "" {
			rpt = "\n" + e.name + ":"
			indent = "  "
		}

		for _, s := range msg {
			rpt += "\n" + indent + s
		}
		errorFn(rpt)

		// errMsg returns a string or []string, so we can safely use a type
		// switch here to handle both cases without a default case
	}
}

func (e *expectation[T]) errMsg(msg any) any {
	const failed = "test failed"

	if msg == nil {
		return failed
	} else if s, ok := msg.(string); ok && len(s) == 0 {
		return failed
	} else if s, ok := msg.([]string); ok && len(s) == 0 {
		return failed
	}

	switch msg := msg.(type) {
	case string:
		return msg

	case []string:
		return msg

	default:
		return fmt.Sprintf("%s: %v", failed, msg)
	}
}

// defaultFailureReport presents a default test failure report
// for the expectation. It is used when a matcher does not provide
// a specific failure report and no failure reporting option is
// present.
func (e *expectation[T]) defaultFailureReport(reporter any, matcher any, opts ...any) {
	e.t.Helper()

	exp := e.getExpected(matcher)

	var ef, gf string
	if f, ok := reporter.(interface{ FormatValue(any, ...any) string }); ok {
		ef = f.FormatValue(exp, opts...)
		gf = f.FormatValue(e.subject, opts...)
	} else {
		ef = fmt.Sprintf("%v", exp)
		gf = fmt.Sprintf("%v", e.subject)
	}

	switch {
	// no expected value, just report the got value
	case exp == nil:
		e.err("got " + gf)

	// expected and got values are small, use a one line report
	case len(ef) < 10 && len(gf) < 10:
		e.err(fmt.Sprintf("expected %s, got %s", ef, gf))

	// otherwise, use a multi-line report
	default:
		e.err([]string{
			"expected: " + ef,
			"got     : " + gf,
		})
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

	// expectation.required may be preset or may be specified as an option
	e.required = e.required || opt.IsSet(opts, opt.IsRequired(true))

	switch reporter := report.(type) {
	case interface{ OnTestFailure(...any) string }:
		e.err([]string{reporter.OnTestFailure(opts...)})
	case interface{ OnTestFailure(...any) []string }:
		e.err(reporter.OnTestFailure(opts...))
	case interface{ OnTestFailure(T, ...any) string }:
		e.err([]string{reporter.OnTestFailure(e.subject, opts...)})
	case interface{ OnTestFailure(T, ...any) []string }:
		e.err(reporter.OnTestFailure(e.subject, opts...))
	case interface{ OnTestFailure(any, ...any) string }:
		e.err([]string{reporter.OnTestFailure(e.subject, opts...)})
	case interface{ OnTestFailure(any, ...any) []string }:
		e.err(reporter.OnTestFailure(e.subject, opts...))
	default:
		e.defaultFailureReport(reporter, matcher, opts...)
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
