package test //nolint: testpackage // tests private types and functions

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestRecord(t *testing.T) {
	With(t)

	// ARRANGE
	writeOutput := func() {
		_, _ = os.Stdout.WriteString("to stdout (1)\n")
		_, _ = os.Stdout.WriteString("to stdout (2)\n")
		_, _ = os.Stderr.WriteString("to stderr (1)\n")
		_, _ = os.Stderr.WriteString("to stderr (2)\n")
		log.Println("to log")
	}

	// ACT
	stdout, stderr := Record(func() {
		writeOutput()
	})

	// ASSERT
	Expect(stdout).To(EqualSlice([]string{
		"to stdout (1)",
		"to stdout (2)",
	}), strings.Contains)
	Expect(stderr).To(ContainSlice([]string{
		"to stderr (1)",
		"to stderr (2)",
		"to log",
	}), strings.Contains)

	t.Run("when copy fails", func(t *testing.T) {
		With(t)

		// ARRANGE
		cpyerr := errors.New("copy error")
		sut := stdioCapture{
			copy: func(dst io.Writer, src io.Reader) (int64, error) {
				return 0, cpyerr
			},
			getLogOutput: getLogOutput,
		}

		defer Expect(Panic(ErrRecordingFailed)).DidOccur()

		// ACT
		// This will panic because the copy function is mocked to return an error
		stdout, stderr := record(sut, func() { writeOutput() })

		// ASSERT
		Expect(stdout, "stdout").Should(BeEmpty())
		Expect(stderr, "stderr").Should(BeEmpty())
	})

	t.Run("when unable to get log output", func(t *testing.T) {
		With(t)

		// ARRANGE
		sut := stdioCapture{
			copy: io.Copy,
			getLogOutput: func() (io.Writer, bool) {
				return nil, false // simulate failure to get log output
			},
		}

		defer Expect(Panic(ErrRecordingUnableToRedirectLogger)).DidOccur()

		// ACT
		// This will panic because the copy function is mocked to return an error
		stdout, stderr := record(sut, func() { writeOutput() })

		// ASSERT
		Expect(stdout, "stdout").Should(BeEmpty())
		Expect(stderr, "stderr").Should(BeEmpty())
	})

	RunParallelScenarios(func(_ *string, _ int) {
		defer Expect(Panic(ErrInvalidOperation)).DidOccur()

		_ = Test(func() {})
	}, []string{"not used"})

	RunParallel("in a RunParallel test", func() {
		defer Expect(Panic(ErrInvalidOperation)).DidOccur()

		_ = Test(func() {})
	})

	Run("in a Parallel() test", func() {
		defer Expect(Panic(ErrInvalidOperation)).DidOccur()

		Parallel()
		_ = Test(func() {})
	})
}

func ExampleRecord() {
	// remove date and time from log output otherwise this
	// example will fail because the date and time will be
	// different each time it is run!
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// capture the output of a function that explicitly writes
	// to stdout, stderr and emits logs (which also go to stderr,
	// by default)
	stdout, stderr := Record(func() {
		fmt.Println("to stdout")
		fmt.Fprintln(os.Stderr, "to stderr")
		log.Println("to log")
	})

	// write what was captured to stdout (for the Example to test)
	// (in a test, you would use Expect() to test the output)

	fmt.Println("captured stdout:")
	for i, s := range stdout {
		fmt.Printf("  %d: %s\n", i+1, s)
	}

	fmt.Println("captured stderr:")
	for i, s := range stderr {
		fmt.Printf("  %d: %s\n", i+1, s)
	}

	// Output:
	// captured stdout:
	//   1: to stdout
	// captured stderr:
	//   1: to stderr
	//   2: to log
}
