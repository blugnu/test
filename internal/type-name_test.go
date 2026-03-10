package internal_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal"
)

func TestTypeName(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "basic type",
			Act: func() {
				Expect(internal.TypeName[int]()).To(Equal("int"))
			},
		},

		{Scenario: "pointer to basic type",
			Act: func() {
				Expect(internal.TypeName[*int]()).To(Equal("*int"))
			},
		},

		{Scenario: "struct type",
			Act: func() {
				type MyStruct struct{}
				Expect(internal.TypeName[MyStruct]()).To(Equal("internal_test.MyStruct"))
			},
		},

		{Scenario: "pointer to struct type",
			Act: func() {
				type MyStruct struct{}
				Expect(internal.TypeName[*MyStruct]()).To(Equal("*internal_test.MyStruct"))
			},
		},

		{Scenario: "type based on basic type",
			Act: func() {
				type MyInt int
				Expect(internal.TypeName[MyInt]()).To(Equal("internal_test.MyInt"))
			},
		},

		{Scenario: "pointer to type based on basic type",
			Act: func() {
				type MyInt int
				Expect(internal.TypeName[*MyInt]()).To(Equal("*internal_test.MyInt"))
			},
		},

		{Scenario: "interface type",
			Act: func() {
				Expect(internal.TypeName[error]()).To(Equal("error"))
			},
		},

		{Scenario: "nil interface",
			Act: func() {
				var err error = nil
				Expect(internal.TypeName(err)).To(Equal("error"))
			},
		},

		{Scenario: "anonymous struct",
			Act: func() {
				Expect(internal.TypeName[struct{ foo int }]()).To(Equal("struct { foo int }"))
			},
		},

		{Scenario: "anonymous interface",
			Act: func() {
				Expect(internal.TypeName[interface{ Foo() }]()).To(Equal("<anon>"))
			},
		},
	}...))
}
