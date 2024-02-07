package test

import "errors"

var (
	ErrCapture           = errors.New("test.Capture")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInvalidOptionType = errors.New("invalid option type")
	ErrInvalidTest       = errors.New("invalid test")
	ErrNotNilable        = errors.New("values of this type are not nilable")
)
