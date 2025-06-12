package testframe

import "errors"

var (
	// Test frame errors
	ErrEmptyStack            = errors.New("empty stack for this goroutine")
	ErrNoStack               = errors.New("error determining goroutine id: unable to obtain stack information for goroutine")
	ErrNoTestFrame           = errors.New("no test frame; did you forget to call With(t)?")
	ErrUnexpectedStackFormat = errors.New("unable to determine goroutine id: unexpected stack format")
)
