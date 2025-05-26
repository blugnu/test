package test

import (
	"fmt"
	"regexp"

	"github.com/blugnu/test/matchers/strings"
)

func ContainString(expected string) strings.ContainsMatch {
	if expected == "" {
		panic(fmt.Errorf("ContainString: %w: empty string is not valid", ErrInvalidArgument))
	}

	return strings.ContainsMatch{Expected: expected}
}

func MatchRegEx(regex string) strings.RegExMatch {
	if regex == "" {
		panic(fmt.Errorf("Match: %w: empty string is not valid; a valid regex must be provided", ErrInvalidArgument))
	}

	ex, err := regexp.Compile(regex)
	if err != nil {
		panic(fmt.Errorf("invalid regex: %w: %w", ErrInvalidArgument, err))
	}

	return strings.RegExMatch{Expected: ex}
}
