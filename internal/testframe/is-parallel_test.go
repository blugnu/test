package testframe_test

import (
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

func TestIsParallel(t *testing.T) {
	With(t)

	Run(Test("not parallel", func() {
		Expect(testframe.IsParallel()).To(BeFalse())
		T().Parallel()
		Expect(testframe.IsParallel()).To(BeTrue())
	}))

	Run(Test("parallel", func() {
		T().Parallel()
		Expect(testframe.IsParallel()).To(BeTrue())
	}))

	Run(Test("example runner is always non-parallel", func() {
		test.Example()
		defer testframe.Pop()

		T().Parallel() // NO-OP in example mode

		Expect(testframe.IsParallel()).To(BeFalse())
	}))
}
