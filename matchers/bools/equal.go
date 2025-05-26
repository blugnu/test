package bools

import (
	"fmt"

	"github.com/blugnu/test/opt"
)

// implements the Matcher interface for testing expected bool values.
type BooleanMatcher struct {
	Expected bool
}

// Match returns true if the actual value matches the expected value.
func (m BooleanMatcher) Match(got bool, _ ...any) bool {
	return m.Expected == got
}

func (m BooleanMatcher) OnTestFailure(got bool, opts ...any) string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return fmt.Sprintf("did not expect %v", m.Expected)
	}
	return fmt.Sprintf("expected %v, got %v", m.Expected, got)
}
