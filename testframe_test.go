package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestGetT(t *testing.T) {
	With(t)

	Run("GetT", func() {
		t1 := GetT()
		Expect(t1).IsNotNil()
	})
}

func TestT(t *testing.T) {
	With(t)
	Run("T", func() {
		t1 := T()
		Expect(t1).IsNotNil()
	})
}

func TestT_Name(t *testing.T) {
	With(t)

	Run("Subtest", func() {
		s := T()
		Expect(t.Name()).To(Equal("TestT_Name"))
		Expect(s.Name()).To(Equal("TestT_Name/Subtest"))
	})
}

func TestWith(t *testing.T) {
	With(t)

	// pushing nil is invalid and should panic
	Run("nil", func() {
		// the panic test must be setup before niling the test frame
		defer Expect(Panic(testframe.ErrNoTestFrame)).DidOccur()

		With(nil)
	})

	// to test a simulated lack of test frame, we can use a test.NilFrame()
	Run("test.NilFrame()", func() {
		// the panic test must be setup before pushing the test frame
		defer Expect(Panic(testframe.ErrNoTestFrame)).DidOccur()

		// pushing test.NilFrame() to simulate a nil test frame
		With(test.NilFrame())

		// this should panic with ErrNoTestFrame
		T()

		// so this should not be reached at all
		Expect(false).To(BeTrue(), opt.OnFailure("this test should not be reached"))
	})

	T().Run("in an explicit subtest", func(st *testing.T) {
		With(st)

		t := T()
		Expect(t).ShouldNot(BeNil())
		Expect(t.Name()).To(Equal("TestWith/in_an_explicit_subtest"))
	})

	// verify that the testframe stack is consistent
	Expect(T().Name()).To(Equal("TestWith"))
}
