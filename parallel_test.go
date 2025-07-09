package test_test

import (
	"testing"
	"time"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

func TestIsParallel(t *testing.T) {
	With(t)

	Run(Test("not parallel", func() {
		Expect(IsParallel()).To(BeFalse())
		T().Parallel()
		Expect(IsParallel()).To(BeTrue())
	}))

	Run(Test("parallel", func() {
		T().Parallel()
		Expect(IsParallel()).To(BeTrue())
	}))

	Run(Test("example runner is always non-parallel", func() {
		test.Example()
		defer testframe.Pop()

		T().Parallel()
		Expect(IsParallel()).To(BeFalse())
	}))
}

func TestParallel(t *testing.T) {
	With(t)

	Run(Test("with single argument", func() {
		Expect(IsParallel(), "before").To(BeFalse())

		Parallel(T())

		Expect(IsParallel(), "after").To(BeTrue())
	}))

	Run(Test("with nil", func() {
		result := TestHelper(func() {
			Expect(IsParallel(), "before").To(BeFalse())

			Parallel(nil)

			test.Warning("should not have been reached")
		})
		result.ExpectInvalid("Parallel() cannot be called with nil")
	}))
}

func TestParallelTests(t *testing.T) {
	With(t)

	type testcase func()
	var (
		sleepFor10ms testcase = func() {
			time.Sleep(10 * time.Millisecond)
			Expect(true).To(BeTrue())
		}

		testcases = []testcase{
			sleepFor10ms,
			sleepFor10ms,
			sleepFor10ms,
		}
	)

	Run(Test("in series", func() {
		// this test establishes that running the testcases in series
		// takes longer than running them in parallel
		elapsed := StopWatch(func() {
			Run(Testcases(
				ForEach(func(tc testcase) {
					(tc)()
				}),
				Cases(testcases),
			))
		})
		Expect(elapsed).ToNot(BeLessThan(30 * time.Millisecond))
	}))

	Run(Test("in parallel using individual ParallelTest() calls", func() {
		elapsed := StopWatch(func() {
			Run(ParallelTest("test-1", sleepFor10ms))
			Run(ParallelTest("test-2", sleepFor10ms))
			Run(ParallelTest("test-3", sleepFor10ms))
		})
		Expect(elapsed).ToNot(BeGreaterThan(15 * time.Millisecond))
	}))

	Run(Test("in parallel using ParallelCases()", func() {
		elapsed := StopWatch(func() {
			Run(ParallelCases(
				ForEach(func(tc testcase) {
					(tc)()
				}),
				Cases(testcases),
			))
		})
		Expect(elapsed).ToNot(BeGreaterThan(15 * time.Millisecond))
	}))

	Run(Test("in parallel using Testcases(ParallelCase())", func() {
		elapsed := StopWatch(func() {
			Run(Testcases(
				ForEach(func(tc testcase) {
					(tc)()
				}),
				ParallelCase("test-1", sleepFor10ms),
				ParallelCase("test-2", sleepFor10ms),
				ParallelCase("test-3", sleepFor10ms),
			))
		})
		Expect(elapsed).ToNot(BeGreaterThan(15 * time.Millisecond))
	}))

	Run(Test("invalid use", func() {
		// note: we cannot use ParallelTest() to run tests that run
		// TestHelper() as this would require recording in parallel
		// which is not supported.
		//
		// We must run these tests in series, running parallel tests
		// within each test in series.

		Run(Test("Parallel() in a parallel test", func() {
			result := TestHelper(func() {
				Parallel(T()) // now this test is parallel

				// so this is now invalid
				Parallel(T())
			})
			result.ExpectInvalid(
				"Parallel() cannot be called from a parallel test",
			)
		}))

		Run(Test("ParallelTest() in a parallel test", func() {
			result := TestHelper(func() {
				Parallel(T()) // now this test is parallel

				// so this is now invalid
				Run(ParallelTest("parallel test", func() {}))
			})
			result.ExpectInvalid(
				"ParallelTest() cannot be run from a parallel test",
			)
		}))

		Run(Test("ParallelTest().Run() in a parallel test", func() {
			// create a parallel test runner, but don't run it
			para := ParallelTest("parallel test", func() {})

			result := TestHelper(func() {
				Parallel(T()) // now this test is parallel

				// so attempting to run the parallel test now is now invalid
				Run(para)
			})
			result.ExpectInvalid(
				"ParallelTest() cannot be run from a parallel test",
			)
		}))

		Run(Test("ParallelCases() in a parallel test", func() {
			result := TestHelper(func() {
				Parallel(T()) // now this test is parallel

				// so this is now invalid
				Run(ParallelCases(ForEach(func(struct{}) {}), Case("anon", struct{}{})))
			})
			result.ExpectInvalid(
				"ParallelCases() cannot be run from a parallel test",
			)
		}))
	}))
}
