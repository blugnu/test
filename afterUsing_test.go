package test

import (
	"fmt"
	"testing"
	"time"
)

func TestAfterUsing(t *testing.T) {
	With(t)

	// ARRANGE
	var now = time.Now
	var fakeTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	restoreFn := AfterUsing(&now, func() time.Time { return fakeTime })

	// ACT
	result := now()

	// ASSERT
	Expect(result).To(Equal(fakeTime))

	// ACT
	restoreFn()
	result = now()

	// ASSERT
	Expect(result).ToNot(Equal(fakeTime))
}

func ExampleAfterUsing() {
	// ARRANGE
	var now = time.Now
	var fakeTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	defer AfterUsing(&now, func() time.Time { return fakeTime })()

	fmt.Println("Current time: ", now())

	// Output:
	// Current time:  2020-01-01 00:00:00 +0000 UTC
}
