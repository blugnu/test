package testcase

import (
	"testing"

	"github.com/blugnu/test/internal/testframe"
)

// TestingT is an interface that describes the methods of a testing.T
// instance that are used in this package.
type TestingT interface {
	Errorf(string, ...any)
	Helper()
	Parallel()
	Run(string, func(t *testing.T)) bool
}

// GetT() is a helper function to retrieve the current TestingT instance.
// It is used to ensure that the test helper is called correctly.
func GetT() TestingT {
	return testframe.MustPeek[TestingT]()
}
