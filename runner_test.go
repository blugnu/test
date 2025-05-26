package test

import (
	"fmt"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	With(t)

	type testCase struct {
		int
	}

	Run("no setup or tear down", func() {
		r := Runner[testCase]{}

		Run("RunScenarios", func() {
			r.RunScenarios(
				func(tc *testCase, i int) {},
				[]testCase{{}},
			)
		})

		Run("RunParallel", func() {
			r.RunParallel(
				func(tc *testCase, i int) {},
				[]testCase{{}},
			)
		})
	})

	Run("with setup and tear down", func() {
		r := Runner[testCase]{
			BeforeEach: func(tc *testCase, i int) {
				tc.int = (i + 1) * 2
				Expect(tc.int, fmt.Sprintf("before each: %d", i)).ToNot(BeLessThan(0))
			},
			AfterEach: func(tc *testCase, i int) {
				Expect(tc.int, fmt.Sprintf("after each: %d", i)).To(Equal(0))
			},
		}

		Run("RunScenarios", func() {
			elapsed := RunWithTimer(func() {
				r.RunScenarios(
					func(tc *testCase, i int) {
						time.Sleep(10 * time.Millisecond)
						Expect(tc.int, fmt.Sprintf("test case %d", i)).To(Equal((i + 1) * 2))

						tc.int = 0
					},
					[]testCase{{}, {}, {}})
			})
			Expect(elapsed, "elapsed time").To(BeGreaterThan(30 * time.Millisecond))
		})

		Run("RunParallel", func() {
			elapsed := RunWithTimer(func() {
				r.RunParallel(
					func(tc *testCase, i int) {
						time.Sleep(10 * time.Millisecond)
						Expect(tc.int, fmt.Sprintf("int: tc %d", i)).ToNot(Equal(i))

						tc.int = 0
					},
					[]testCase{{1}, {2}, {3}})
			})
			Expect(elapsed, "elapsed time").To(BeLessThan(30 * time.Millisecond))
		})
	})
}

func TestRunner_RunParallel(t *testing.T) {
	With(t)

	type testCase struct {
		int
	}
	r := Runner[testCase]{
		BeforeEach: func(tc *testCase, i int) {
			tc.int = (i + 1) * 2
			Expect(tc.int, fmt.Sprintf("before each: %d", i)).ToNot(BeLessThan(0))
		},
		AfterEach: func(tc *testCase, i int) {
			Expect(tc.int, fmt.Sprintf("after each: %d", i)).To(Equal(0))
		},
	}

	Run("RunParallel", func() {
		elapsed := RunWithTimer(func() {
			r.RunParallel(
				func(tc *testCase, i int) {
					time.Sleep(10 * time.Millisecond)
					Expect(tc.int, fmt.Sprintf("test case %d", i)).To(Equal((i + 1) * 2))

					tc.int = 0
				},
				[]testCase{{}, {}, {}})
		})
		Expect(elapsed, "elapsed time").To(BeLessThan(15 * time.Millisecond))
	})

	Run("called from a parallel test", func() {
		result := Test(func() {
			RunParallel("parallel", func() {
				r := Runner[int]{}
				r.RunParallel(func(_ *int, _ int) {}, []int{})
			})
		})
		result.ExpectInvalid(
			"Runner.RunParallel() must not be called from a parallel test",
		)
	})
}
