package emptiness

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type Matcher struct {
	// if true, nil is considered empty
	TreatNilAsEmpty bool

	// the following fields are set during evaluation of the matcher and
	// used to generate an appropriate failure report

	// hasLength is true if a length was determined for the subject
	hasLength bool

	// isEmpty is true if the subject has been determined to be empty
	isEmpty bool

	// isNil is true if the subject is nil
	isNil bool

	// method is the name of the method used to determine the length
	// (e.g., "len", "Count", "Len", or "Length")
	method string

	// _type is the type of the subject, if known
	_type string

	// len is the length of the subject as a string; a string representation
	// is used to avoid issues with large integers and the different integer
	// types used for the various supported length methods
	len string
}

func (m *Matcher) Match(subject any, opts ...any) bool {
	m.tryLen(subject)

	switch {
	case m.isNil:
		// if TreatNilAsEmpty has been set then the nilness of a `nil` value
		// is ignored and the value is treated as empty
		//
		// i.e. if isNil is true then the value is nil AND is not treated as empty,
		// therefore the match fails
		return false

	case m.hasLength:
		return m.isEmpty
	}

	// otherwise try the various supported methods for determining length/count
	// if hasLength is set after any of these calls, then the value is supported
	// and the emptiness result is available

	if tryMethods[int](m, subject); m.hasLength {
		return m.isEmpty
	}

	if tryMethods[uint](m, subject); m.hasLength {
		return m.isEmpty
	}

	if tryMethods[int64](m, subject); m.hasLength {
		return m.isEmpty
	}

	tryMethods[uint64](m, subject)

	return m.isEmpty
}

func (m *Matcher) OnTestFailure(subject any, opts ...any) []string {
	if !m.hasLength && m.isNil && m._type == "" {
		test.T().Helper()
		test.Invalid(
			"emptiness.Matcher: requires a value that is a slice, channel, or map, or is of",
			"                   a type that implements a Count(), Len(), or Length() function",
			"                   returning an int, int64, uint, or uint64.",
			"",
			fmt.Sprintf("                   A value of type %T does not meet these criteria.", subject),
		)
	}

	if m.isNil && !m.TreatNilAsEmpty && m._type != "" {
		return []string{
			"expected: <empty " + m._type + ">",
			"got     : nil " + m._type,
		}
	}

	switch m._type {
	case "slice":
		return []string{
			"expected: <empty slice>",
			"got     : len() == " + m.len,
		}
	case "string":
		return []string{
			"expected: <empty string>",
			"got     : " + opt.ValueAsString(subject, opts...),
		}
	default:
		return []string{
			"expected: <empty " + m._type + ">",
			"got     : " + m.method + "() == " + m.len,
		}
	}
}

// tryLen attempts to determine the length of the subject using the built-in
// len() function
func (m *Matcher) tryLen(v any) {
	setResult := func(isNil, isEmpty bool, slen, typ string) {
		m.isNil = isNil
		m.isEmpty = isEmpty
		m.hasLength = true
		m.len = slen
		m.method = "len"
		m._type = typ
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
	// is supported by len()

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	kind := val.Kind()
	if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && val.IsNil() {
		setResult(!m.TreatNilAsEmpty, true, "nil "+kind.String(), typ.Kind().String())
		return
	}

	var n int
	switch typ.Kind() { //nolint: exhaustive // dealing only with types that support len()
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
	typ := fmt.Sprintf("%T", v)

	setResult := func(l T, method string) {
		m.hasLength = true
		m.len = fmt.Sprintf("%d", l)
		m.isEmpty = l == 0
		m.method = method
		m._type = typ
	}

	switch v := v.(type) {
	case interface{ Count() T }:
		setResult(v.Count(), "Count")

	case interface{ Len() T }:
		setResult(v.Len(), "Len")

	case interface{ Length() T }:
		setResult(v.Length(), "Length")
	}
}
