package test

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/blugnu/test/matchers/slices"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type lengthResult struct {
	hasLength bool
	string
	isZero bool
	isNil  bool
	method string
	_type  string
}

func hasLen(v any) lengthResult {
	switch got := v.(type) {
	case string:
		return lengthResult{
			hasLength: true,
			string:    strconv.Itoa(len(got)),
			isZero:    len(got) == 0,
			method:    "len",
			_type:     "string",
		}
	case nil:
		return lengthResult{isNil: true}
	default:
		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)
		kind := val.Kind()
		if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && val.IsNil() {
			return lengthResult{isNil: true, _type: typ.Kind().String()}
		}

		var n int
		switch typ.Kind() {
		case reflect.Array:
			n = typ.Len()
		case reflect.Chan, reflect.Map, reflect.Slice:
			n = val.Len()
		default:
			return lengthResult{hasLength: false}
		}
		return lengthResult{
			hasLength: true,
			string:    strconv.Itoa(n),
			isZero:    n == 0,
			method:    "len",
			_type:     typ.Kind().String(),
		}
	}

}

func hasLength[T int | uint | int64 | uint64](v any) lengthResult {
	switch v := v.(type) {
	case interface{ Count() T }:
		return lengthResult{
			hasLength: true,
			string:    fmt.Sprintf("%d", v.Count()),
			isZero:    v.Count() == 0,
			method:    "Count",
		}

	case interface{ Len() T }:
		return lengthResult{
			hasLength: true,
			string:    fmt.Sprintf("%d", v.Len()),
			isZero:    v.Len() == 0,
			method:    "Len",
		}

	case interface{ Length() T }:
		return lengthResult{
			hasLength: true,
			string:    fmt.Sprintf("%d", v.Length()),
			isZero:    v.Length() == 0,
			method:    "Length",
		}
	default:
		return lengthResult{hasLength: false}
	}
}

// getLength returns the length of the value v as a string
func (e expectation[T]) getLengthResult(v any) lengthResult {
	var result lengthResult
	if result = hasLen(v); result.hasLength || result.isNil {
		return result
	}
	if result = hasLength[int](v); result.hasLength {
		return result
	}
	if result = hasLength[uint](v); result.hasLength {
		return result
	}
	if result = hasLength[int64](v); result.hasLength {
		return result
	}
	return hasLength[uint64](v)
}

func (e expectation[T]) isEmpty(orNil bool, opts ...any) {
	e.t.Helper()

	valsFn := func(result lengthResult) []string {
		switch got := any(e.subject).(type) {
		case []string:
			report := make([]string, 1, len(got)+1)
			report[0] = "expected: <empty []string>"
			report = slices.AppendToReport(report, got, "got:", append(opts, opt.PrefixInlineWithFirstItem(true))...)
			return report
		default:
			if !orNil && result.isNil {
				if result._type == "" {
					return []string{
						"expected: <empty>",
						"got     : nil",
					}
				}
				return []string{
					"expected: <empty " + result._type + ">",
					"got     : nil " + result._type,
				}
			}

			switch result._type {
			case "string":
				return []string{
					"expected: <empty string>",
					"got     : " + opt.ValueAsString(got, opts...),
				}
			case "":
				return []string{
					"expected: <empty>",
					"got     : " + result.method + "() == " + result.string,
				}
			default:
				return []string{
					"expected: <empty " + result._type + ">",
					"got     : " + result.method + "() == " + result.string,
				}
			}
		}
	}

	result := e.getLengthResult(e.subject)
	if orNil && result.isNil {
		return
	}

	if !result.hasLength && !result.isNil {
		if orNil {
			test.Invalid(
				"IsEmptyOrNil: requires a value that is a slice, channel, or map, or is of",
				"              a type that implements a Count(), Len(), or Length() function",
				"              returning an int, int64, uint, or uint64.",
				"",
				fmt.Sprintf("              A value of type %T does not meet these criteria.", e.subject),
			)
		} else {
			test.Invalid(
				"IsEmpty: requires a value that is a string, array, slice, channel or map,",
				"         or is of a type that implements a Count(), Len(), or Length()",
				"         function returning an int, int64, uint, or uint64.",
				"",
				fmt.Sprintf("         A value of type %T does not meet these criteria.", e.subject),
			)
		}
		return
	}

	if !result.isZero {
		if report, ok := opt.Get[opt.FailureReport](opts); ok {
			e.err(report.OnTestFailure(false))
			return
		}
		e.err(valsFn(result))
	}
}

// IsEmpty checks that the value of the expectation is empty.
//
// The test passes if the value is empty, otherwise it fails.
//
// NOTE: A nil value is not considered empty. To test for an empty value
// that may be nil, use IsEmptyOrNil() instead.
//
// i.e. an empty slice will pass this test, but a nil slice will not.
//
// This test may be used to check for empty strings, arrays, slices,
// channels, maps and any type that implement a Count(), Len() or
// Length() function returning int, int64, uint, or uint64.
//
// If a value of any other type is tested, the test fails with a message
// similar to:
//
//	test.IsEmpty: requires a value of type string, array, slice, channel or map,
//	               or a type that implements a Count(), Len(), or Length() function
//	               returning int, int64, uint, or uint64.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsEmpty(opts ...any) {
	e.t.Helper()
	e.isEmpty(false, opts...)
}

// IsEmptyOrNil checks that the value of the expectation is empty or nil.
//
// The test passes if the value is empty or nil, otherwise it fails.
//
// This test may be used to check for empty slices, channels, maps or value
// of any nilable type that implement a Count(), Len() or Length() function
// returning int, int64, uint, or uint64.
//
// If a value of any other type is tested, the test fails with a message
// similar to:
//
//		test.IsEmpty: requires a value of type slice, channel or map, or a type
//	                that implements a Count(), Len(), or Length() function
//		               returning int, int64, uint, or uint64.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsEmptyOrNil(opts ...any) {
	e.t.Helper()
	e.isEmpty(true, opts...)
}

// IsNotEmpty checks that the value of the expectation is not empty.
//
// The test passes if the value is not empty, otherwise it fails.
//
// This test may be used to check for non-empty strings, arrays, slices,
// channels, maps and any type that implement a Count(), Len() or
// Length() function returning int, int64, uint, or uint64.
//
// If a value of any other type is tested, the test fails with a message
// similar to:
//
//	test.IsNotEmpty: requires a value of type string, array, slice, channel or map,
//	                 or a type that implements a Count(), Len(), or Length() function
//	                 returning int, int64, uint, or uint64.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsNotEmpty(opts ...any) {
	e.t.Helper()

	valsFn := func(result lengthResult) string {
		switch any(e.subject).(type) {
		case []string:
			return "expected: <non-empty []string>"
		default:
			switch result._type {
			case "string":
				return "expected: <non-empty string>"
			case "":
				return "expected: " + result.method + "() > 0"
			default:
				return "expected: <non-empty " + result._type + ">, " + result.method + "() > 0 "
			}
		}
	}

	result := e.getLengthResult(e.subject)
	if !result.hasLength {
		test.Invalid(
			"IsNotEmpty: requires a value of type string, array, slice, channel or map,",
			"            or a type that implements a Count(), Len(), or Length() function",
			"            returning int, int64, uint, or uint64.",
			"",
			fmt.Sprintf("            A value of type %T does not meet these criteria.", e.subject),
		)
		return
	}

	if result.isZero {
		if report, ok := opt.Get[opt.FailureReport](opts); ok {
			e.err(report.OnTestFailure(false))
			return
		}
		e.err(valsFn(result))
	}
}
