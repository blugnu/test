package opt

import "github.com/blugnu/test/internal/optiontype"

// AsDeclaration is an option supported by report.ValueAsString that may be used
// to format values as a declaration, i.e. with the type name and value
// included in the output.
type AsDeclaration bool

// CaseInsensitive directs compatible matchers to compare strings in a
// case-insensitive manner. Most matchers that support this option
// will use case-sensitive comparisons by default.
//
// Check the documentation of specific matchers for details.
//
// # Equivalent to:
//
//	opt.CaseIsSignificant(false)
const CaseInsensitive = optiontype.CaseSensitivity(false)

// CaseSensitive directs compatible matchers to compare strings in a
// case-sensitive manner. Most matchers that support this option
// will already use case-sensitive comparisons by default.
//
// Check the documentation of specific matchers for details.
//
// # Equivalent to:
//
//	opt.CaseIsSignificant(true)
const CaseSensitive = optiontype.CaseSensitivity(true)

// ExactOrder may be used to indicate that the order of elements in a
// collection is significant (or not).
type ExactOrder bool

const (
	ValueAsDeclaration optiontype.ValueFormat = iota // values are formatted using %#v (i.e. Go-syntax representation)
	ValueAsString                                    // values are formatted using %q (strings) or %v (all other types); if [PlainStrings] is set, strings are also formatted using %v
)

// IgnoreReport may be used to indicate that the contents of any test report are not
// significant when testing the result of testing a test, i.e. R.Expect().
//
// This is useful when the test report is not significant to the test outcome.
type IgnoreReport bool

const AnyReport = IgnoreReport(true)

// ItemMatch may be used to indicate how items in a collection should be
// matched when comparing collections.
// The possible values are:
//
//	opt.AllItems       // all items in the expected collection must be present
//	opt.AnyItem        // any item in the expected collection may match
//	opt.Subset         // all items in the expected collection must match; there may be extra items
//	opt.Exact          // all items in the expected collection must match; no extra items allowed
type ItemMatch int

const (
	// AllItems indicates that all items in the expected collection
	// must be present in the actual collection.
	AllItems ItemMatch = iota + 1

	// AnyItem indicates that any item in the expected collection
	// may match an item in the actual collection.
	AnyItem

	// Subset indicates that all items in the expected collection
	// must match items in the actual collection; there may be extra items
	Subset

	// Exact indicates that all items in the expected collection
	// must match items in the actual collection; no extra items allowed
	Exact
)

// MaxItems is an option that may be used to limit the number of items
// reported in a collection. It is used by [report.AppendSlice] and
// [report.AppendMap] to limit the number of items reported in a slice or map.
// If the number of items exceeds the limit, the remaining items are reported
// as "... and N more".
type MaxItems int

// NoPanicExpected is an internal option used as a sentinel recover value by
// the panic testing mechanism to signal that a panic is NOT expected to occur
type NoPanicExpected bool

type ReportType bool

const (
	NoTypeNames      = ReportType(false) // directs that type names should not be included in the output
	IncludeTypeNames = ReportType(true)  // directs that type names should be included in the output (default)
)

// StackTrace may be used to indicate that a stack trace should be included
// in the test report when a test fails.  Where a stack trace is supported
// it is generally included by default, so this option may be used to
// disable it by including opt.StackTrace(false) or opt.NoStackTrace() in
// the options.
type StackTrace bool

const (
	QuotedStrings   = quoteStrings(true)  // directs that strings should be formatted with quotes (default)
	UnquotedStrings = quoteStrings(false) // directs that strings should be formatted without quotes
)

// ToNotMatch is set internally when a matcher is invoked in a ToNot() or
// ShouldNot() test.
//
// Matchers should test for this option to phrase the test report correctly
// when the test fails (and/or to modify the behaviour of the matcher when
// matching, if appropriate; it usually isn't).
type ToNotMatch bool

// MARK: Convenience Functions

// AnyOrder is a convenience function that returns ExactOrder(false)
func AnyOrder() ExactOrder {
	return ExactOrder(false)
}

// NoStackTrace is a convenience function that returns StackTrace(false)
func NoStackTrace() StackTrace {
	return StackTrace(false)
}

// quoteStrings is applied by [ValueAsString] when formatting string values
// to quote strings (true/default) or not (false).
//
// This option type is not exported; it is specified using the constants
// [QuotedStrings] and [UnquotedStrings].
type quoteStrings bool
