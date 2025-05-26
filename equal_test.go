package test

import "testing"

func TestEqual(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "expected equal and was equal",
			Act: func() { Expect(1).To(Equal(1)) },
		},
		{Scenario: "expected to not be equal and was not equal",
			Act: func() { Expect(1).ToNot(Equal(2)) },
		},
	})
}

func TestDeepEqual(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "expected deep equal and was deep equal",
			Act: func() {
				Expect([]byte{1}).To(DeepEqual([]byte{1}))
			},
		},
		{Scenario: "expected to not be deep equal and was not deep equal",
			Act: func() {
				Expect([]byte{1}).ToNot(DeepEqual([]byte{2}))
			},
		},
	})
}

func ExampleEqual() {
	With(ExampleTestRunner{})

	Expect(1).To(Equal(2))
	Expect("the hobbit").To(Equal("the lord of the rings"))

	// this will not compile because the types are not the same:
	// Expect(42).To(Equal("the answer"))

	// this will not compile because the types are not comparable:
	// Expect([]int{1, 2, 3}).To(Equal([]int{1, 2, 3}))
	//
	// instead, use:
	//   Expect(..).To(DeepEqual(..))
	//
	// or, for slices:
	//   Expect(..).To(EqualSlice(..))

	// Output:
	// expected 2, got 1
	//
	// expected: "the lord of the rings"
	// got     : "the hobbit"
}

func ExampleDeepEqual() {
	With(ExampleTestRunner{})

	Expect([]byte{1, 2, 3}).To(DeepEqual([]byte{1, 2, 4}))
	Expect([]uint8{1, 1, 2, 3, 5}).To(DeepEqual([]uint8{1, 2, 4, 8, 16}))

	// this will not compile because the types are not the same:
	// Expect([]uint8{1, 1, 2, 3, 5}).To(DeepEqual([]int{1,1,2,3,5}))

	// Output:
	// expected [1 2 4], got [1 2 3]
	//
	// expected: [1 2 4 8 16]
	// got     : [1 1 2 3 5]
}
