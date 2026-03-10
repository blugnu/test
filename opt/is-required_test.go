package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestIsRequired_String(t *testing.T) {
	With(t)

	sut := opt.IsRequired(true)
	Expect(sut.String()).To(Equal("opt.IsRequired(true)"))

	sut = opt.IsRequired(false)
	Expect(sut.String()).To(Equal("opt.IsRequired(false)"))
}

func TestIsRequired(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "option is present and true",
			Act: func() {
				opts := []any{
					opt.IsRequired(true),
				}

				value, ok := opt.Get[opt.IsRequired](opts)

				Expect(ok).To(BeTrue())
				Expect(value).To(Equal(opt.IsRequired(true)))
			},
		},
		{Scenario: "option is present and false",
			Act: func() {
				opts := []any{
					opt.IsRequired(false),
				}

				value, ok := opt.Get[opt.IsRequired](opts)

				Expect(ok).To(BeTrue())
				Expect(value).To(Equal(opt.IsRequired(false)))
			},
		},
		{Scenario: "option is absent",
			Act: func() {
				opts := []any{}

				_, ok := opt.Get[opt.IsRequired](opts)
				Expect(ok).To(BeFalse())
			},
		},
	}...))
}
