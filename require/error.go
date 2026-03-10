package require

import (
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"

	"github.com/blugnu/test"
)

// Error tests that an error satisfies [errors.As] for a specified
// error type E.  If the test passes, the value is returned as that type.
// If the test fails, the test fails immediately without evaluating any
// further expectations.
//
// If a test does not use the returned value, consider using the [test.BeError]
// matcher with [test.Require] instead, to avoid lint warnings about
// unused return values:
//
//	Require(err).Should(BeError[E]())
//
// If a failure to match the expected error type is not a fatal error for the
// test, consider using [expect.ErrorAs] instead:
//
//	e, ok := expect.Error[E](err)
//	if !ok {
//	    // handle the error not being of type E
//	}
func Error[E error](got error, opts ...any) E {
	test.T().Helper()

	target, _ := expect.ErrorAs[E](got, append(opts, opt.Required())...)

	return target
}
