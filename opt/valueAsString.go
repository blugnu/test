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

	spec := "%q"
	if reflect.TypeOf(v).Kind() != reflect.String || IsSet(opts, QuotedStrings(false)) {
		spec = "%v"
	}

	return fmt.Sprintf(spec, v)
}
