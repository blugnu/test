package emptiness

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
	"github.com/blugnu/test/test"
)

type Matcher struct {
	// if true, nil is considered empty
	NilIsEmpty bool

	// the following fields are set during evaluation of the matcher and
	// used to generate an appropriate failure report

	// isValid is true if a valid determination of emptiness was made
	// for the subject (whether empty or not)
	isValid bool

	// hasLength is true if a length was determined for the subject
	hasLength bool

	// isEmpty is true if the subject has been determined to be empty
	isEmpty bool

	// isNil is true if the subject is nil
	isNil bool

	// method is the name of the method used to determine the length
	// (e.g., "len", "Count", "Len", or "Length")
	method string

	// customMethod is true if a method other than "len" was used to
	// determine length
	customMethod bool

	// typeName is the type of the subject, if known
	typeName string

	// len is the length of the subject as a string; a string representation
	// is used to avoid issues with large integers and the different integer
	// types used for the various supported length methods
	len string
}

func (m *Matcher) Match(subject any, opts ...any) bool {
	result := func() bool {
		switch {
		case !m.isValid:
			return false
		case m.isNil:
			return m.NilIsEmpty
		default:
			return m.isEmpty
		}
	}

	if m.tryLen(subject); m.isValid {
		return result()
	}

	// otherwise try the various supported methods for determining length/count
	// if hasLength is set after any of these calls, then the value is supported
	// and the emptiness result is available

	if tryMethods[int](m, subject); m.isValid {
		return result()
	}

	if tryMethods[uint](m, subject); m.isValid {
		return result()
	}

	if tryMethods[int64](m, subject); m.isValid {
		return result()
	}

	tryMethods[uint64](m, subject)

	return result()
}

func (m *Matcher) OnTestFailure(subject any, opts ...any) []string {
	if !m.isValid {
		test.T().Helper()
		test.Invalid(
			"emptiness.Matcher: requires a value that is a slice, channel, or map, or of a type",
			"                   that implements a Count(), Len(), or Length() function returning",
			"                   int, int64, uint, or uint64.",
			"",
			fmt.Sprintf("                   A value of type %s does not meet these criteria.", report.TypeName(subject)),
		)
	}

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return m.negatedFailureReport()
	}

	return m.failureReport(subject, opts...)
}

// failureReport generates a failure report for the case where
// the matcher was used in a non-negated context (i.e., Expect(...).Should(...))
func (m *Matcher) failureReport(subject any, opts ...any) []string {
	switch {
	case m.isNil:
		return report.ExpectedGot(
			report.Empty(m.typeName),
			report.Nil,
		)

	case m.typeName == "slice" || m.typeName == "array":
		return report.AppendSliceOrArray(
			report.Expected(report.Empty(m.typeName)),
			subject, opt.WithName(opts, "got     :")...,
		)

	case m.typeName == "string":
		return report.ExpectedGot(
			report.Empty("string"),
			report.Value(subject, opts...),
		)

	default:
		return report.ExpectedGot(
			report.Empty(m.typeName),
			m.method+"() == "+m.len,
		)
	}
}

// negatedFailureReport generates a failure report for the case where
// the matcher was used in a negated context (i.e., Expect(...).ShouldNot(...))
func (m *Matcher) negatedFailureReport() []string {
	if m.typeName == "string" || m.typeName == "array" {
		return report.Expected("<not empty>")
	}

	expected := "<not empty>"
	if m.NilIsEmpty {
		expected = "<not empty or nil>"
	}
	if m.customMethod {
		expected += " (using " + m.method + "() > 0)"
	}

	switch {
	case m.isNil:
		return report.ExpectedGot(expected, report.Nil)

	case m.customMethod && m.hasLength:
		return report.ExpectedGot(
			expected+" (using "+m.method+"() > 0)",
			m.method+"() == "+m.len,
		)

	case m.NilIsEmpty:
		return report.ExpectedGot(expected, "<empty>")

	default:
		return report.Expected(expected)
	}
}

// tryLen attempts to determine the length of the subject using the built-in
// len() function
func (m *Matcher) tryLen(v any) {
	setResult := func(isNil, isEmpty bool, slen, typ string) {
		m.isValid = true
		m.isNil = isNil
		m.isEmpty = isEmpty
		m.hasLength = true
		m.len = slen
		m.method = "len"
		m.typeName = typ
	}

	switch got := v.(type) {
	case string:
		setResult(false, len(got) == 0, strconv.Itoa(len(got)), "string")
		return

	case nil:
		m.isNil = true
		return
	}

	// not a string and not nil so we need to determine whether the value
	// is supported by len(), for which we need the type and value reflections
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	// before calling len(), check for nil slices, maps, and channels
	kind := val.Kind()
	if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && val.IsNil() {
		setResult(true, true, "nil "+kind.String(), typ.Kind().String())
		return
	}

	var n int
	switch typ.Kind() { //nolint: exhaustive // only dealing with types that support len()
	case reflect.Array:
		n = typ.Len()
	case reflect.Chan, reflect.Map, reflect.Slice:
		n = val.Len()
	default:
		return
	}

	setResult(false, n == 0, strconv.Itoa(n), typ.Kind().String())
}

// tryMethods attempts to determine the length of the subject using
// Count(), Len(), or Length() methods, if they are implemented
func tryMethods[T int | uint | int64 | uint64](m *Matcher, v any) {
	supports := func(method string) {
		m.isValid = true
		m.method = method
		m.customMethod = true
		m.typeName = report.TypeName(v)

		// check for nil pointer receiver
		if val := reflect.ValueOf(v); val.Kind() == reflect.Pointer {
			m.isNil = val.IsNil()
		}
	}

	setLength := func(l T, method string) {
		m.hasLength = true
		m.len = fmt.Sprintf("%d", l)
		m.isEmpty = l == 0
		m.method = method
	}

	type implementsCount interface{ Count() T }
	type implementsLen interface{ Len() T }
	type implementsLength interface{ Length() T }

	switch v.(type) {
	case implementsCount:
		supports("Count")
	case implementsLen:
		supports("Len")
	case implementsLength:
		supports("Length")
	default:
		// type does not implement any supported method
		return
	}

	if m.isNil {
		// nil receiver; length cannot be determined
		return
	}

	// for a non-nil receiver, call whichever supported method is available
	switch v := v.(type) {
	case implementsCount:
		setLength(v.Count(), "Count")

	case implementsLen:
		setLength(v.Len(), "Len")

	case implementsLength:
		setLength(v.Length(), "Length")
	}
}
