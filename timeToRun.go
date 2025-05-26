package test

import "time"

// RunWithTimer executes the provided function and returns the duration it took to execute.
func RunWithTimer(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}
