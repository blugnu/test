package expect_test

import (
	"errors"
	"fmt"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/test"
)

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

func ExampleError() {
	test.Example()

	var (
		errWrapped error = errors.New("wrapped error")
		errOther   error = errors.New("some other error")
		err        error = fmt.Errorf("an error: %w", errWrapped)
	)

	expect.Error(err)             // will pass
	expect.Error(err, errWrapped) // will pass
	expect.Error(nil)             // will fail
	expect.Error(err, errOther)   // will fail because err does not match or wrap errOther

	// Output:
	// expected error
	//
	// expected error: "some other error"        [*errors.errorString]
	// got           : "an error: wrapped error" [*fmt.wrapError]
}

func ExampleErrorAs() {
	test.Example()

	var (
		errString error = errors.New("a standard error")
		errCustom error = CustomError("a custom error")
	)

	_, _ = expect.ErrorAs[CustomError](errCustom) // this test will pass
	_, _ = expect.ErrorAs[CustomError](errString) // this test will fail

	// Output:
	// expected error to match or wrap: expect_test.CustomError
	// got: "a standard error" [*errors.errorString]
}

func ExampleNoError() {
	test.Example()

	var err error = errors.New("an error")

	expect.NoError(err) // will fail
	expect.NoError(nil) // will pass

	// Output:
	// unexpected error
	// got: "an error" [*errors.errorString]
}
