package test

type TestScenario struct {
	Scenario string
	Act      func()
	Assert   func(R)
}

func RunTestScenarios(scns []TestScenario) {
	T().Helper()

	fn := func(tc TestScenario) {
		if tc.Act == nil {
			panic("RunTestScenario: act phase is nil")
		}

		T().Helper()

		test := Test(tc.Act)
		if tc.Assert == nil {
			test.Assert(nil)
		} else {
			tc.Assert(test)
		}
	}
	RunScenarios(fn, scns)
}
