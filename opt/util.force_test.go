package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestForce(t *testing.T) {
	With(t)

	Run(Test("replaces options of the same type with specified option at the end", func() {
		opts := opt.Force([]any{1, "old string", 2}, "new string")

		Expect(opts).To(DeepEqual([]any{1, 2, "new string"}))
	}))

	Run(Test("places the specified option at the start when ForceAtStart is used", func() {
		opts := opt.Force([]any{1, "old string", 2}, "new string", opt.ForceAtStart)

		Expect(opts).To(DeepEqual([]any{"new string", 1, 2}))
	}))

	Run(Test("adds to original options when there are no existing options of the same type", func() {
		original := []any{1, 2, 3}
		opts := opt.Force(original, "new string")

		Expect(opts).To(DeepEqual([]any{1, 2, 3, "new string"}))
	}))
}
