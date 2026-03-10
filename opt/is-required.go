package opt

import "fmt"

// IsRequired may be used to indicate that an expectation is required to pass; if
// the expectation is not met the current test is failed and execution continues
// with the *next* test.  No further expectations in the current test will be
// evaluated.
//
// see also: Require()
type IsRequired bool

// String returns a string representation of the IsRequired option.
func (r IsRequired) String() string {
	return fmt.Sprintf("opt.IsRequired(%t)", bool(r))
}

// Required is a convenience function that returns IsRequired(true)
func Required() IsRequired {
	return IsRequired(true)
}
