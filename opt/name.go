package opt

import "fmt"

// Name returns the first of any string values in the specified
// options.
func Name(opts []any) string {
	for _, opt := range opts {
		if s, ok := opt.(string); ok {
			return s
		}
	}

	return ""
}

// Namef returns a string formed by formatting the specified string
// with the specified arguments.
//
// This is identical to fmt.Sprintf() and is provided to make the
// intent of the code slightly shorter and (more importantly) clearer
// when naming a test expectation:
//
//	Expect(got, opt.Namef("result of case %d", i)).To(Equal("expected result"))
//
// vs
//
//	Expect(got, fmt.Sprintf("result of case %d", i)).To(Equal("expected result"))
func Namef(s string, args ...any) string {
	return fmt.Sprintf(s, args...)
}
