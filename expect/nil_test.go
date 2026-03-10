package expect_test

import (
	"errors"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/test"
)

func ExampleNil() {
	test.Example()

	err := errors.New("an error")
	expect.Nil(nil)
	expect.Nil(err) // this test will fail

	// this test will fail as invalid (int is not nilable)
	expect.Nil(42)

	// Output:
	// expected nil, got error: an error
	// <== INVALID TEST
	// nilness.Matcher: values of type 'int' are not nilable
}

func ExampleNotNil() {
	test.Example()

	err := errors.New("an error")
	expect.NotNil(nil) // this test will fail
	expect.NotNil(err) // this test is valid and will pass (error is nilable)

	// this test is valid and will also pass.  Non-nilable types are
	// inherently "not nil"; the test is unnecessary, but not invalid
	expect.NotNil(42)

	// Output:
	// expected not nil
}
