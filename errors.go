package test

import (
	"errors"
)

var (
	// general errors
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrInvalidOperation = errors.New("invalid operation")

	// mock and fake errors
	ErrExpectationsNotMet = errors.New("expectations not met")
	ErrExpectedArgs       = errors.New("arguments were expected but not recorded")
	ErrNoResultForArgs    = errors.New("no result for arguments")
	ErrUnexpectedArgs     = errors.New("the arguments recorded did not match those expected")
	ErrUnexpectedCall     = errors.New("unexpected call")
	ErrResultNotUsed      = errors.New("result not used")

	// recording errors
	ErrRecordingFailed                 = errors.New("recording failed")
	ErrRecordingStdout                 = errors.New("error recording stdout")
	ErrRecordingStderr                 = errors.New("error recording stderr")
	ErrRecordingUnableToRedirectLogger = errors.New("unable to redirect logger output")
)
