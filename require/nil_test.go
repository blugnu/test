package require_test

import (
	"errors"

	"github.com/blugnu/test/require"
	"github.com/blugnu/test/test"
)

func ExampleNil() {
	test.Example()

	err := errors.New("an error")
	require.Nil(nil) // will pass
	require.Nil(err) // will fail
	require.Nil(42)  // will not be reached

	// Output:
	// expected nil, got error: an error
}

func ExampleNotNil() {
	test.Example()

	err := errors.New("an error")
	require.NotNil(err) // will pass
	require.NotNil(nil) // will fail
	require.NotNil(42)  // will not be reached

	// Output:
	// expected not nil
}
