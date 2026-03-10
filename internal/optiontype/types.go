package optiontype

// CaseSensitivity is an option supported by matchers that perform
// string comparisons to determine whether string comparisons should be
// case-sensitive (true) or case-insensitive (false).
type CaseSensitivity bool

// ValueFormat is an option supported by report.ValueAsString that
// may be used to format values in different ways. The possible values are:
//
//	ValueAsDeclaration // values are formatted using %#v (i.e. Go-syntax representation)
//	ValueAsString      // values are formatted using %q (strings) or %v (all other types); if [PlainStrings] is set, strings are also formatted using %v
//
// Support for this option depends on a matcher using report.ValueAsString.
// Check the documentation of specific matchers for details.
type ValueFormat int
