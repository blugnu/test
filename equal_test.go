package test

import "testing"

func TestEqual(t *testing.T) {
	With(t)

	type testcase struct {
		scenario string
		act      func()
		assert   func(test R)
	}
	RunScenarios(
		func(tc testcase) {
			tc.assert(Test(tc.act))
		},
		[]testcase{
			{scenario: "Equal(int)",
				act: func() { Expect(1).To(Equal(2)) },
				assert: func(test R) {
					Expect(test.Outcome).To(Equal(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						"expected 2, got 1",
					}))
				},
			},
			{scenario: "Equal(string)",
				act: func() { Expect("the quick brown fox").To(Equal("jumped over the lazy dog")) },
				assert: func(test R) {
					Expect(test.Outcome).To(Equal(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						"expected: \"jumped over the lazy dog\"",
						"got     : \"the quick brown fox\"",
					}))
				},
			},
			{scenario: "Equal(struct)",
				act: func() {
					type foo struct {
						name string
					}
					Expect(foo{"ford"}).To(Equal(foo{"arthur"}))
				},
				assert: func(test R) {
					Expect(test.Outcome).To(Equal(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						`expected: test.foo{name:"arthur"}`,
						`got     : test.foo{name:"ford"}`,
					}))
				},
			},
		},
	)
}

func TestDeepEqual(t *testing.T) {
	With(t)

	type testcase struct {
		scenario string
		act      func()
		assert   func(test R)
	}
	RunScenarios(
		func(tc testcase) {
			tc.assert(Test(tc.act))
		},
		[]testcase{
			{scenario: "DeepEqual(int)",
				act: func() { Expect(1).To(DeepEqual(2)) },
				assert: func(test R) {
					Expect(test.Outcome).To(DeepEqual(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						"expected 2, got 1",
					}))
				},
			},
			{scenario: "DeepEqual(string)",
				act: func() { Expect("the quick brown fox").To(DeepEqual("jumped over the lazy dog")) },
				assert: func(test R) {
					Expect(test.Outcome).To(DeepEqual(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						"expected: \"jumped over the lazy dog\"",
						"got     : \"the quick brown fox\"",
					}))
				},
			},
			{scenario: "DeepEqual(struct)",
				act: func() {
					type foo struct {
						name string
					}
					Expect(foo{"ford"}).To(DeepEqual(foo{"arthur"}))
				},
				assert: func(test R) {
					Expect(test.Outcome).To(DeepEqual(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						`expected: test.foo{name:"arthur"}`,
						`got     : test.foo{name:"ford"}`,
					}))
				},
			},
			{scenario: "DeepEqual(struct)/equal",
				act: func() {
					type foo struct {
						bytes []byte
					}
					Expect(foo{[]byte{1}}).To(DeepEqual(foo{[]byte{1}}))
				},
				assert: func(test R) {
					Expect(test.Outcome).To(DeepEqual(TestPassed))
					Expect(test.Report).IsEmpty()
				},
			},
			{scenario: "DeepEqual(struct)/inequal",
				act: func() {
					type foo struct {
						bytes []byte
					}
					Expect(foo{[]byte{65}}).To(DeepEqual(foo{[]byte{97}}))
				},
				assert: func(test R) {
					Expect(test.Outcome).To(DeepEqual(TestFailed))
					Expect(test.Report).To(ContainStrings([]string{
						`expected: test.foo{bytes:[]uint8{0x61}}`,
						`got     : test.foo{bytes:[]uint8{0x41}}`,
					}))
				},
			},
		},
	)
}
