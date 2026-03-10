package test

import (
	"errors"
)

// general errors
var (
	// ErrInvalidArgument indicates that an argument provided to a function or method is invalid.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrInvalidOperation indicates that an attempted operation is invalid in the current context.
	ErrInvalidOperation = errors.New("invalid operation")
)

// mock and fake errors
var (
	// ErrExpectationsNotMet indicates that the expectations set on a mock or fake
	// were not met.
	ErrExpectationsNotMet = errors.New("expectations not met")

	// ErrExpectedArgs indicates that arguments were expected but not recorded.
	ErrExpectedArgs = errors.New("arguments were expected but not recorded")

	// ErrNoResultForArgs indicates that no result was found for the given arguments.
	ErrNoResultForArgs = errors.New("no result for arguments")

	// ErrUnexpectedArgs indicates that the arguments recorded did not match those expected.
	ErrUnexpectedArgs = errors.New("the arguments recorded did not match those expected")

	// ErrUnexpectedCall indicates that a call was made that was not expected.
	ErrUnexpectedCall = errors.New("unexpected call")

	// ErrResultNotUsed indicates that a result was provided but not used.
	ErrResultNotUsed = errors.New("result not used")
)

// recording errors
var (
	// ErrRecordingFailed indicates that an error occurred while trying to
	// record output.
	ErrRecordingFailed = errors.New("recording failed")

	// ErrRecordingStdout indicates that an error occurred while trying to
	// record stdout.
	ErrRecordingStdout = errors.New("error recording stdout")

	// ErrRecordingStderr indicates that an error occurred while trying to
	// record stderr.
	ErrRecordingStderr = errors.New("error recording stderr")

	// ErrFailedToRedirectLogger indicates that an error occurred while trying
	// to redirect logger output.
	ErrFailedToRedirectLogger = errors.New("failed to redirect logger output")

	// Deprecated: this error is deprecated and will be removed in a future version;
	// use [ErrFailedToRedirectLogger] instead.
	ErrRecordingUnableToRedirectLogger = ErrFailedToRedirectLogger
)
