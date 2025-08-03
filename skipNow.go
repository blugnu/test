package test

import "testing"

// SkipNow is a function that is used to skip the current test immediately.
// It is a wrapper around the SkipNow method of the current test frame.
//
// If the current testframe is a *testing.T (the usual case), the supplied
// reason for skipping the test is logged. Otherwise, the reason is required
// for documentation purposes only.
func SkipNow(reason string) {
	t := T()
	t.Helper()

	if t, ok := t.(*testing.T); ok {
		t.Log("<== SKIPPED:", reason)
	}

	t.SkipNow()
}
