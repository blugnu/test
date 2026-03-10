package internal

import "github.com/blugnu/test/opt"

// WithStringAsOnFailure checks the provided options for a string value. If one
// is found, it is removed from the options and replaced with an equivalent
// opt.OnFailure option.
//
// If no string is found, the options are returned unchanged.
//
// This allows a string to be provided as a shorthand for the opt.OnFailure
// option.
//
// Example:
//
//	expect.True(value, "value should be true") // equivalent to
//	expect.True(value, opt.OnFailure("value should be true"))
func WithStringAsOnFailure(opts []any) []any {
	if msg, opts, ok := opt.Extract[string](opts); ok {
		return append(opts, opt.OnFailure(msg))
	}

	return opts
}
