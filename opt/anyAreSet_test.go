package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestAnyAreSet(t *testing.T) {
	With(t)

	var result bool

	result = opt.AnyAreSet([]any{0, 1, "test", true}, 0, true)
	Expect(result).To(BeTrue())

	result = opt.AnyAreSet([]any{0, 1, true}, false)
	Expect(result).To(BeFalse())
}
