package report_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/report"
)

func TestTypeName(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "basic type",
			Act: func() {
				Expect(report.TypeName[int]()).To(Equal("int"))
			},
		},

		{Scenario: "pointer to basic type",
			Act: func() {
				Expect(report.TypeName[*int]()).To(Equal("*int"))
			},
		},

		{Scenario: "struct type",
			Act: func() {
				type MyStruct struct{}
				Expect(report.TypeName[MyStruct]()).To(Equal("report_test.MyStruct"))
			},
		},

		{Scenario: "pointer to struct type",
			Act: func() {
				type MyStruct struct{}
				Expect(report.TypeName[*MyStruct]()).To(Equal("*report_test.MyStruct"))
			},
		},

		{Scenario: "type based on basic type",
			Act: func() {
				type MyInt int
				Expect(report.TypeName[MyInt]()).To(Equal("report_test.MyInt"))
			},
		},

		{Scenario: "pointer to type based on basic type",
			Act: func() {
				type MyInt int
				Expect(report.TypeName[*MyInt]()).To(Equal("*report_test.MyInt"))
			},
		},

		{Scenario: "interface type",
			Act: func() {
				Expect(report.TypeName[error]()).To(Equal("error"))
			},
		},

		{Scenario: "nil interface",
			Act: func() {
				var err error = nil
				Expect(report.TypeName(err)).To(Equal("error"))
			},
		},

		{Scenario: "anonymous struct",
			Act: func() {
				Expect(report.TypeName[struct{ foo, bar int }]()).To(Equal("struct{foo int; bar int}"))
			},
		},

		{Scenario: "anonymous interface",
			Act: func() {
				Expect(report.TypeName[interface{ Foo() }]()).To(Equal("<anonymous interface>"))
			},
		},
	}...))
}
