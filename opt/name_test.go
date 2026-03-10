package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestHasName(t *testing.T) {
	With(t)

	Run(Test("options contain a Name", func() {
		result, ok := opt.GetName([]any{0, 1, opt.Name("test"), true})
		Expect(result).To(Equal(opt.Name("test")))
		Expect(ok).To(BeTrue())
	}))

	Run(Test("options contain a string", func() {
		result, ok := opt.GetName([]any{0, 1, "test", true})
		Expect(result).To(Equal(opt.Name("test")))
		Expect(ok).To(BeTrue())
	}))

	Run(Test("options contain a Name and string with Name first", func() {
		result, ok := opt.GetName([]any{0, 1, opt.Name("test"), "test2", true})
		Expect(result).To(Equal(opt.Name("test")))
		Expect(ok).To(BeTrue())
	}))

	Run(Test("options contain a Name and string with string first", func() {
		result, ok := opt.GetName([]any{0, 1, "test2", opt.Name("test"), true})
		Expect(result).To(Equal(opt.Name("test2")))
		Expect(ok).To(BeTrue())
	}))

	Run(Test("options does not contain a string or Name", func() {
		result, ok := opt.GetName([]any{0, 1, true})
		Expect(result).To(Equal(opt.Name("")))
		Expect(ok).To(BeFalse())
	}))
}

func TestName_String(t *testing.T) {
	With(t)

	Run(Test("returns the specified string as a Name", func() {
		result := opt.Name("test").String()
		Expect(result).To(Equal("test"))
	}))
}

func TestNamef(t *testing.T) {
	With(t)

	Run(Test("returns the formatted string as a Name", func() {
		result := opt.Namef("test %d", 42)
		Expect(result).To(Equal(opt.Name("test 42")))
	}))
}

func TestWithName(t *testing.T) {
	With(t)

	Run(Test("adds Name if not present", func() {
		opts := []any{opt.IsRequired(true)}
		newOpts := opt.WithName(opts, "test subject")

		Expect(newOpts).To(EqualSlice([]any{opt.IsRequired(true), opt.Name("test subject")}))
	}))

	Run(Test("does not add Name if already present", func() {
		opts := []any{opt.IsRequired(true), opt.Name("existing subject")}
		newOpts := opt.WithName(opts, "test subject")

		Expect(newOpts).To(EqualSlice(opts))
	}))

	Run(Test("does not add Name if empty", func() {
		opts := []any{opt.IsRequired(true)}
		newOpts := opt.WithName(opts, "")

		Expect(newOpts).To(EqualSlice(opts))
	}))
}

func TestWithNamef(t *testing.T) {
	With(t)

	Run(Test("adds formatted Name if not present", func() {
		opts := []any{opt.IsRequired(true)}
		newOpts := opt.WithNamef(opts, "subject %d", 1)

		Expect(newOpts).To(EqualSlice([]any{opt.IsRequired(true), opt.Name("subject 1")}))
	}))

	Run(Test("does not add formatted Name if already present", func() {
		opts := []any{opt.IsRequired(true), opt.Name("existing subject")}
		newOpts := opt.WithNamef(opts, "subject %d", 1)

		Expect(newOpts).To(EqualSlice(opts))
	}))

	Run(Test("does not add formatted Name if empty", func() {
		opts := []any{opt.IsRequired(true)}
		newOpts := opt.WithNamef(opts, "")

		Expect(newOpts).To(EqualSlice(opts))
	}))
}
