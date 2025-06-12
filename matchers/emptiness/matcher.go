package emptiness

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type Matcher struct {
	// if true, nil is considered empty
	TreatNilAsEmpty bool

	// the follwoing fields are set during evaluation of the matcher and
	// used to generate an appropriate failure report

	// hasLength is true if a length was determined for the subject
	hasLength bool

	// isZero is true if the subject has a length of zero
	isZero bool

	// isNil is true if the subject is nil
	isNil bool

	// method is the name of the method used to determine the length
	// (e.g., "len", "Count", "Len", or "Length")
	method string

	// _type is the type of the subject, if known
	_type string

	// len is the length of the subject, if known to a maximum of math.MaxInt
	len int

	// string is the string representation of the length of the subject
	string
}

func (m *Matcher) Match(subject any, opts ...any) bool {
	m.hasLen(subject)
	switch {
	case m.isNil:
		return false
	case m.hasLength:
		return m.isZero
	}

	if hasLength[int](m, subject); m.hasLength {
		return m.isZero
	}

	if hasLength[uint](m, subject); m.hasLength {
		return m.isZero
	}

	if hasLength[int64](m, subject); m.hasLength {
		return m.isZero
	}

	hasLength[uint64](m, subject)

	return m.isZero
}

func (m Matcher) OnTestFailure(subject any, opts ...any) []string {
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
			"got     : len() == " + m.string,
		}
	case "string":
		return []string{
			"expected: <empty string>",
			"got     : " + opt.ValueAsString(subject, opts...),
		}
	default:
		return []string{
			"expected: <empty " + m._type + ">",
			"got     : " + m.method + "() == " + m.string,
		}
	}
}

func (m *Matcher) hasLen(v any) {
	switch got := v.(type) {
	case string:
		m.hasLength = true
		m.string = strconv.Itoa(len(got))
		m.isZero = len(got) == 0
		m.method = "len"
		m._type = "string"
		m.len = len(got)
	case nil:
		m.isNil = true
	default:

		typ := reflect.TypeOf(v)
		val := reflect.ValueOf(v)
		kind := val.Kind()
		if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && val.IsNil() {
			m.isNil = !m.TreatNilAsEmpty
			m._type = typ.Kind().String()
			m.hasLength = m.TreatNilAsEmpty
			m.string = "nil " + kind.String()
			m.isZero = true
			m.method = "len"
			return
		}

		var n int
		switch typ.Kind() {
		case reflect.Array:
			n = typ.Len()
		case reflect.Chan, reflect.Map, reflect.Slice:
			n = val.Len()
		default:
			return
		}

		m.hasLength = true
		m.string = strconv.Itoa(n)
		m.isZero = n == 0
		m.method = "len"
		m._type = typ.Kind().String()
		m.len = n
	}
}

func hasLength[T int | uint | int64 | uint64](m *Matcher, v any) {
	typ := fmt.Sprintf("%T", v)

	switch v := v.(type) {
	case interface{ Count() T }:
		m.hasLength = true
		m.string = fmt.Sprintf("%d", v.Count())
		m.isZero = v.Count() == 0
		m.method = "Count"
		m._type = typ
		if v.Count() < math.MaxInt {
			m.len = int(v.Count())
		}
		return

	case interface{ Len() T }:
		m.hasLength = true
		m.string = fmt.Sprintf("%d", v.Len())
		m.isZero = v.Len() == 0
		m.method = "Len"
		m._type = typ
		if v.Len() < math.MaxInt {
			m.len = int(v.Len())
		}
		return

	case interface{ Length() T }:
		m.hasLength = true
		m.string = fmt.Sprintf("%d", v.Length())
		m.isZero = v.Length() == 0
		m.method = "Length"
		m._type = typ
		if v.Length() < math.MaxInt {
			m.len = int(v.Length())
		}
		return
	}
}
