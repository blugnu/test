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

	result = opt.ValueAsString("test 1")
	if result != "\"test 1\"" {
		t.Errorf("Expected '\"test 1\"', got '%s'", result)
	}

	result = opt.ValueAsString("test 2", opt.QuotedStrings(true))
	if result != "\"test 2\"" {
		t.Errorf("Expected '\"test 2\"', got '%s'", result)
	}

	result = opt.ValueAsString("test 3", opt.QuotedStrings(false))
	if result != "test 3" {
		t.Errorf("Expected 'test 3', got '%s'", result)
	}

	result = opt.ValueAsString("test 4", opt.UnquotedStrings())
	if result != "test 4" {
		t.Errorf("Expected 'test 4', got '%s'", result)
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
