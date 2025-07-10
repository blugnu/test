package opt

// AsDeclaration is an option supported by opt.ValueAsString that may be used
// to format values as a declaration, i.e. with the type name and value
// included in the output.
type AsDeclaration bool

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

// IsRequired may be used to indicate that an expectation is required to pass; if
// the expectation is not met the current test is failed and execution continues
// with the *next* test.  No further expectations in the current test will be
// evaluated.
//
// see also: Require()
type IsRequired bool

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

// StackTrace may be used to indicate that a stack trace should be included
// in the test report when a test fails.  Where a stack trace is supported
// it is generally included by default, so this option may be used to
// disable it by including opt.StackTrace(false) or opt.NoStackTrace() in
// the options.
type StackTrace bool

// ToNotMatch is set internall when a matcher is invoked in a ToNot() or
// ShouldNot() test.
//
// Matchers should test for this option to phrase the test report correctly
// when the test fails (and/or to modify the behaviour of the matcher when
// matching, if appropriate; it usually isn't).
type ToNotMatch bool

// AnyOrder is a convenience function that returns ExactOrder(false)
func AnyOrder() ExactOrder {
	return ExactOrder(false)
}

// NoStackTrace is a convenience function that returns StackTrace(false)
func NoStackTrace() StackTrace {
	return StackTrace(false)
}

// Required is a convenience function that returns IsRequired(true)
func Required() IsRequired {
	return IsRequired(true)
}

// UnquotedStrings is a convenience function that returns QuotedStrings(false)
func UnquotedStrings() QuotedStrings {
	return QuotedStrings(false)
}
