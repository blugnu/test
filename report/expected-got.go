package report

// Expected returns a slice of strings representing an expected value
// in a failure report.
//
// The returned slice will contain a single string in the format
// "expected: expectedValue".
func Expected(expected string) []string {
	return []string{
		"expected: " + expected,
	}
}

// ExpectedGot returns a slice of strings representing an expected value
// and an actual value in a failure report.
//
// The returned slice will contain two strings in the format:
//
//	expected: expectedValue
//	got     : gotValue
func ExpectedGot(expected, got string) []string {
	return []string{
		"expected: " + expected,
		"got     : " + got,
	}
}
