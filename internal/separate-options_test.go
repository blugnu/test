package internal_test

import (
	"testing"

	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"

	. "github.com/blugnu/test"
)

func TestSeparateOptions(t *testing.T) {
	With(t)

	Run(Test("returns SubjectName and IsRequired as expect options, others as match options", func() {
		src := []any{
			opt.Name("subject name"),
			opt.IsRequired(true),
			123,
			struct{ A int }{A: 1},
			true,
		}

		expectOpts, matchOpts := internal.SeparateOptions(src)

		Expect(expectOpts).To(EqualSlice([]any{
			opt.Name("subject name"),
			opt.IsRequired(true),
		}))

		Expect(matchOpts).To(EqualSlice([]any{
			123,
			struct{ A int }{A: 1},
			true,
		}))
	}))

	Run(Test("returns only first SubjectName and IsRequired as expect options", func() {
		src := []any{
			opt.Name("subject name"),
			"another subject name",
			opt.IsRequired(true),
			opt.IsRequired(false),
		}

		expectOpts, matchOpts := internal.SeparateOptions(src)

		Expect(expectOpts).To(EqualSlice([]any{
			opt.Name("subject name"),
			opt.IsRequired(true),
		}))

		Expect(matchOpts).Should(BeEmpty())
	}))

	Run(Test("returns a string as SubjectName expectation option", func() {
		src := []any{
			"subject name",
		}

		expectOpts, matchOpts := internal.SeparateOptions(src)

		Expect(expectOpts).To(EqualSlice([]any{
			opt.Name("subject name"),
		}))

		Expect(matchOpts).Should(BeEmpty())
	}))
}
