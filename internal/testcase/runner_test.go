package testcase_test

import (
	"testing"
	"time"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testcase"
)

func TestNewRunner(t *testing.T) {
	With(t)

	Run(Test("nil executor", func() {
		var runner testcase.Runner[int]

		result := TestHelper(func() {
			runner = testcase.NewRunner[int](nil)
		})

		result.ExpectInvalid("test executor cannot be nil")

		Expect(runner.TestingT).IsNil()
		Expect(runner.TestExecutor).IsNil()
		Expect(runner.CaseControllers).Should(BeEmptyOrNil())
	}))

	Run(Test("no initial cases", func() {
		runner := testcase.NewRunner(testcase.Executor[int]{})

		Expect(runner.TestingT).IsNotNil()
		Expect(runner.TestExecutor).IsNotNil()
		Expect(runner.CaseControllers).Should(BeEmptyOrNil())
	}))
}

func TestRunner_Run(t *testing.T) {
	With(t)

	Run(Test("run with no cases", func() {
		result := TestHelper(func() {
			runner := testcase.NewRunner(testcase.Executor[int]{})
			runner.Run()
		})

		result.ExpectInvalid("no test cases provided")
	}))

	Run(Test("parallel cases", func() {
		elapsed := StopWatch(func() {
			runner := testcase.NewRunner(ForEach(func(tc int) {
				time.Sleep(10 * time.Millisecond) // simulate work
			}))
			runner.AddCase("parallel", 1, testcase.Parallel)
			runner.AddCase("parallel", 2, testcase.Parallel)
			runner.AddCase("parallel", 3, testcase.Parallel)
			runner.AddCase("parallel", 4, testcase.Parallel)

			runner.Run()
		})
		Expect(elapsed).To(BeLessThan(12 * time.Millisecond))
	}))

	Run(Test("some skipped cases", func() {
		evals := []string{}
		result := TestHelper(func() {
			runner := testcase.NewRunner(For(func(name string, tc int) {
				evals = append(evals, name)
			}))
			runner.AddCase("skip-1", 1, testcase.Skip)
			runner.AddCase("eval-2", 2)
			runner.AddCase("skip-3", 3, testcase.Skip)

			runner.Run()
		})

		result.ExpectWarning("2 of 3 cases were skipped")

		Expect(evals).To(EqualSlice([]string{
			"eval-2",
		}))
	}))

	Run(Test("all cases skipped", func() {
		evals := []string{}
		result := TestHelper(func() {
			runner := testcase.NewRunner(For(func(name string, tc int) {
				evals = append(evals, name)
			}))

			runner.AddCase("skip-1", 1, testcase.Skip)
			runner.AddCase("skip-2", 2, testcase.Skip)

			runner.Run()
		})

		result.ExpectWarning("all cases were skipped")

		Expect(evals).Should(BeEmptyOrNil())
	}))

	Run(Test("some debug cases", func() {
		result := TestHelper(func() {
			var executed int
			runner := testcase.NewRunner(ForEach(func(tc int) {
				executed++
			}))

			runner.AddCase("test", 1)
			runner.AddCase("debug", 2, testcase.Debug)
			runner.AddCase("skip", 3, testcase.Skip)
			runner.AddCase("parallel", 4, testcase.Parallel)

			runner.Run()

			Expect(executed).To(Equal(1))
		})
		result.ExpectWarning("only 1 of 4 cases were evaluated (debug mode)")
	}))

	Run(Test("all debug cases", func() {
		result := TestHelper(func() {
			var executed int
			runner := testcase.NewRunner(ForEach(func(tc int) {
				executed++
			}))

			runner.AddCase("debug-1", 1, testcase.Debug)
			runner.AddCase("debug-2", 2, testcase.Debug)

			runner.Run()

			Expect(executed).To(Equal(2))
		})

		result.Expect(TestPassed)
	}))
}
