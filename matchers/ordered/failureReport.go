package ordered

import "github.com/blugnu/test/opt"

func failureReport[T any](cond string, expected, got T, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not " + cond
	}
	return []string{
		"expected: " + cond + " " + opt.ValueAsString(expected, opts...),
		"got     : " + opt.ValueAsString(got, opts...),
	}
}
