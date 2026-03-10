package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestUnset(t *testing.T) {
	With(t)

	opts := []any{"a", "b", "c"}

	// Test the Unset function with various scenarios
	Expect(opt.Unset(opts, "b")).To(EqualSlice([]any{"a", "c"}))
}
