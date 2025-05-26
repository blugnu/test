package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestTestFailureFunc(t *testing.T) {
	With(t)

	impl := opt.FailureReport(func(...any) []string {
		return []string{"test"}
	})

	result := impl.OnTestFailure()
	Expect(result).To(EqualSlice([]string{"test"}))
}

func TestOnFailure(t *testing.T) {
	With(t)

	impl := opt.OnFailure("custom message")

	result := impl.OnTestFailure()
	Expect(result).To(EqualSlice([]string{"custom message"}))
}
