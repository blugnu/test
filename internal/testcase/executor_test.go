package testcase_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testcase"
)

func TestNewExecutor(t *testing.T) {
	With(t)

	Run(Test("unsupported function", func() {
		result := TestHelper(func() {
			_ = testcase.NewExecutor[int](func() {})
		})

		result.ExpectWarning("func() is not a valid test function")
	}))
}

func TestExecutor_Execute(t *testing.T) {
	With(t)

	Run(Test("no function configured", func() {
		result := TestHelper(func() {
			data := 42
			exec := testcase.Executor[int]{}
			exec.Execute("not used", data)
		})

		result.ExpectInvalid("a test function must be provided")
	}))
}
