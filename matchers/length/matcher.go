package length

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/test"
)

type Matcher struct {
	Length int

	got     any
	gotLen  int
	isValid bool
}

func (m *Matcher) Match(got any, opts ...any) bool {
	handleN := func(n int) bool {
		m.gotLen = n
		m.isValid = true
		return n == m.Length
	}

	m.got = got

	val := reflect.ValueOf(got)
	kind := val.Kind()
	if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && val.IsNil() {
		return handleN(0)
	}

	typ := reflect.TypeOf(got)
	if typ == nil {
		return false
	}

	switch typ.Kind() { //nolint: exhaustive // dealing only with types that support len()
	case reflect.Array:
		return handleN(typ.Len())

	case reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return handleN(val.Len())
	}

	return false
}

func (m *Matcher) OnTestFailure(subject any, opts ...any) []string {
	if !m.isValid {
		test.T().Helper()
		test.Invalid(fmt.Sprintf("length.Matcher: requires a value that is a string, slice, channel, or map: got %T", m.got))
	}

	return []string{
		fmt.Sprintf("expected: len() == %d", m.Length),
		fmt.Sprintf("got     : len() == %d", m.gotLen),
	}
}
