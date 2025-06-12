package opt

import (
	"fmt"
	"reflect"
)

// FormatString formats a string according to the provided options.
//
// This function should be used in TestFailure methods to format expected and
// actual values consistently with other matcher failure reports.
//
// Supported options:
//
//	opt.QuotedStrings(bool)    // determines whether strings should be quoted
//	                           // (default is true).
//	                           //
//	                           // Has no effect if the value is not a string.
func ValueAsString(v any, opts ...any) string {
	if v == nil {
		return "nil"
	}

	isString := reflect.TypeOf(v).Kind() == reflect.String

	spec := "%v"
	switch {
	case isString && IsSet(opts, QuotedStrings(false)):
		spec = "%v"
	case isString:
		spec = "%q"
	case IsSet(opts, AsDeclaration(true)):
		spec = "%#v"
	}

	return fmt.Sprintf(spec, v)
}
