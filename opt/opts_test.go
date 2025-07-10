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

func TestRequired(t *testing.T) {
	With(t)

	result := opt.Required()

	if value, ok := ExpectType[opt.IsRequired](result); ok {
		Expect(value).To(Equal(opt.IsRequired(true)))
	}
}

func TestNoStackTrace(t *testing.T) {
	With(t)

	result := opt.NoStackTrace()

	if value, ok := ExpectType[opt.StackTrace](result); ok {
		Expect(value).To(Equal(opt.StackTrace(false)))
	}
}

func TestUnquotedStrings(t *testing.T) {
	With(t)

	result := opt.UnquotedStrings()

	if value, ok := ExpectType[opt.QuotedStrings](result); ok {
		Expect(value).To(Equal(opt.QuotedStrings(false)))
	}
}
