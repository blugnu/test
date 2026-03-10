package internal

import (
	"sync"

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
	expectOpts := make([]any, 0, len(opts))
	matchOpts := make([]any, 0, len(opts))
	addName := sync.Once{}
	addRequired := sync.Once{}

	for _, v := range opts {
		switch v.(type) {
		case string, opt.Name:
			addName.Do(func() {
				switch s := v.(type) {
				case string:
					expectOpts = append(expectOpts, opt.Name(s))
				case opt.Name:
					expectOpts = append(expectOpts, s)
				}
			})

		case opt.IsRequired:
			addRequired.Do(func() {
				expectOpts = append(expectOpts, v)
			})

		default:
			matchOpts = append(matchOpts, v)
		}
	}

	return expectOpts, matchOpts
}
