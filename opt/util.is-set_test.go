package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestIsSet(t *testing.T) {
	With(t)

	sut := []any{opt.ExactOrder(true)}

	var result bool

	result = opt.IsSet(sut, opt.ExactOrder(true))
	Expect(result).To(BeTrue())

	result = opt.IsSet(sut, opt.ExactOrder(false))
	Expect(result).To(BeFalse())
}
