package test

import (
	"fmt"
	"strings"
	"time"
)

// flakyRunner runs a test function that is prone to flakiness. It implements
// the [Runnable] interface, allowing it to be used with the [Run] function.
type flakyRunner struct {
	name        string
	fn          func()
	maxAttempts uint
	maxDuration time.Duration
	wait        time.Duration
}

// Run runs the named test function as a subtest in a separate test runner using
// [TestHelper]. This allows for retries if the test fails, without failing the
// test overall, allowing the test to pass if the flaky test ultimately passes.
//
// Failure reports are collected for each failed attempt.
//
// If the test passes before exceeding the allowed attempts, the failure reports
// are discarded, otherwise the test is failed with a report that details all of
// the failed attempts.
func (tr flakyRunner) Run() {
	T().Helper()

	Run(Test(tr.name, func() {
		T().Helper()

		outcome, attempts, reports, elapsed := tr.runTestWithRetries()
		if outcome == TestPassed {
			return
		}

		// FUTURE: consider condensing the reports to identify ranges of attempts
		// that failed with the same error consecutively, rather than reporting
		// each attempt separately.

		summary := "Flaky test failed after 1 attempt"
		if attempts != 1 {
			summary = fmt.Sprintf("Flaky test failed after %d attempts", attempts)
		}

		report := []string{summary + " in " + elapsed.String()}
		for i, failure := range reports {
			report = append(report, "", fmt.Sprintf("attempt %d:", i+1))
			report = append(report, failure...)
		}

		Error(strings.Join(report, "\n"))
	}))
}

// runTestWithRetries runs the test function with retries based on the configured
// maximum number of attempts and maximum duration. It returns the outcome of the
// test, the number of attempts made, the reports of each failed attempt, and the
// total elapsed time.
//
// If the test passes, the reports are discarded and the outcome is [TestPassed].
// If the test fails after all attempts, the reports contain the details of each
// failed attempt, and the outcome is [TestFailed].
//
// The elapsed time is truncated to milliseconds for consistency in reporting.
func (tr flakyRunner) runTestWithRetries() (TestOutcome, uint, [][]string, time.Duration) {
	T().Helper()

	var (
		start   time.Time
		elapsed time.Duration
		attempt = uint(0)
		reports [][]string
	)

	handleOutcome := func(outcome TestOutcome) (TestOutcome, uint, [][]string, time.Duration) {
		elapsed = elapsed.Truncate(time.Millisecond)

		if outcome == TestPassed {
			reports = nil // discard reports if the test passed
		}

		return outcome, attempt, reports, elapsed
	}

	retry := func() bool {
		if tr.maxAttempts != 0 && attempt >= tr.maxAttempts {
			return false
		}

		if tr.maxDuration != 0 && elapsed >= tr.maxDuration {
			return false
		}

		return true
	}

	start = time.Now()
	for retry() {
		result := TestHelper(tr.fn)
		elapsed = time.Since(start)
		attempt++

		if result.Outcome == TestPassed {
			return handleOutcome(TestPassed)
		}

		reports = append(reports, result.Report)

		time.Sleep(tr.wait)
	}

	return handleOutcome(TestFailed)
}

// FlakyOption is an option function type for an option that modifies
// the behavior of a [FlakyTest]
type FlakyOption func(*flakyRunner)

// MaxAttempts sets the maximum number of attempts for a [FlakyTest].
//
// If the test does not pass within the specified number of attempts, it
// will fail with a report detailing all failed attempts.
//
// The default is 3 attempts.
//
// [MaxAttempts] is ignored if the [MaxDuration] is reached before the
// number of attempts reaches the maximum.
//
// Setting [MaxAttempts] to 0 (zero) will allow the test to run indefinitely
// until it passes, [MaxDuration] is reached, or the 'go test' timeout
// is reached.
func MaxAttempts(n uint) FlakyOption {
	return func(r *flakyRunner) {
		r.maxAttempts = n
	}
}

// MaxDuration sets the maximum duration for a [FlakyTest].
//
// If the test does not pass within the specified duration, it will fail
// with a report detailing all failed attempts.
//
// The default is 1 second.
//
// [MaxDuration] is ignored if [MaxAttempts] is reached before the duration.
//
// Setting [MaxDuration] to 0 (zero) will allow the test to run indefinitely
// until it passes, [MaxAttempts] is reached, or the 'go test' timeout
// is reached.
func MaxDuration(d time.Duration) FlakyOption {
	return func(r *flakyRunner) {
		r.maxDuration = d
	}
}

// WaitBetweenAttempts sets the duration to wait between attempts
// for a [FlakyTest].
//
// The default is 10ms.
func WaitBetweenAttempts(d time.Duration) FlakyOption {
	return func(r *flakyRunner) {
		r.wait = d
	}
}

// FlakyTest creates a runner for a test that may fail intermittently. The test
// will be retried up to a specified number of [MaxAttempts] or until [MaxDuration]
// has passed (whichever occurs first).
//
// If the test passes, no failure report is produced; reports from any failed
// attempts are discarded.
//
// If the test does not pass on any attempt, a report is generated that includes
// the details of all failed attempts.
//
// By default, the test will be attempted up to 3 times for up to 1 second.
//
// You can configure the maximum number of attempts and the maximum duration using
// the [MaxAttempts] and [MaxDuration] options, respectively.  Setting these options
// to 0 (zero) disables each limit.  Setting both options to 0 (zero) allows the
// test to run indefinitely until it passes, the 'go test' timeout is reached, or
// the process is terminated.
//
// Example usage:
//
//	// run a test for a maximum of 5 attempts over a 200 millisecond period
//	Run(FlakyTest("test that fails intermittently", func() {
//		// .. test code ..
//	}, MaxDuration(200*time.Millisecond), MaxAttempts(5)))
func FlakyTest(name string, fn func(), opts ...FlakyOption) flakyRunner {
	const (
		defaultMaxAttempts         = 3
		defaultMaxDuration         = time.Second
		defaultWaitBetweenAttempts = 10 * time.Millisecond
	)

	runner := flakyRunner{
		name:        name,
		fn:          fn,
		maxAttempts: defaultMaxAttempts,
		maxDuration: defaultMaxDuration,
		wait:        defaultWaitBetweenAttempts,
	}

	for _, opt := range opts {
		opt(&runner)
	}

	return runner
}
