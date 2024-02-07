package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

var copy = io.Copy

// converts a string to a slice of strings, splitting on newlines.
func toSlice(s string) []string {
	if l := strings.Split(s, "\n"); len(l) > 1 || (len(l) == 1 && l[0] != "") {
		if l[len(l)-1:][0] == "" {
			l = l[:len(l)-1]
		}
		return l
	}
	return nil
}

// sets up the redirect of a *os.File (such as os.Stdout or os.Stderr).
//
// The function returns:
//   - a function that must be called to restore the original *os.File,
//   - a function that must be called to complete the redirect and
//     return any captured output.
//
// Example:
//
//	  func DoSomething() {
//		restore, close := redirect(&os.Stdout)
//		defer restore()
//
//		fmt.Println("some output")
//		s, err := close()
//
//		fmt.Println(s) // "some output"
//	  }
func redirect(t **os.File) (func(), func() (string, error)) {
	og := *t
	r, w, _ := os.Pipe()
	*t = w

	c := make(chan string)
	e := make(chan error)
	go func() {
		var buf bytes.Buffer
		_, err := copy(&buf, r)
		c <- buf.String()
		e <- err
	}()

	return func() { *t = og }, func() (string, error) { w.Close(); return <-c, <-e }
}

// captures the output of a function that may write to stdout or stderr and which
// may emit logs using the default standard library logger.
//
// The function returns:
//   - a *Capture struct containing any captured output and
//   - any error(s) in the capture process failed.
func capture(fn func()) (stdout []string, stderr []string, err error) {

	unhookStdout, completeStdout := redirect(&os.Stdout)
	defer unhookStdout()

	unhookStderr, completeStderr := redirect(&os.Stderr)
	defer unhookStderr()

	// The standard logger is initialised with os.Stderr before
	// the capture is setup. This means that the standard logger
	// will not write to the captured stderr.
	//
	// We need to replace the standard logger's output with the
	// captured stderr, which we can do using log.SetOutput().
	// But we also need to restore it afterwards so we need to
	// retrieve the original output (we cannot just assume it is
	// the current value of os.Stderr).
	//
	// There is no direct mechanism to obtain the current default
	// logger's output, but we can use reflection to obtain the
	// address of the output field and then use unsafe to obtain
	// its value.

	v := reflect.ValueOf(log.Default()).Elem().FieldByName("out")
	lo := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Interface().(io.Writer)
	defer func() { log.SetOutput(lo) }()

	log.SetOutput(os.Stderr) // os.Stderr is redirected at this point

	// call the function that may write to stdout, stderr and log
	fn()

	// check the captured output for any errors (in the capture process)
	var (
		s    string
		errs = []error{}
	)
	if s, err = completeStdout(); err != nil {
		s = ""
		errs = append(errs, fmt.Errorf("%w: Stdout: %w", ErrCapture, err))
	}
	stdout = toSlice(s)

	if s, err = completeStderr(); err != nil {
		s = ""
		errs = append(errs, fmt.Errorf("%w: Stderr: %w", ErrCapture, err))
	}
	stderr = toSlice(s)

	return stdout, stderr, errors.Join(errs...)
}

// captures the stdout and stderr output of a test helper or other
// function. If the capture mechanism fails the function will panic.
//
// # Example:
//
//	  func TestSomething(t *testing.T) {
//		// ACT
//		stdout, _ :=test.CaptureOutput(t, func () {
//		   return doSomething()
//		})
//
//		// ASSERT
//		test.IsNil(t, err)
//		stdout.Contains("expected output from doSomething")
//	  }
func CaptureOutput(t *testing.T, fn func()) (stdout StringsTest, stderr StringsTest) {
	so, se, err := capture(fn)
	if err != nil {
		panic(fmt.Errorf("test.CaptureOutput: unexpected error: %w", err))
	}
	return Strings(t, so, "stdout"), Strings(t, se, "stderr")
}
