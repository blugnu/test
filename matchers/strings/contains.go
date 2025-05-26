package strings

import (
	"fmt"
	"strings"

	"github.com/blugnu/test/opt"
)

type ContainsMatch struct {
	Expected string
}

func (m ContainsMatch) Match(got string, opts ...any) bool {
	if opt.IsSet(opts, opt.CaseSensitive(false)) {
		return strings.Contains(strings.ToLower(string(got)), strings.ToLower(string(m.Expected)))
	}

	return strings.Contains(string(got), string(m.Expected))
}

func (m ContainsMatch) OnTestFailure(got string, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{
			"expected: string not containing: " + opt.ValueAsString(m.Expected, opts...),
			"got     : " + opt.ValueAsString(got, opts...),
			fmt.Sprintf("             %s", strings.Repeat("^", len(m.Expected))),
		}
	}
	return []string{
		"expected: string containing: " + opt.ValueAsString(m.Expected, opts...),
		"got     : " + opt.ValueAsString(got, opts...),
	}
}
