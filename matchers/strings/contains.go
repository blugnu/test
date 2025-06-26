package strings

import (
	"strings"

	"github.com/blugnu/test/opt"
)

type ContainsMatch struct {
	Expected string
}

func (m ContainsMatch) Match(got string, opts ...any) bool {
	if opt.IsSet(opts, opt.CaseSensitive(false)) {
		return strings.Contains(strings.ToLower(got), strings.ToLower(m.Expected))
	}

	return strings.Contains(got, m.Expected)
}

func (m ContainsMatch) OnTestFailure(got string, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		offset := strings.Index(got, m.Expected)
		if !opt.IsSet(opts, opt.QuotedStrings(false)) {
			offset += 1
		}
		pad := strings.Repeat(" ", offset)

		return []string{
			"expected: string not containing: " + opt.ValueAsString(m.Expected, opts...),
			"got     : " + opt.ValueAsString(got, opts...),
			"          " + pad + strings.Repeat("^", len(m.Expected)),
		}
	}
	return []string{
		"expected: string containing: " + opt.ValueAsString(m.Expected, opts...),
		"got     : " + opt.ValueAsString(got, opts...),
	}
}
