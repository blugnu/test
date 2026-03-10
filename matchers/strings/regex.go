package strings

import (
	"regexp"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type RegExMatch struct {
	Expected *regexp.Regexp
}

func (m RegExMatch) Match(got string, opts ...any) bool {
	return m.Expected.Match([]byte(got))
}

func (m RegExMatch) OnTestFailure(got string, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{
			"expected: string with no match for: " + report.Value(m.Expected.String(), opts...),
			"got     : " + report.Value(got, opts...),
			"matched : " + report.Value(string(m.Expected.Find([]byte(got))), opts...),
		}
	}
	return []string{
		"expected: string containing match for: " + report.Value(m.Expected.String(), opts...),
		"got     : " + report.Value(got, opts...),
	}
}
