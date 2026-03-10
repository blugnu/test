package opt

import (
	"fmt"
)

// Name may be used to specify the name of the subject under test.
type Name string

// String returns the name as a string.
func (s Name) String() string {
	return string(s)
}

// GetName returns the first Name in opts.  A string is considered equivalent
// to a Name.  If no Name or string is found, it returns an empty string
// and false.
func GetName(opts []any) (Name, bool) {
	for _, opt := range opts {
		switch opt := opt.(type) {
		case Name:
			return opt, true
		case string:
			return Name(opt), true
		}
	}

	return "", false
}

// Namef returns a Name formed by formatting the specified string
// with the specified arguments.
//
// This is a convenience function equivalent to:
//
//	opt.Name(fmt.Sprintf(s, args...))
func Namef(s string, args ...any) Name {
	return Name(fmt.Sprintf(s, args...))
}

// WithName returns opts with the specified name added if it is not
// empty and opts does not already contain a Name.
//
// If name is empty, opts is returned unmodified.
func WithName(opts []any, name string) []any {
	if name == "" {
		return opts
	}

	if _, ok := GetName(opts); !ok {
		return append(opts, Name(name))
	}

	return opts
}

// WithNamef returns opts with a Name added, formed by formatting
// the specified arguments, if the resulting Name is not empty and opts
// does not already contain a Name.
//
// If format is empty, args are ignored and opts is returned unmodified.
func WithNamef(opts []any, format string, args ...any) []any {
	if format == "" {
		return opts
	}

	return WithName(opts, fmt.Sprintf(format, args...))
}
