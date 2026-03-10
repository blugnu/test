package test_test

import (
	"testing"
	"time"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestFlaky(t *testing.T) {
	With(t)

	result := TestHelper(func() {
		Run(FlakyTest("fails to pass within max attempts", func() {
			Expect(false).To(BeTrue())
		}))
	})
	result.Expect(
		"Flaky test failed after 3 attempts in ", // time is variable so cannot be precisely checked
		"",
		"attempt 1:",
		"runnable-flaky_test.go",
		"  expected true",
		"",
		"attempt 2:",
		"runnable-flaky_test.go",
		"  expected true",
		"",
		"attempt 3:",
		"runnable-flaky_test.go",
		"  expected true",
	)

	result = TestHelper(func() {
		Run(FlakyTest("fails to pass within default time limit", func() {
			time.Sleep(500 * time.Millisecond)
			Expect(false).To(BeTrue())
		}))
	})
	result.Expect(
		"Flaky test failed after 2 attempts in 1.0",
		"",
		"attempt 1:",
		"runnable-flaky_test.go",
		"  expected true",
		"",
		"attempt 2:",
		"runnable-flaky_test.go",
		"  expected true",
	)

	result = TestHelper(func() {
		first := true
		Run(FlakyTest("test passes after a failed attempt", func() {
			if first {
				Expect(false).To(BeTrue())
				first = false
			}
		}))
	})
	result.Expect(test.Passed)

	result = TestHelper(func() {
		Run(FlakyTest("configured max attempts", func() {
			Expect(false).To(BeTrue())
		}, MaxAttempts(1)))
	})
	result.Expect(
		"Flaky test failed after 1 attempt",
		"",
		"attempt 1:",
		"runnable-flaky_test.go",
		"  expected true",
	)

	result = TestHelper(func() {
		Run(FlakyTest("configured max duration", func() {
			time.Sleep(200 * time.Millisecond)
			Expect(false).To(BeTrue())
		}, MaxDuration(100*time.Millisecond)))
	})
	result.Expect(
		"Flaky test failed after 1 attempt in 2", // precise time is not guaranteed but should be at least 2(00ms)
		"",
		"attempt 1:",
		"runnable-flaky_test.go",
		"  expected true",
	)

	result = TestHelper(func() {
		Run(FlakyTest("no wait between attempts", func() {
			Expect(false).To(BeTrue())
		}, MaxAttempts(1), WaitBetweenAttempts(0)))
	})
	result.Expect(
		"Flaky test failed after 1 attempt",
		"",
		"attempt 1:",
		"runnable-flaky_test.go",
		"  expected true",
	)
}
