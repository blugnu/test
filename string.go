package test

import (
	"strings"
)

type StringMatcher struct {
	Expecting[string]
	match func(string, string) bool
}

func (m StringMatcher) Match(got string, opts ...any) bool {
	return m.match(got, m.value)
}

func ContainString(expected string) StringMatcher {
	if expected == "" {
		// TODO: figure a way to offter suggestions that take account of whether
		//  To() or ToNot() is involved.  The suggested alternative here does not make
		//  sense if used in ToNot()
		panic("ContainString(<empty string>) is invalid: did you mean Expect(<string>).IsEmpty()")
	}

	return StringMatcher{Expecting[string]{value: expected},
		func(got, expected string) bool {
			return strings.Contains(got, expected)
		}}
}
