package test

import "testing"

func TestExpectType(t *testing.T) {
	With(t)

	type testcase struct {
		scenario string
		act      func() (int, bool)
		assert   func(test R, result int, ok bool)
	}
	RunScenarios(
		func(tc testcase) {
			var result int
			var ok bool
			test := Test(func() {
				result, ok = tc.act()
			})

			tc.assert(test, result, ok)
		},
		[]testcase{
			{scenario: "[int](int)",
				act: func() (int, bool) {
					return ExpectType[int](1)
				},
				assert: func(test R, result int, ok bool) {
					Expect(test.Outcome).To(Equal(TestPassed))
					ExpectEmpty(test.Report)

					Expect(result).To(Equal(1))
					Expect(ok).To(BeTrue())
				},
			},
			{scenario: "[int](string)",
				act: func() (int, bool) {
					return ExpectType[int]("string")
				},
				assert: func(test R, result int, ok bool) {
					Expect(test.Outcome).To(Equal(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						"expected type int, got string",
					}))

					Expect(result).To(Equal(0))
					Expect(ok).To(BeFalse())
				},
			},
		},
	)
}
