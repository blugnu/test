package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

func TestAnyOrder(t *testing.T) {
	With(t)

	result := opt.AnyOrder()

	if value, ok := expect.Type[opt.ExactOrder](result); ok {
		Expect(value).To(Equal(opt.ExactOrder(false)))
	}
}

func TestRequired(t *testing.T) {
	With(t)

	result := opt.Required()

	if value, ok := expect.Type[opt.IsRequired](result); ok {
		Expect(value).To(Equal(opt.IsRequired(true)))
	}
}

func TestNoStackTrace(t *testing.T) {
	With(t)

	result := opt.NoStackTrace()

	if value, ok := expect.Type[opt.StackTrace](result); ok {
		Expect(value).To(Equal(opt.StackTrace(false)))
	}
}
