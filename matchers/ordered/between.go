package ordered

import (
	"cmp"
	"fmt"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type IsBetweenInitializer[T cmp.Ordered] struct {
	Value T
}

func (init IsBetweenInitializer[T]) And(v T) *IsBetween[T] {
	a := init.Value
	b := v

	if a > b {
		a, b = b, a
	}

	return &IsBetween[T]{
		Min: a,
		Max: b,
	}
}

type IsBetween[T cmp.Ordered] struct {
	Min, Max T

	// interval captures the interval used for the comparison
	interval opt.IntervalClosure

	// true if the interval is invalid
	invalidInterval bool
}

func (m *IsBetween[T]) Match(got T, opts ...any) bool {
	if interval, ok := opt.Get[opt.IntervalClosure](opts); ok {
		m.interval = interval
	}

	switch m.interval {
	case opt.IntervalOpen:
		return got > m.Min && got < m.Max
	case opt.IntervalOpenMin:
		return got > m.Min && got <= m.Max
	case opt.IntervalOpenMax:
		return got >= m.Min && got < m.Max
	case opt.IntervalClosed:
		return got >= m.Min && got <= m.Max
	}

	m.invalidInterval = true
	return false
}

func (m *IsBetween[T]) OnTestFailure(got T, opts ...any) []string {
	if m.invalidInterval {
		test.T().Helper()
		test.Invalid(fmt.Sprintf("unsupported option: %v", m.interval))
	}

	cond := "between "
	switch m.interval {
	case opt.IntervalOpen:
		cond += fmt.Sprintf("(%[1]v, %[2]v): %[1]v < x < %[2]v", m.Min, m.Max)
	case opt.IntervalOpenMin:
		cond += fmt.Sprintf("(%[1]v, %[2]v]: %[1]v < x <= %[2]v", m.Min, m.Max)
	case opt.IntervalOpenMax:
		cond += fmt.Sprintf("[%[1]v, %[2]v): %[1]v <= x < %[2]v", m.Min, m.Max)
	case opt.IntervalClosed:
		cond += fmt.Sprintf("[%[1]v, %[2]v]: %[1]v <= x <= %[2]v", m.Min, m.Max)
	}

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not " + cond
	}

	return []string{
		"expected: " + cond,
		"got     : " + opt.ValueAsString(got, opts...),
	}
}
