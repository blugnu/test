package test

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/matchers/matcher"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
	"github.com/blugnu/test/test"
)

// Expect creates an expectation for the given value.  The value
// may be of any type.
//
// # Supported Options
//
//	string                 // a name for the expectation; the name is used in
//	                       // the failure message if the expectation fails.
//
//	opt.Name(string)       // a name for the expectation; this is an alternate
//	                       // way to supply the name, equivalent to passing a
//	                       // string directly.
//
//	opt.Namef(s, args...)  // a formatted name for the expectation; this is an
//	                       // alternate way to supply the name, equivalent to
//	                       // passing the result of fmt.Sprintf() directly.
//
//	opt.IsRequired(bool)   // if true, the expectation is required to pass;
//	                       // no further expectations in the current test will be
//	                       // evaluated if the expectation fails.
//
//	opt.Required()         // equivalent to opt.IsRequired(true)
func Expect[T any](value T, opts ...any) *expectation[T] {
	return newExpectation(value, false, opts...)
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
//	string                 // a name for the expectation; the name is used in
//	                       // the failure message if the expectation fails.
//
//	opt.Name(string)       // a name for the expectation; this is an alternate
//	                       // way to supply the name, equivalent to passing a
//	                       // string directly.
//
//	opt.Namef(s, args...)  // a formatted name for the expectation; this is an
//	                       // alternate way to supply the name, equivalent to
//	                       // passing the result of fmt.Sprintf() directly.
func Require[T any](value T, opts ...any) *expectation[T] {
	return newExpectation(value, true, opts...)
}

// expectation[T] is a type that represents an expectation in a test. It
// holds the TestingT from the test frame in scope at the time the
// expectation was expressed, the name of the expectation (optional) and
// the value to which the expectation applies.
type expectation[T any] struct {
	// t holds the TestingT from the test frame in scope at the time
	// the expectation was created.
	t TestingT

	// subject holds the value that the expectation applies to.
	subject T

	// name holds the name of the expectation, if any.  If the expectation
	// is unnamed, name is an empty string.
	name opt.Name

	// required indicates whether the expectation is required to pass.
	required bool
}

// newExpectation creates a new expectation for the given value.  The expectation
// is marked as required if the required argument is true or if the options
// contain the opt.IsRequired(true) option.
//
// If the options contain a string or opt.Name value, that is used as the name
// of the expectation; otherwise the expectation is unnamed.
func newExpectation[T any](value T, required bool, opts ...any) *expectation[T] {
	t := GetT()

	subject, _ := opt.GetName(opts)

	return &expectation[T]{
		t:        t,
		subject:  value,
		name:     subject,
		required: required || opt.IsSet(opts, opt.IsRequired(true)),
	}
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
			msg = e.name.String() + ": " + msg
		}
		errorFn(msg)

	case []string:
		rpt := slices.Clone(msg)
		if e.name != "" {
			for i := range rpt {
				rpt[i] = "  " + rpt[i]
			}
			rpt = append([]string{e.name.String() + ":"}, rpt...)
		}

		errorFn("\n" + strings.Join(rpt, "\n"))

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
func (e *expectation[T]) defaultFailureReport(reporter any, matcher any, _ ...any) {
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
		e.err(report.ExpectedGot(ef, gf))
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
	_ = e // Receiver is not used in this method, but kept for consistency
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

// Is tests the value of the expectation against some expected
// value.
//
// The function behaves differently depending on the values and
// types of the subject and expected values.  The different behaviours
// are intuitive and sensible for common cases:
//
//   - If both values are nil, the test passes;
//
//   - If either value is nil and the other is not, the test fails;
//
//   - If both values implement the error interface, the test passes
//     if the subject error satisfies [errors.Is] for the expected error;
//
//   - Otherwise, the values are compared using [reflect.DeepEqual]
//     or a comparison function supplied in the options;
//
// i.e. for non-nil, non-error values, with no custom comparison function,
// an Is(target) test is equivalent to:
//
//	Expect(got).To(DeepEqual(target), opts...)
//
// For error values, an Is(target) test is equivalent to:
//
//	Expect(got).To(BeError(target), opts...)
//
// # Common Options
//
//	[opt.OnFailure]
//	[opt.IsRequired]
//	[opt.Required]
//	string                      // equivalent to opt.OnFailure(string)
//
// # Additional Options
//
//	func(a, b T) bool           // a function to compare the values;
//								// ignored if either value is nil, or both
//								// values implement the error interface, or
//								// the type T implements Equal(T) bool
//
//	func(a, b any) bool         // a function to compare the values;
//								// ignored if either value is nil, or both
//								// values implement the error interface, or
//								// the type T implements Equal(T) bool, or a
//								// func(a, b T) bool option is provided
//
// All options are passed to matchers used by the function.  For additional options
// refer to the documentation for [expectation.IsNil], [] [DeepEqual].
func (e *expectation[T]) Is(expected T, opts ...any) {
	e.t.Helper()

	// if a plain string is provided as an option, convert it to
	// an opt.OnFailure option
	opts = internal.WithStringAsOnFailure(opts)

	// identify whether subject or expected (or both) are errors
	subjectError, subjectIsError := any(e.subject).(error)
	expectedError, expectedIsError := any(expected).(error)

	switch {
	case internal.IsNil(expected):
		e.IsNil(opts...)

	case any(expected) != nil && internal.IsNil(e.subject):
		// opts are cloned to avoid modifying the caller's slice of options when we
		// add expectation properties to be passed
		opts := slices.Clone(opts)

		// the testframe is included in opts to ensure that the test failure
		// report correctly reflects the testframe of the expectation
		opts = append([]any{e.t}, opts...)

		if e.required {
			opts = append(opts, opt.Required())
		}

		if e.name != "" {
			opts = append(opts, e.name)
		}

		// set custom failure report if none was provided
		if _, ok := opt.Get[opt.FailReporter](opts); !ok {
			opts = append(opts, opt.OnFailure(fmt.Sprintf("expected %s, got nil", report.Value(expected, opts...))))
		}

		Fail(opts...)

	case expectedIsError && subjectIsError:
		Expect(subjectError, e.name).To(BeError(expectedError), opts...)

	default:
		Expect(e.subject, e.name).To(DeepEqual(expected), opts...)
	}
}

// IsNot tests the value of the expectation against some expected
// value for a non-match.
//
// The test behaves differently depending on the value and type of
// the expected and actual values:
//
//   - IsNot(nil) is equivalent to [expectation.IsNotNil];
//
//   - If either value is [nil] and the other is not, the test passes;
//
//   - If both values implement the [error] interface, the test passes
//     if the error being tested does not satisfy [errors.Is];
//
//   - Otherwise, the values are compared using [reflect.DeepEqual],
//     or a comparison function supplied in the options;
//
// i.e. for non-nil, non-error values, an IsNot() test is equivalent to:
//
//	Expect(got).ToNot(DeepEqual(expected), opts...)
//
// and for error values, an IsNot() test is equivalent to:
//
//	Expect(errors.Is(got, expected)).To(BeFalse(), opts...)
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
func (e *expectation[T]) IsNot(expected T, opts ...any) {
	e.t.Helper()

	switch {
	case (any(expected) == nil) != (any(e.subject) == nil):
		return

	case any(expected) == nil:
		e.IsNotNil(opts...)
		return

	default:
		experr, _ := any(expected).(error)
		goterr, _ := any(e.subject).(error)
		if experr != nil && goterr != nil {
			if _, ok := opt.Get[opt.FailureReport](opts); !ok {
				opts = append(opts, opt.FailureReport(func(...any) []string {
					return []string{
						fmt.Sprintf("expected error that is not: %v", experr),
						fmt.Sprintf("got                     : %v", goterr),
					}
				}))
			}

			FailIf(errors.Is(goterr, experr), opts...)
			return
		}

		Expect(e.subject, e.name).ToNot(DeepEqual(expected), opts...)
	}
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
		test.Invalid("Should: a matcher must be specified")
	}

	opts = internal.WithStringAsOnFailure(opts)

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
		test.Invalid("ShouldNot: a matcher must be specified")
	}

	opts = append(internal.WithStringAsOnFailure(opts), opt.ToNotMatch(true))

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
// The following options are supported directly by the function:
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
//
//	string                       // equivalent to opt.OnFailure(string)
//
// All options are passed to the matcher, which may support additional
// options. Refer to the documentation for individual matchers for
// details of any additional options that are supported.
func (e *expectation[T]) To(matcher matcher.ForType[T], opts ...any) {
	e.t.Helper()

	if matcher == nil {
		test.Invalid("To: a matcher must be specified")
		return
	}

	opts = internal.WithStringAsOnFailure(opts)

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

	if matcher == nil {
		test.Invalid("ToNot: a matcher must be specified")
		return
	}

	opts = append(internal.WithStringAsOnFailure(opts), opt.ToNotMatch(true))

	if matcher.Match(e.subject, opts...) {
		e.fail(matcher, opts...)
	}
}
