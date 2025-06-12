package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestNilFrame(t *testing.T) {
	With(t)

	Run("with no test frame", func() {
		defer Expect(Panic(test.ErrNoTestFrame)).DidOccur()
		With(test.NilFrame())

		test.NilFrame()
	})
}
