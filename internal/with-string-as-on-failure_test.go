package internal_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
)

func TestWithStringAsOnFailure(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "with string option",
			Act: func() {
				// arrange
				opts := []any{"custom failure message", opt.UnquotedStrings}

				// act
				newOpts := internal.WithStringAsOnFailure(opts)

				// assert
				Expect(newOpts).Should(HaveLen(2))
				Expect(newOpts).To(ContainItem(opts[1]))

				newOpt0, ok := opt.Get[opt.FailReporter](newOpts)
				Expect(ok).To(BeTrue())

				report := newOpt0.OnTestFailure()
				Expect(report).To(EqualSlice([]string{"custom failure message"}))
			},
		},
		{Scenario: "without string option",
			Act: func() {
				// arrange
				opts := []any{opt.UnquotedStrings}

				// act
				newOpts := internal.WithStringAsOnFailure(opts)

				// assert
				Expect(newOpts).Should(HaveLen(1))
				Expect(newOpts).To(ContainItem(opts[0]))
			},
		},
	}...))
}
