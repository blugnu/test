package opt

// CaseSensitive may be used to indicate that the contents of strings
// should be compared in a case-sensitive manner (or not).
//
// Most string-based matchers will compare strings in a case-sensitive manner
// by default, which may be overridden by using opt.CaseSensitive(false) if
// supported by the matcher.
type CaseSensitive bool

// ExactOrder may be used to indicate that the order of elements in a
// collection is significant (or not).
type ExactOrder bool

// IgnoreReport may be used to indicate that the contents of any test report are not
// significant when testing the result of testing a test, i.e. R.Expect().
//
// This is useful when the test report is not significant to the test outcome.
type IgnoreReport bool

// NoPanic is an internal option used as a sentinel recover value by the panic
// testing mechanism to signal that a panic is NOT expected to occur
type NoPanicExpected bool

// PrefixInlineWithFirstItem may be used to indicate that the first item
// in a collection should be output on the same line as any prefix when
// appending to a test report
type PrefixInlineWithFirstItem bool

// QuotedStrings may be used to indicate that the contents of string values
// should not be quoted when reported in test failures.
//
// In matchers that support this option, string values in test failure reports
// are quoted by default. This option may be used to override that behavior,
// typically by using the opt.UnquotedStrings() convenience function.
type QuotedStrings bool

// ToNotMatch is set automatically when a matcher is invoked in a ToNot() test.
// The option is not set when the matcher is invoked in a To() test.
// Matchers should test for this option to phrase the test report correctly
// when the test fails.
type ToNotMatch bool

// AnyOrder is a convenience function that returns ExactOrder(false).
func AnyOrder() ExactOrder {
	return ExactOrder(false)
}

// UnquotedStrings is a convenience function that returns QuotedStrings(false).
func UnquotedStrings() QuotedStrings {
	return QuotedStrings(false)
}
