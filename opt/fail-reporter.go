package opt

import (
	"fmt"
	"strings"
)

// FailReporter is a function type that implements an OnTestFailure function.
//
// This option allows a test to override the default error report when a test
// fails, replacing it with a custom message.  All options supplied to the
// matcher are also passed to the function.  The custom function must accept
// the options as a variadic number of arguments, but may choose to ignore them.
//
//	Expect(got).To(BeTrue(), opt.FailureReport(func(...any) []string {
//		return []string{"custom failure message"}
//	}))
type FailReporter func(opts ...any) []string

// OnTestFailure implements a custom test failure report function, calling
// the receiver FailReporter with the provided options.
func (f FailReporter) OnTestFailure(opts ...any) []string {
	return f(opts...)
}

// OnFailure returns a FailReporter appropriate to the type of the provided report.
//
// If report is a string, it is used as a single-line failure message. If the string
// contains newline characters the string will be split into multiple lines.
//
// If report is a []string, it is used as a multi-line failure report. Newlines in
// the report strings are preserved.
//
// If report is a FailReporter (or function of the same signature), it is returned
// unchanged.
//
// If any other type is provided, the current test is failed with a message
// indicating the value of the provided report, similar to:
//
//	"test failed with: <report>"
func OnFailure(report any) FailReporter {
	reportString := func(msg string) FailReporter {
		return func(...any) []string {
			return []string{msg}
		}
	}

	reportLines := func(lines ...string) FailReporter {
		return func(...any) []string {
			return lines
		}
	}

	switch r := report.(type) {
	case string:
		if strings.Contains(r, "\n") {
			return reportLines(strings.Split(r, "\n")...)
		}
		return reportString(r)
	case []string:
		return reportLines(r...)
	case FailReporter:
		return r
	case func(...any) []string:
		return FailReporter(r)
	default:
		return reportString(fmt.Sprintf("test failed with: %v", report))
	}
}

// MARK: deprecated

// Deprecated: use FailMessage, FailReport or FailReporter instead.
type FailureReport = FailReporter
