package ordered

import (
	"cmp"
	"fmt"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type Comparison int

const orEqual int = 1

const (
	LessThan Comparison = iota
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
)

type RelativeMatcher[T cmp.Ordered] struct {
	Expected   T
	Comparison Comparison
}

func (m RelativeMatcher[T]) OrEqual() RelativeMatcher[T] {
	// NOTE: this relies on specific ordering of comparison constants
	m.Comparison += Comparison(orEqual)
	return m
}

func (m RelativeMatcher[T]) Match(got T, opts ...any) bool {
	if cmp, ok := opt.Get[func(T, T) bool](opts); ok {
		return cmp(got, m.Expected)
	}

	switch m.Comparison {
	case LessThan:
		return got < m.Expected
	case LessThanOrEqual:
		return got <= m.Expected
	case GreaterThan:
		return got > m.Expected
	case GreaterThanOrEqual:
		return got >= m.Expected
	}

	return false
}

func (m RelativeMatcher[T]) OnTestFailure(got T, opts ...any) []string {
	cond, ok := map[Comparison]string{
		LessThan:           "less than",
		LessThanOrEqual:    "less than or equal to",
		GreaterThan:        "greater than",
		GreaterThanOrEqual: "greater than or equal to",
	}[m.Comparison]

	if !ok {
		test.T().Helper()
		test.Invalid(fmt.Sprintf("unknown comparison type: Comparison(%d)", m.Comparison))
	}

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not " + cond
	}

	return []string{
		"expected: " + cond + " " + opt.ValueAsString(m.Expected, opts...),
		"got     : " + opt.ValueAsString(got, opts...),
	}
}
