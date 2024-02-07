package test

import "testing"

func TestFormat(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		v        any
		Format
		result string
	}{
		{scenario: "FormatNotSet", v: 42, result: "42"},
		{scenario: "FormatDefault", v: 42, Format: FormatDefault, result: "42"},
		{scenario: "FormatDecl", v: 42, Format: FormatDecl, result: "42"},
		{scenario: "FormatBin", v: 42, Format: FormatBin, result: "00101010"},
		{scenario: "FormatHex", v: 42, Format: FormatHex, result: "2a"},
		{scenario: "custom", v: 42, Format: Format("foo: %d"), result: "foo: 42"},
		{scenario: "struct {int}{42}, FormatDefault", v: struct{ int }{42}, Format: FormatDefault, result: "{42}"},
		{scenario: "struct {int}{42}, FormatDecl", v: struct{ int }{42}, Format: FormatDecl, result: "struct { int }{int:42}"},
		{scenario: "struct {int}{42}, FormatBin", v: struct{ int }{42}, Format: FormatHex, result: "{2a}"},
		{scenario: "struct {int}{42}, FormatHex", v: struct{ int }{42}, Format: FormatHex, result: "{2a}"},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			got := format(tc.v, tc.Format)

			// ASSERT
			Equal(t, tc.result, got)
		})
	}
}
