package opt_test

import (
	"fmt"
	"testing"

	"github.com/blugnu/test/opt"
)

func TestValueAsString(t *testing.T) {
	var result string

	result = opt.ValueAsString(nil)
	if result != "nil" {
		t.Errorf("Expected 'nil', got '%s'", result)
	}

	result = opt.ValueAsString("test")
	if result != "\"test\"" {
		t.Errorf("Expected '\"test\"', got '%s'", result)
	}

	result = opt.ValueAsString("test", opt.QuotedStrings(true))
	if result != "\"test\"" {
		t.Errorf("Expected '\"test\"', got '%s'", result)
	}

	result = opt.ValueAsString("test", opt.QuotedStrings(false))
	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}

	result = opt.ValueAsString("test", opt.UnquotedStrings())
	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}

	result = opt.ValueAsString(42, opt.QuotedStrings(false))
	if result != "42" {
		t.Errorf("Expected '42', got '%s'", result)
	}

	result = opt.ValueAsString([]int{1}, opt.AsDeclaration(true))
	if result != "[]int{1}" {
		t.Errorf("Expected '[]int{1}', got '%s'", result)
	}
}

func ExampleValueAsString() {
	// non-string values are returned as unquoted strings
	fmt.Println(opt.ValueAsString(42))

	// string values are returned as quoted strings
	fmt.Println(opt.ValueAsString("example"))

	// to suppress the quotes, use opt.UnquotedString(true)
	fmt.Println(opt.ValueAsString("example", opt.UnquotedStrings()))

	// Output:
	// 42
	// "example"
	// example
}
