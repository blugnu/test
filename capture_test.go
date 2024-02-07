package test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestCapture(t *testing.T) {
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
	Equal(t, stdout.got[0], "to stdout (1)")
	Equal(t, stdout.got[1], "to stdout (2)")
	Equal(t, stderr.got[0], "to stderr (1)")
	Equal(t, stderr.got[1], "to stderr (2)")
	// TODO: String(stderr.got[2]).Contains("to log")

	t.Run("when capture fails", func(t *testing.T) {
		// ARRANGE
		cpyerr := fmt.Errorf("copy error")
		og := copy
		defer func() { copy = og }()
		copy = func(dst io.Writer, src io.Reader) (int64, error) { _, _ = io.Copy(dst, src); return 0, cpyerr }
		defer ExpectPanic(ErrCapture).Assert(t)

		stdout, stderr := CaptureOutput(t, func() { writeOutput() })

		// ASSERT
		stdout.IsEmpty()
		stderr.IsEmpty()
	})
}
