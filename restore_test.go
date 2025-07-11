package test_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/blugnu/test"
)

func TestOriginal(t *testing.T) {
	With(t)

	result := TestHelper(func() {
		Original((*int)(nil))
	})

	result.ExpectInvalid("test.Original: cannot create an override for a nil pointer")
}

func TestRestore(t *testing.T) {
	With(t)

	i := 42
	Run(Test("replaced in subtest", func() {
		defer Restore(Original(&i).ReplacedBy(12))

		Expect(i).To(Equal(12))
	}))

	Expect(i).To(Equal(42))
}

// ExampleRestore demonstrates how to use Restore to temporarily
// replace a variable with a new value for the duration of a test
//
// The example simulates a naive approach to mocking the time.Now function
// and demonstrates that Restore can be used to temporarily replace
// variables in a test, including variables that are function references.
//
// To ensure predictable output for the example, instead of using
// the real time.Now function, we define a variable `now` that returns
// a fixed time value.  This allows us to control the output of the example
// without relying on the current time, which would vary depending on when
// the example is run.
//
// This is not a robust mechanism for mocking time.Now, but it serves
// as an illustration.  For a more robust approach to mocking time,
// consider using a package like github.com/blugnu/time, or similar.
func ExampleRestore() {
	// establish "now" as a function that returns the release date of Go 1.0
	// as the "current time".
	var now = func() time.Time {
		return time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)
	}

	// simulate a sub-test where the "current time" is temporarily
	// replaced with a fixed time value.
	func() {
		var fakeTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		defer Restore(Original(&now).ReplacedBy(func() time.Time { return fakeTime }))

		fmt.Println("now() returns: ", now())
	}()

	fmt.Println("now() returns: ", now())

	// Output:
	// now() returns:  2020-01-01 00:00:00 +0000 UTC
	// now() returns:  2012-03-28 00:00:00 +0000 UTC
}
