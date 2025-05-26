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
	"unsafe"
)

// record is the internal function that captures the output of stdout and stderr
// using a specified stdioCapture struct. This enables a mock implementation
// to be used in tests.
func record(cap stdioCapture, fn func()) (stdout []string, stderr []string) {
	stdout, stderr, err := cap.execute(fn)
	if err != nil {
		panic(fmt.Errorf("test.Capture: unexpected error: %w", err))
	}
	return
}

// Record captures the stdout and stderr output resulting from the execution
// of some function.
//
// In the unlikely event that the mechanism fails, the function will
// panic to avoid returning misleading results or require error handling.
func Record(fn func()) (stdout []string, stderr []string) {
	if IsParallel() {
		panic(fmt.Errorf("%w: test.Record cannot be used in a parallel test", ErrInvalidOperation))
	}
	return record(stdioCapture{io.Copy}, fn)
}

// stdioCapture is a struct that captures the output of stdout and stderr
// using a specified copy function.
type stdioCapture struct {
	copy func(dst io.Writer, src io.Reader) (written int64, err error)
}

// redirect sets up the redirect of a *os.File (such as os.Stdout or os.Stderr).
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
func (cap stdioCapture) redirect(t **os.File) (func(), func() (string, error)) {
	og := *t
	r, w, _ := os.Pipe()
	*t = w

	c := make(chan string)
	e := make(chan error)
	go func() {
		var buf bytes.Buffer
		_, err := cap.copy(&buf, r)
		c <- buf.String()
		e <- err
	}()

	return func() { *t = og }, func() (string, error) { w.Close(); return <-c, <-e }
}

// execute performs a provided function, capturing any output to stdout or stderr
// produced by that function.
func (cap stdioCapture) execute(fn func()) (stdout []string, stderr []string, err error) {
	unhookStdout, completeStdout := cap.redirect(&os.Stdout)
	defer unhookStdout()

	unhookStderr, completeStderr := cap.redirect(&os.Stderr)
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
	// the value.

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
	stdout = cap.toSlice(s)

	if s, err = completeStderr(); err != nil {
		s = ""
		errs = append(errs, fmt.Errorf("%w: Stderr: %w", ErrCapture, err))
	}
	stderr = cap.toSlice(s)

	return stdout, stderr, errors.Join(errs...)
}

// converts a string to a slice of strings, splitting on newlines.
func (stdioCapture) toSlice(s string) []string {
	if l := strings.Split(s, "\n"); len(l) > 1 || (len(l) == 1 && l[0] != "") {
		for len(l) > 0 && l[len(l)-1] == "" {
			l = l[:len(l)-1]
		}
		return l
	}
	return nil
}
