package require

import (
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// Map creates a MapTest for the provided map and options.  If any assertion
// using the returned MapTest fails, the current test is halted immediately.
//
// The returned expectation provides the same assertions as [expect.Map].
func Map[K comparable, V any](got map[K]V, opts ...any) expect.MapTest[K, V] {
	return expect.Map(got, append(opts, opt.Required())...)
}
