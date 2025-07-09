package test_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestContainsString(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// ContainString tests
		{Scenario: "expecting substring to exist in string",
			Act: func() { Expect("abc").To(ContainString("a")) },
		},
		{Scenario: "expecting an empty substring",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("abc").To(ContainString(""))
			},
		},
		{Scenario: "expecting a substring that does not exist in string",
			Act: func() {
				Expect("abc").To(ContainString("d"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: string containing: "d"`,
					`got     : "abc"`,
				)
			},
		},
		{Scenario: "not expecting a substring to exist in string",
			Act: func() {
				Expect("abc").ToNot(ContainString("d"))
			},
		},
	}...))
}

func TestMatchRegEx(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{

		// Match tests
		{Scenario: "expecting to match a regex and does",
			Act: func() {
				Expect("email: someone@somewhere.com").To(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"))
			},
		},
		{Scenario: "expecting to match a regex and does not",
			Act: func() {
				Expect("email: <not present>").To(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: string containing match for: "[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"`,
					`got     : "email: <not present>"`,
				)
			},
		},
		{Scenario: "expecting to match an empty regex",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").To(MatchRegEx(""))
			},
		},
		{Scenario: "expecting to match an invalid regex",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").To(MatchRegEx("["))
			},
		},

		{Scenario: "expecting to not match a regex and does",
			Act: func() {
				Expect("email: someone@somewhere.com").ToNot(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: string with no match for: "[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"`,
					`got     : "email: someone@somewhere.com"`,
					`matched : "someone@somewhere.com"`,
				)
			},
		},
		{Scenario: "expecting to not match a regex and does not",
			Act: func() { Expect("email: <not present>").ToNot(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}")) },
		},
	}...))
}
