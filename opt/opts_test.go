package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestAnyOrder(t *testing.T) {
	With(t)

	result := opt.AnyOrder()

	if value, ok := ExpectType[opt.ExactOrder](result); ok {
		Expect(value).To(Equal(opt.ExactOrder(false)))
	}
}
