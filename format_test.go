package test

import "testing"

func TestFormat(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		name string
		v    any
		Format
		result string
	}{
		{name: "FormatNotSet", v: 42, result: "42"},
		{name: "FormatDefault", v: 42, Format: FormatDefault, result: "42"},
		{name: "FormatDecl", v: 42, Format: FormatDecl, result: "42"},
		{name: "FormatBin", v: 42, Format: FormatBin, result: "00101010"},
		{name: "FormatHex", v: 42, Format: FormatHex, result: "2a"},
		{name: "custom", v: 42, Format: Format("foo: %d"), result: "foo: 42"},
		{name: "struct {int}{42}, FormatDefault", v: struct{ int }{42}, Format: FormatDefault, result: "{42}"},
		{name: "struct {int}{42}, FormatDecl", v: struct{ int }{42}, Format: FormatDecl, result: "struct { int }{int:42}"},
		{name: "struct {int}{42}, FormatBin", v: struct{ int }{42}, Format: FormatHex, result: "{2a}"},
		{name: "struct {int}{42}, FormatHex", v: struct{ int }{42}, Format: FormatHex, result: "{2a}"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			got := format(tc.v, tc.Format)

			// ASSERT
			Equal(t, tc.result, got)
		})
	}
}
