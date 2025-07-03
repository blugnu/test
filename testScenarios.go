package test

import "github.com/blugnu/test/test"

type TestScenario struct {
	Scenario string
	Debug    bool
	Skip     bool
	Act      func()
	Assert   func(*R)
}

func RunTestScenarios(scns []TestScenario) {
	T().Helper()

	skipped := 0

	fn := func(tc *TestScenario, _ int) {
		t := T()
		t.Helper()

		if tc.Skip {
			skipped++
			t.SkipNow()
			return
		}

		if tc.Act == nil {
			test.Invalid("test.RunTestScenarios: no Act function defined")
		}

		result := TestHelper(tc.Act)
		if tc.Assert == nil {
			result.Expect(TestPassed)
			return
		}

		tc.Assert(&result)
		if !result.checked {
			test.Invalid("test.RunTestScenarios: result not tested; *R.Expect(...) or *R.ExpectInvalid(...) must be called")
		}
	}

	// if any of the scenarios are marked as Debug then only run those. This is useful
	// for debugging individual scenarios without having to run all of them.
	//
	// HACK: a bit of a hack, but it works for now.
	debug := make([]TestScenario, 0, len(scns))
	for _, scn := range scns {
		if scn.Debug {
			debug = append(debug, scn)
		}
	}
	if len(debug) > 0 {
		scns = debug
	}

	RunScenarios(fn, scns)

	if len(debug) > 0 {
		T().Errorf("WARNING: %d tests were run with Debug: true", len(debug))
	}
	if skipped > 0 {
		T().Errorf("WARNING: %d tests were skipped", skipped)
	}
}
