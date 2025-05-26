package test

import (
	"errors"
	"strings"
)

var (
	// Test frame errors
	ErrNoStack               = errors.New("unable to obtain stack information for goroutine; unable to determine goroutine id")
	ErrUnexpectedStackFormat = errors.New("unexpected stack format; unable to determine goroutine id")

	// Test runner errors
	ErrNoTestFrame = errors.New("no test frame; did you forget to call With(t)?")

	// Mock and Fake errors
	ErrExpectationsNotMet = errors.New("expectations not met")
	ErrExpectedArgs       = errors.New("arguments were expected but not recorded")
	ErrNoResultForArgs    = errors.New("no result for arguments")
	ErrUnexpectedArgs     = errors.New("the arguments recorded did not match those expected")
	ErrUnexpectedCall     = errors.New("unexpected call")
	ErrResultNotUsed      = errors.New("result not used")

	// Other errors
	ErrCapture          = errors.New("test.Capture")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrNotNilable       = errors.New("values of this type are not nilable")
)

func invalidTest(msg ...string) {
	t := T()
	t.Helper()
	s := strings.Join(msg, "\n")
	if len(s) > 0 {
		s = "\n" + s
	}
	t.Error("<== INVALID TEST" + s)
}
