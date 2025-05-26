package bytes_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestEqualMatcher(t *testing.T) {
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

	RunTestScenarios([]TestScenario{
		{Scenario: "expected equal and was equal",
			Act: func() { Expect(makeSlice(4)).To(EqualBytes(makeSlice(4))) },
		},
		{Scenario: "expected to not be equal and was equal",
			Act: func() { Expect(makeSlice(4)).ToNot(EqualBytes(makeSlice(4))) },
			Assert: func(result *R) {
				result.Expect(
					"unexpected: []byte should not be equal",
				)
			},
		},
		{Scenario: "expected empty and was not empty",
			Act: func() { Expect(makeSlice(3)).To(EqualBytes([]byte{})) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 0, got 3",
					"expected: <empty>",
					"        | ++ ++ ++",
					"got     : 01 02 03",
				)
			},
		},
		{Scenario: "got empty",
			Act: func() { Expect(makeSlice(0)).To(EqualBytes(makeSlice(3))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 3, got 0",
					"expected: 01 02 03",
					"        | -- -- --",
					"got     : <empty>",
				)
			},
		},
		{Scenario: "got fewer",
			Act: func() { Expect(makeSlice(4)).To(EqualBytes(makeSlice(5))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 5, got 4",
					"expected: ...03 04 05",
					"        |          --",
					"got     : ...03 04",
				)
			},
		},
		{Scenario: "got extra",
			Act: func() { Expect(makeSlice(5)).To(EqualBytes(makeSlice(3))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 3, got 5",
					"expected: ...02 03",
					"        |          ++ ++",
					"got     : ...02 03 04 05",
				)
			},
		},
		{Scenario: "short/1 diff at initial",
			Act: func() { Expect(makeSlice(3, 0)).To(EqualBytes(makeSlice(3))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  differences at: [0]",
					"expected: 01 02 03",
					"        | **",
					"got     : ff 02 03",
				)
			},
		},
		{Scenario: "short/1 diff in middle",
			Act: func() { Expect(makeSlice(3, 1)).To(EqualBytes(makeSlice(3))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  differences at: [1]",
					"expected: 01 02 03",
					"        |    **",
					"got     : 01 ff 03",
				)
			},
		},
		{Scenario: "short/1 diff at end",
			Act: func() { Expect(makeSlice(3, 2)).To(EqualBytes(makeSlice(3))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  differences at: [2]",
					"expected: 01 02 03",
					"        |       **",
					"got     : 01 02 ff",
				)
			},
		},
		{Scenario: "long/different lengths",
			Act: func() { Expect(makeSlice(13)).To(EqualBytes(makeSlice(20))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 20, got 13",
					"expected: ...0c 0d 0e 0f 10...",
					"        |          -- -- --",
					"got     : ...0c 0d",
				)
			},
		},
		{Scenario: "long/different lengths/other diff shown",
			Act: func() { Expect(makeSlice(13, 12)).To(EqualBytes(makeSlice(20))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 20, got 13",
					"  differences at: [12, 13]",
					"expected: ...0c 0d 0e 0f 10...",
					"        |       ** -- -- --",
					"got     : ...0c ff",
				)
			},
		},
		{Scenario: "long/different lengths/other diff not shown",
			Act: func() { Expect(makeSlice(13, 8)).To(EqualBytes(makeSlice(20))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  different lengths: expected 20, got 13",
					"  differences at: [8, 13]",
					"expected: ...0c 0d 0e 0f 10...",
					"        |          -- -- --",
					"got     : ...0c 0d",
				)
			},
		},
		{Scenario: "long/diff in middle",
			Act: func() { Expect(makeSlice(13, 5)).To(EqualBytes(makeSlice(13))) },
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  differences at: [5]",
					"expected: ...04 05 06 07 08...",
					"        |          **",
					"got     : ...04 05 ff 07 08...",
				)
			},
		},
		{Scenario: "test passes with custom byte type",
			Act: func() {
				type MyByte byte
				Expect([]MyByte{1, 2, 3}).To(EqualBytes([]MyByte{1, 2, 3}))
			},
		},
		{Scenario: "test fails with custom byte type",
			Act: func() {
				type MyByte byte
				Expect([]MyByte{1, 2, 3}).To(EqualBytes([]MyByte{1, 2, 4}))
			},
			Assert: func(result *R) {
				result.Expect(
					"bytes not equal:",
					"  differences at: [2]",
					"expected: 01 02 04",
					"        |       **",
					"got     : 01 02 03",
				)
			},
		},
	})
}
