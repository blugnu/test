package strings

import (
	"strings"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type ContainsMatch struct {
	Expected string
}

func (m ContainsMatch) Match(got string, opts ...any) bool {
	if opt.IsSet(opts, opt.CaseInsensitive) {
		return strings.Contains(strings.ToLower(got), strings.ToLower(m.Expected))
	}

	return strings.Contains(got, m.Expected)
}

func (m ContainsMatch) OnTestFailure(got string, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		offset := strings.Index(got, m.Expected)
		if !opt.IsSet(opts, opt.UnquotedStrings) {
			offset += 1
		}
		pad := strings.Repeat(" ", offset)

		return []string{
			"expected: string not containing: " + report.Value(m.Expected, opts...),
			"got     : " + report.Value(got, opts...),
			"          " + pad + strings.Repeat("^", len(m.Expected)),
		}
	}
	return []string{
		"expected: string containing: " + report.Value(m.Expected, opts...),
		"got     : " + report.Value(got, opts...),
	}
}
