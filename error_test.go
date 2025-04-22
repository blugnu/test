package test

import (
	"errors"
	"testing"
)

func TestExpectError(t *testing.T) {
	With(t)

	type testcase struct {
		scenario string
		act      func()
		assert   func(R)
	}
	RunScenarios(
		func(tc testcase) {
			result := Test(tc.act)
			tc.assert(result)
		},
		[]testcase{
			{
				scenario: "any (nil) is nil",
				act:      func() { var a any; ExpectError(a).Is(nil) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{
				scenario: "any (non-nil) is nil",
				act:      func() { var a any = 1; ExpectError(a).Is(nil) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"expected error, got int",
					}))
				},
			},
			{
				scenario: "error nil (error)",
				act:      func() { var err error; ExpectError(err).Is(nil) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{
				scenario: "matching sentinel",
				act:      func() { sent := errors.New("sentinel"); ExpectError(sent).Is(sent) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{
				scenario: "different sentinel",
				act: func() {
					senta := errors.New("sentinel-a")
					sentb := errors.New("sentinel-b")
					ExpectError(senta).Is(sentb)
				},
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"expected error: sentinel-b",
						"got           : sentinel-a",
					}))
				},
			},
		},
	)
}
