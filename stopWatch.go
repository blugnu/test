package test

import "time"

// StopWatch executes the provided function and returns the
// duration it took to execute.
func StopWatch(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}
