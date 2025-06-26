package strings_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestStringMatcher(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		// ContainString tests
		{Scenario: "ContainString",
			Act: func() { Expect("abc").To(ContainString("a")) },
		},
		{Scenario: "ContainString, case-insensitive",
			Act: func() { Expect("abc").To(ContainString("AB"), opt.CaseSensitive(false)) },
		},
		{Scenario: "ContainString/empty string",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("abc").To(ContainString(""))
			},
		},
		{Scenario: "ContainString/does not contain",
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
		{Scenario: "ToNot(ContainString)/does not contain",
			Act: func() {
				Expect("abc").ToNot(ContainString("d"))
			},
		},
		{Scenario: "ToNot(ContainString)/contains string",
			Act: func() {
				Expect("abc").ToNot(ContainString("c"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: string not containing: "c"`,
					`got     : "abc"`,
					`             ^`,
				)
			},
		},
		{Scenario: "ToNot(ContainString)/contains string (unquoted)",
			Act: func() {
				Expect("abcdef").ToNot(ContainString("cd"), opt.UnquotedStrings())
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: string not containing: c`,
					`got     : abcdef`,
					`            ^^`,
				)
			},
		},

		// Match tests
		{Scenario: "To(Match)/matches",
			Act: func() {
				Expect("email: someone@somewhere.com").To(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}"))
			},
		},
		{Scenario: "To(Match)/does not match",
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
		{Scenario: "To(Match)/empty expression",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").To(MatchRegEx(""))
			},
		},
		{Scenario: "To(Match)/invalid regex",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").To(MatchRegEx("["))
			},
		},

		{Scenario: "ToNot(Match)/matches",
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
		{Scenario: "ToNot(Match)/does not match",
			Act: func() { Expect("email: <not present>").ToNot(MatchRegEx("[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+[a-z]{2,3}")) },
		},
		{Scenario: "ToNot(Match)/empty expression",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").ToNot(MatchRegEx(""))
			},
		},
		{Scenario: "ToNot(Match)/invalid regex",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				Expect("email: someone@somewhere.com").ToNot(MatchRegEx("["))
			},
		},

		{Scenario: "matching with a custom type based on string",
			Act: func() {
				type myString string
				got := myString("abc")
				Expect(string(got)).To(MatchRegEx("a"))
			},
		},
	})
}
