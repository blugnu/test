package require

import (
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// Slice returns an expectation for the provided slice with opt.Required() added
// to the options. This ensures that any assertion that fails using the returned
// expectation will halt further execution of the test.
//
// The returned expectation provides the same assertions as [expect.Slice].
func Slice[T any](got []T, opts ...any) expect.SliceTest[T] {
	return expect.Slice(got, append(opts, opt.Required())...)
}
