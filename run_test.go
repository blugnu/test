package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestRun(t *testing.T) {
	With(t)

	Run(Test("called from an Example", func() {
		defer Expect(Panic()).DidOccur()

		test.Example()

		Run(Test("this is invalid", func() {}))
	}))
}

func TestRun_Test(t *testing.T) {
	With(t)

	result := TestHelper(func() {
		Run(Test("named_test", func() {
			Expect(false).To(BeTrue())
		}))
	})

	Expect(result.FailedTests).To(ContainItem("TestRun_Test/named_test"))
	result.Expect("expected true, got false")
}

func TestRun_Testcases(t *testing.T) {
	With(t)

	type tc struct {
		a, b   int
		result int
	}

	Run(Test("no test cases", func() {
		result := TestHelper(func() {
			Run(Testcases(
				ForEach(func(tc tc) { /* NO-OP */ }),
			))
		})

		result.Expect("no test cases provided")
	}))

	Run(Test("For(nil)", func() {
		result := TestHelper(func() {
			Run(Testcases(For[tc](nil)))
		})

		result.ExpectInvalid("For() function cannot be nil")
	}))

	Run(Test("ForEach(nil)", func() {
		result := TestHelper(func() {
			Run(Testcases(ForEach[tc](nil)))
		})

		result.ExpectInvalid("ForEach() function cannot be nil")
	}))

	Run(Test("anonymous", func() {
		result := TestHelper(func() {
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Cases([]tc{
					{a: 2, b: 2, result: 4},
					{a: 2, b: 2, result: 5},
				}),
			))
		})

		Expect(result.FailedTests).To(ContainItem("TestRun_Testcases/anonymous/testcase-002"))
		result.Expect("expected 5, got 4")
	}))

	Run(Test("debug anonymous", func() {
		result := TestHelper(func() {
			type tc struct {
				a, b   int
				result int
				debug  bool
			}

			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Cases([]tc{
					{a: 2, b: 2, result: 4},
					{a: 2, b: 2, result: 5, debug: true},
					{a: 2, b: 2, result: 6},
				}),
			))
		})

		Expect(result.FailedTests).To(ContainItem("TestRun_Testcases/debug_anonymous/testcase-002"))
		Expect(result.FailedTests).ToNot(ContainItem("TestRun_Testcases/debug_anonymous/testcase-003"))
		result.Expect("expected 5, got 4")
	}))

	Run(Test("cases", func() {
		testcaseNames := []string{}
		result := TestHelper(func() {
			Run(Testcases(
				For(func(name string, tc tc) {
					testcaseNames = append(testcaseNames, name)
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Case("2+2==4", tc{a: 2, b: 2, result: 4}),
				Case("2+2==5", tc{a: 2, b: 2, result: 5}),
			))
		})

		Expect(testcaseNames).To(EqualSlice([]string{
			"2+2==4",
			"2+2==5",
		}), opt.AnyOrder())

		Expect(result.FailedTests).To(ContainItem("TestRun_Testcases/cases/2+2==5"))
		result.Expect("expected 5, got 4")
	}))

	Run(Test("debug method", func() {
		result := TestHelper(func() {
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Case("2+2==4", tc{a: 2, b: 2, result: 4}),
				Debug("2+2==5", tc{a: 2, b: 2, result: 5}),
				Case("2+2==6", tc{a: 2, b: 2, result: 6}),
			))
		})

		Expect(result.FailedTests).To(ContainItem("TestRun_Testcases/debug_method/2+2==5"))
		Expect(result.FailedTests).ToNot(ContainItem("TestRun_Testcases/debug_method/2+2==6"))

		result.Expect("expected 5, got 4")
	}))

	Run(Test("debug field", func() {
		result := TestHelper(func() {
			type tc struct {
				a, b   int
				result int
				debug  bool
			}

			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Case("2+2==5", tc{a: 2, b: 2, result: 5, debug: true}),
				Case("2+2==6", tc{a: 2, b: 2, result: 6}),
			))
		})

		Expect(result.FailedTests).To(ContainItem("TestRun_Testcases/debug_field/2+2==5"))
		Expect(result.FailedTests).ToNot(ContainItem("TestRun_Testcases/debug_field/2+2==6"))

		result.Expect("expected 5, got 4")
	}))

	Run(Test("debug method overrides skip", func() {
		result := TestHelper(func() {
			type tc struct {
				a, b   int
				result int
				skip   bool
			}
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Debug("2+2==5", tc{a: 2, b: 2, result: 5, skip: true}),
			))
		})

		result.Expect(TestFailed, "expected 5, got 4")
	}))

	Run(Test("skip method", func() {
		result := TestHelper(func() {
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Case("2+2==4", tc{a: 2, b: 2, result: 4}),
				Skip("2+2==5", tc{a: 2, b: 2, result: 5}),
			))
		})

		result.ExpectWarning("1 of 2 cases were skipped")
	}))

	Run(Test("skip field", func() {
		result := TestHelper(func() {
			type tc struct {
				a, b   int
				result int
				skip   bool
			}
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Skip("2+2==5", tc{a: 2, b: 2, result: 5, skip: true}),
			))
		})

		result.ExpectWarning("all cases were skipped")
	}))

	Run(Test("skip method overrides debug field", func() {
		result := TestHelper(func() {
			type tc struct {
				a, b   int
				result int
				debug  bool
			}
			Run(Testcases(
				ForEach(func(tc tc) {
					Expect(tc.a + tc.b).To(Equal(tc.result))
				}),
				Skip("2+2==5", tc{a: 2, b: 2, result: 5, debug: true}),
			))
		})

		result.ExpectWarning("all cases were skipped")
	}))
}
