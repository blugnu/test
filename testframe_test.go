package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

func TestGetT(t *testing.T) {
	With(t)

	t1 := GetT()
	Expect(t1).IsNotNil()
}

func TestT(t *testing.T) {
	With(t)

	t1 := T()
	Expect(t1).IsNotNil()
}

func TestT_Name(t *testing.T) {
	With(t)

	var s TestingT
	Run(Test("Subtest", func() {
		s = T()
	}))

	Expect(t.Name()).To(Equal("TestT_Name"))
	Expect(s.Name()).To(Equal("TestT_Name/Subtest"))
}

func TestWith(t *testing.T) {
	With(t)

	// pushing nil is invalid and should panic
	Run(Test("nil", func() {
		// the panic test must be setup before niling the test frame
		defer Expect(Panic(testframe.ErrNoTestFrame)).DidOccur()

		With(nil)
	}))

	// to test a simulated lack of test frame, we can use test.NilFrame()
	Run(Test("test.NilFrame", func() {
		// the panic test must be setup before pushing the test frame
		defer Expect(Panic(testframe.ErrNoTestFrame)).DidOccur()

		// pushing test.NilFrame() to simulate a nil test frame
		With(test.NilFrame())

		// this should panic with ErrNoTestFrame
		T()

		// so this should not be reached at all
		test.Warning("this should not be reached")
	}))

	T().Run("in a testing.T subtest", func(st *testing.T) {
		With(st)

		t := T()
		Expect(t).ShouldNot(BeNil())
		Expect(t.Name()).To(Equal("TestWith/in_a_testing.T_subtest"))
	})

	// verify that the testframe stack is consistent
	Expect(T().Name()).To(Equal("TestWith"))
}
