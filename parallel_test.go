package test

import (
	"errors"
	"testing"
	"time"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

func TestIsParallel(t *testing.T) {
	With(t)

	Run("not parallel", func() {
		Expect(IsParallel()).To(BeFalse())
		Parallel()
		Expect(IsParallel()).To(BeTrue())
	})

	Run("parallel", func() {
		Parallel()
		Expect(IsParallel()).To(BeTrue())
	})

	Run("example runner is always non-parallel", func() {
		test.Example()
		defer testframe.Pop()

		Parallel()
		Expect(IsParallel()).To(BeFalse())
	})
}

func TestParallel(t *testing.T) {
	With(t)

	t.Run("with single argument", func(st *testing.T) {
		With(st)
		Expect(IsParallel(), "before").To(BeFalse())
		Parallel(st)
		Expect(IsParallel(), "after").To(BeTrue())
	})

	Run("with no arguments", func() {
		Expect(IsParallel()).To(BeFalse())
		Parallel()
		Expect(IsParallel()).To(BeTrue())
	})

	Run("multiple arguments, with test frame", func() {
		result := Test(func() {
			Parallel(t, t)
		})
		result.ExpectInvalid(
			"ERROR: invalid argument",
			"Parallel() must be called with 0 or 1 test runner arguments",
		)
	})
}

func TestParallel_MultipleArgumentsNoTestFrame(t *testing.T) {
	defer func() {
		r := recover()
		err, ok := r.(error)
		if !ok {
			t.Errorf("expected panic with error, got: %v", r)
			return
		}
		if !errors.Is(err, ErrInvalidArgument) {
			t.Errorf("\nexpected panic with ErrInvalidArgument\ngot: %v", r)
		}
	}()

	Parallel(t, t)
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

	Run("in series", func() {
		elapsed := RunWithTimer(func() {
			RunScenarios(
				func(tc *testcase, _ int) {
					(*tc)()
					Expect(true, "tests evaluate ok in non-parallel tests").To(BeTrue())
				},
				testcases)
		})
		Expect(elapsed).ToNot(BeLessThan(30 * time.Millisecond))
	})

	Run("using Parallel()", func() {
		elapsed := RunWithTimer(func() {
			RunScenarios(
				func(tc *testcase, _ int) {
					Parallel()
					(*tc)()
					Expect(true, "tests evaluate ok in parallel tests").To(BeTrue())
				},
				testcases)
		})
		Expect(elapsed).ToNot(BeGreaterThan(15 * time.Millisecond))
	})

	Run("calling Parallel() from a parallel test", func() {
		result := Test(func() {
			Parallel()
			Parallel()
		})
		result.ExpectInvalid(
			"Parallel() must not be called from a parallel test",
		)
	})

	Run("using RunParallelScenarios()", func() {
		elapsed := RunWithTimer(func() {
			RunParallelScenarios(
				func(tc *testcase, _ int) {
					(*tc)()
					Expect(true, "tests evaluate ok in parallel tests").To(BeTrue())
				},
				testcases)
		})
		Expect(elapsed).To(BeLessThan(12 * time.Millisecond))
	})

	Run("calling RunParallel() from a parallel test", func() {
		result := Test(func() {
			Parallel()
			RunParallel("already parallel", func() {})
		})
		result.ExpectInvalid(
			"RunParallel() must not be called from a parallel test",
		)
	})

	Run("calling RunParallelScenarios() from a parallel test", func() {
		result := Test(func() {
			RunParallel("parallel", func() {
				RunParallelScenarios(func(*int, int) {}, []int{})
			})
		})
		result.ExpectInvalid(
			"RunParallelScenarios() must not be called from a parallel test",
		)
	})
}
