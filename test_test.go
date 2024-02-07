package test

import (
	"path"
	"runtime"
	"testing"
)

// alias for *testing.T to simplify the type signature of testcase functions
// in data-driven tests
type T = *testing.T

// returns the name of the file (without path) containing the calling
// function.  Used to ensure that test failure reports identify the
// correct calling file and that *testing.T.Helper() calls are being
// made as required to ensure this.
func currentFilename() string {
	_, filepath, _, _ := runtime.Caller(2)
	_, filename := path.Split(filepath)
	return filename
}
