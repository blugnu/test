package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestT(t *testing.T) {
	With(t)

	Run(Test("valid testframe", func() {
		t1 := test.T()
		Expect(t1).IsNotNil()
	}))

	Run(Test("no testframe", func() {
		defer Expect(Panic()).DidNotOccur()

		With(test.NilFrame())

		t := test.T()
		if t == nil {
			panic("T() returned nil, expected a valid test.Helper")
		}

		// coverage for the noopHelper
		t.Helper()
	}))
}
