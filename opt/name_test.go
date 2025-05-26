package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestHasName(t *testing.T) {
	With(t)

	Run("options contain a string", func() {
		result := opt.Name([]any{0, 1, "test", true})
		Expect(result).To(Equal("test"))
	})

	Run("options does not contain a string", func() {
		result := opt.Name([]any{0, 1, true})
		Expect(result).To(Equal(""))
	})

	Run("returns first of multiple strings", func() {
		result := opt.Name([]any{"test", "test2"})
		Expect(result).To(Equal("test"))
	})

	Run("options contains a Namef() result", func() {
		result := opt.Name([]any{0, 1, opt.Namef("test %d", 1), true})
		Expect(result).To(Equal("test 1"))
	})
}
