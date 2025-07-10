package mocks

import (
	"fmt"
	"strings"

	"github.com/blugnu/test/test"
)

// Matcher is a struct that implements the test.Matcher interface for
// handling mock expectations in tests. To support a variety of mock
// implementations, it is implemented as an AnyMatcher, accepting
// an any type for the `got` parameter.
//
// It checks for the implementation of supported interfaces to determine
// whether the mock expectations were met.  If no supported interface
// is found, the matcher will fail the test as invalid.
//
// # Supported Interfaces
//
//   - ExpectationsWereMet() error
//   - Reset()
//
// # Reset()
//
// If the mock implements a Reset() function, this will be called
// regardless of whether the expectations were met or not. i.e. testing
// mock expectations will also reset the mock (where supported).
type Matcher struct {
	err     error
	got     any
	isValid bool
}

// HasExpectationsWereMet is an interface that defines a method
// ExpectationsWereMet() that returns an error if the expectations
// were not met.
//
// This is an interface that many mock libraries implement, such as
// DATA-DOG/sqlmock.
type HasExpectationsWereMet interface {
	ExpectationsWereMet() error
}

// HasReset is an interface that defines a method Reset() that
// resets the state of the mock.
type HasReset interface {
	Reset()
}

// Match implements the AnyMatcher interface
func (m *Matcher) Match(got any, opts ...any) bool {
	handle := func(err error) bool {
		m.err = err
		m.isValid = true

		if mock, ok := got.(HasReset); ok {
			mock.Reset()
		}

		return m.err == nil
	}

	m.got = got

	switch mock := got.(type) {
	case HasExpectationsWereMet:
		return handle(mock.ExpectationsWereMet())

	default:
		return false
	}
}

// OnTestFailure implements test failure report for the Matcher
func (m *Matcher) OnTestFailure(subject any, opts ...any) []string {
	if !m.isValid {
		test.T().Helper()
		test.Invalid(fmt.Sprintf("mocks.Matcher: %T does not implement a supported mock interface", m.got))
	}

	return strings.Split(m.err.Error(), "\n")
}
