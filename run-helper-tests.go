package test

import (
	"fmt"

	"github.com/blugnu/test/test"
)

type HelperScenario struct {
	Scenario string
	Act      func()
	Assert   func(*R)
	Debug    bool
	Skip     bool
}

type helpertestRunner struct {
	scenarios []HelperScenario
}

func HelperTests(scns ...HelperScenario) helpertestRunner {
	return helpertestRunner{
		scenarios: scns,
	}
}

func (r helpertestRunner) Run() {
	T().Helper()

	Run(Testcases(
		For(func(name string, scn HelperScenario) {
			T().Helper()

			if scn.Act == nil {
				test.Invalid("no Act function defined")
			}

			result := TestHelper(scn.Act)
			if scn.Assert == nil {
				result.Expect(TestPassed)
				return
			}

			scn.Assert(&result)
			if !result.checked {
				test.Warning(fmt.Sprintf("HelperTests() result not tested (scenario: %s)", name))
			}
		}),
		Cases(r.scenarios),
	))
}
