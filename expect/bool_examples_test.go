package expect_test

import (
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/test"
)

func ExampleFalse() {
	test.Example()

	var a = 1

	expect.False(a == 1)     // will fail
	expect.False(a-1+1 == 1) // will also fail
	expect.False(a == 2)     // will pass

	// Output:
	// expected false
	//
	// expected false
}

func ExampleTrue() {
	test.Example()

	var a = 1

	expect.True(a == 1) // will pass
	expect.True(a == 2) // will fail
	expect.True(a == 3) // will also fail

	// Output:
	// expected true
	//
	// expected true
}
