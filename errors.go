package test

import (
	"errors"
)

var (
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
