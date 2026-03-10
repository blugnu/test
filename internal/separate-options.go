package internal

import (
	"github.com/blugnu/test/opt"
)

// SeparateOptions separates opts into those that are supported by an expectation
// and those supported by a matcher.
//
// Options that are handled as expectation options are:
//
//	string, opt.Name, opt.IsRequired
//
// Any string option is replaced with an opt.Name option in the returned expectation
// options.
//
// All other options are assumed to be matcher options.
//
// Expectation options are de-duped; only the first of each supported option type
// instance is included in the returned expectation options; additional instances
// are ignored and discarded. String and opt.Name options are considered equivalent
// for de-duplication purposes.
//
// If matcher option types are specified more than once, all instances are included
// in the returned matcher options; handling of multiple instances of options is
// left to individual matchers.
func SeparateOptions(opts []any) ([]any, []any) {
	var (
		expectOpts    = make([]any, 0, len(opts))
		matchOpts     = make([]any, 0, len(opts))
		nameAdded     bool
		requiredAdded bool
	)

	for _, v := range opts {
		switch v.(type) {
		case string, opt.Name:
			if nameAdded {
				continue
			}

			switch s := v.(type) {
			case string:
				expectOpts = append(expectOpts, opt.Name(s))
			case opt.Name:
				expectOpts = append(expectOpts, s)
			}
			nameAdded = true

		case opt.IsRequired:
			if requiredAdded {
				continue
			}

			expectOpts = append(expectOpts, v)
			requiredAdded = true

		default:
			matchOpts = append(matchOpts, v)
		}
	}

	return expectOpts, matchOpts
}
