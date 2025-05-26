package test

type TestScenario struct {
	Scenario string
	Debug    bool
	Skip     bool
	Act      func()
	Assert   func(*R)
}

func RunTestScenarios(scns []TestScenario) {
	T().Helper()

	fn := func(tc *TestScenario, _ int) {
		t := T()
		t.Helper()

		if tc.Skip {
			t.SkipNow()
			return
		}

		if tc.Act == nil {
			invalidTest("test.RunTestScenarios: no Act function defined")
			return
		}

		result := Test(tc.Act)
		if tc.Assert == nil {
			result.Expect(TestPassed)
			return
		}
		tc.Assert(&result)
		if !result.checked {
			invalidTest("test.RunTestScenarios: result not tested; *R.Expect(...) or *R.ExpectInvalid(...) must be called")
			return
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
}
