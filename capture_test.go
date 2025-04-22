package test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestCapture(t *testing.T) {
	With(t)

	// ARRANGE
	writeOutput := func() {
		os.Stdout.WriteString("to stdout (1)\n")
		os.Stdout.WriteString("to stdout (2)\n")
		os.Stderr.WriteString("to stderr (1)\n")
		os.Stderr.WriteString("to stderr (2)\n")
		log.Println("to log")
	}

	// ACT
	stdout, stderr := CaptureOutput(t, func() {
		writeOutput()
	})

	// ASSERT
	Expect(stdout).To(EqualSlice([]string{
		"to stdout (1)",
		"to stdout (2)",
	}))
	Expect(stderr).To(ContainStrings([]string{
		"to stderr (1)",
		"to stderr (2)",
		"to log",
	}))

	t.Run("when capture fails", func(t *testing.T) {
		With(t)

		// ARRANGE
		cpyerr := fmt.Errorf("copy error")
		og := copy
		defer func() { copy = og }()
		copy = func(dst io.Writer, src io.Reader) (int64, error) { _, _ = io.Copy(dst, src); return 0, cpyerr }
		defer ExpectPanic(ErrCapture).Assert()

		stdout, stderr := CaptureOutput(t, func() { writeOutput() })

		// ASSERT
		ExpectEmpty(stdout)
		ExpectEmpty(stderr)
	})
}
