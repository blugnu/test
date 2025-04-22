package test

import "testing"

func TestEqualBytes(t *testing.T) {
	With(t)

	makeSlice := func(size int, ffat ...int) []byte {
		ff := -1
		if len(ffat) > 0 {
			ff = ffat[0]
		}
		s := make([]byte, size)
		for i := 0; i < size; i++ {
			s[i] = byte(i + 1)
			if i == ff {
				s[i] = 0xff
			}
		}
		return s
	}

	type testcase struct {
		scenario string
		act      func()
		assert   func(R)
	}
	RunScenarios(
		func(tc testcase) {
			tc.assert(Test(tc.act))
		},

		[]testcase{
			{scenario: "equal",
				act: func() { Expect(makeSlice(4)).To(EqualBytes(makeSlice(4))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{scenario: "expect empty",
				act: func() { Expect(makeSlice(3)).To(EqualBytes([]byte{})) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 0, got 3",
						"expected: <empty>",
						"        | ++ ++ ++",
						"got     : 01 02 03",
					}))
				},
			},
			{scenario: "got empty",
				act: func() { Expect(makeSlice(0)).To(EqualBytes(makeSlice(3))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 3, got 0",
						"expected: 01 02 03",
						"        | -- -- --",
						"got     : <empty>",
					}))
				},
			},
			{scenario: "got fewer",
				act: func() { Expect(makeSlice(4)).To(EqualBytes(makeSlice(5))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 5, got 4",
						"expected: ...03 04 05",
						"        |          --",
						"got     : ...03 04",
					}))
				},
			},
			{scenario: "got extra",
				act: func() { Expect(makeSlice(5)).To(EqualBytes(makeSlice(3))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 3, got 5",
						"expected: ...02 03",
						"        |          ++ ++",
						"got     : ...02 03 04 05",
					}))
				},
			},
			{scenario: "short/1 diff at initial",
				act: func() { Expect(makeSlice(3, 0)).To(EqualBytes(makeSlice(3))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  differences at: [0]",
						"expected: 01 02 03",
						"        | **",
						"got     : ff 02 03",
					}))
				},
			},
			{scenario: "short/1 diff in middle",
				act: func() { Expect(makeSlice(3, 1)).To(EqualBytes(makeSlice(3))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  differences at: [1]",
						"expected: 01 02 03",
						"        |    **",
						"got     : 01 ff 03",
					}))
				},
			},
			{scenario: "short/1 diff at end",
				act: func() { Expect(makeSlice(3, 2)).To(EqualBytes(makeSlice(3))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  differences at: [2]",
						"expected: 01 02 03",
						"        |       **",
						"got     : 01 02 ff",
					}))
				},
			},
			{scenario: "long/different lengths",
				act: func() { Expect(makeSlice(13)).To(EqualBytes(makeSlice(20))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 20, got 13",
						"expected: ...0c 0d 0e 0f 10...",
						"        |          -- -- --",
						"got     : ...0c 0d",
					}))
				},
			},
			{scenario: "long/different lengths/other diff shown",
				act: func() { Expect(makeSlice(13, 12)).To(EqualBytes(makeSlice(20))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 20, got 13",
						"  differences at: [12, 13]",
						"expected: ...0c 0d 0e 0f 10...",
						"        |       ** -- -- --",
						"got     : ...0c ff",
					}))
				},
			},
			{scenario: "long/different lengths/other diff not shown",
				act: func() { Expect(makeSlice(13, 8)).To(EqualBytes(makeSlice(20))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  different lengths: expected 20, got 13",
						"  differences at: [8, 13]",
						"expected: ...0c 0d 0e 0f 10...",
						"        |          -- -- --",
						"got     : ...0c 0d",
					}))
				},
			},
			{scenario: "long/diff in middle",
				act: func() { Expect(makeSlice(13, 5)).To(EqualBytes(makeSlice(13))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						"bytes not equal:",
						"  differences at: [5]",
						"expected: ...04 05 06 07 08...",
						"        |          **",
						"got     : ...04 05 ff 07 08...",
					}))
				},
			},
		},
	)
}

func ExampleEqualBytes() {
	With(ExampleTestRunner{})

	a := []byte{0x01, 0x02, 0x03}
	b := []byte{0x01, 0x03, 0x02}

	Expect(a).To(EqualBytes(b))

	// Output:
	// bytes not equal:
	//   differences at: [1, 2]
	// expected: 01 03 02
	//         |    ** **
	// got     : 01 02 03
}
