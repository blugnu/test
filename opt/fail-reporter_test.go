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

	Run(Test("string input", func() {
		impl := opt.OnFailure("custom message")

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"custom message"}))
	}))

	Run(Test("string input with newline", func() {
		impl := opt.OnFailure("line1\nline2")

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"line1", "line2"}))
	}))

	Run(Test("[]string input", func() {
		impl := opt.OnFailure([]string{"line1", "line2"})

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"line1", "line2"}))
	}))

	Run(Test("[]string input with newline", func() {
		impl := opt.OnFailure([]string{"line1\nalso line1", "line2"})

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"line1\nalso line1", "line2"}))
	}))

	Run(Test("FailReporter input", func() {
		orig := opt.FailReporter(func(a ...any) []string { return []string{"line1", "line2"} })
		impl := opt.OnFailure(orig)

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"line1", "line2"}))
	}))

	Run(Test("function input", func() {
		fn := func(...any) []string {
			return []string{"line1", "line2"}
		}
		impl := opt.OnFailure(fn)

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"line1", "line2"}))
	}))

	Run(Test("other input", func() {
		impl := opt.OnFailure(42)

		result := impl.OnTestFailure()
		Expect(result).To(EqualSlice([]string{"test failed with: 42"}))
	}))
}
