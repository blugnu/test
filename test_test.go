package test

import (
	"path"
	"runtime"
)

// currentFilename returns the name of the file containing the
// calling function.
//
// This is used to ensure that test failure reports identify the correct calling
// file (*_test.go).  This avoids replicating the filename as a literal which
// could result in tests failing due to the test file itself being renamed
// (which would be annoying and inappropriate).
func currentFilename() string {
	_, file, _, _ := runtime.Caller(1)
	_, filename := path.Split(file)
	return filename
}
