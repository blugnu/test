package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestExtract(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "option is present",
			Act: func() {
				opts := []any{
					opt.Name("test subject"),
					opt.IsRequired(true),
				}

				value, remain, ok := opt.Extract[opt.Name](opts)

				Expect(ok).To(BeTrue())
				Expect(value).To(Equal(opt.Name("test subject")))
				Expect(remain).To(EqualSlice([]any{opt.IsRequired(true)}))
			},
		},
		{Scenario: "option is absent",
			Act: func() {
				opts := []any{
					opt.IsRequired(true),
				}

				_, ok := opt.Get[opt.Name](opts)
				Expect(ok).To(BeFalse())
			},
		},
	}...))
}
