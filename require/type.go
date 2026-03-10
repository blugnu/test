package require

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// Type tests that a value is of an expected type. If the test passes,
// the value is returned as that type otherwise the test fails immediately
// without evaluating any further expectations.
//
// If a test does not use the returned value, consider using the [BeOfType]
// matcher instead, to avoid lint warnings about unused return values.
func Type[T any](got any, opts ...any) T {
	test.T().Helper()

	z, _ := expect.Type[T](got, append(opts, opt.Required())...)
	return z
}
