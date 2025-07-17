package opt_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestIntervalClosure_String(t *testing.T) {
	With(t)

	Expect(opt.IntervalOpen.String()).
		To(Equal("IntervalOpen: min < x < max"))

	Expect(opt.IntervalOpenMin.String()).
		To(Equal("IntervalOpenMin: min < x <= max"))

	Expect(opt.IntervalOpenMax.String()).
		To(Equal("IntervalOpenMax: min <= x < max"))

	Expect(opt.IntervalClosed.String()).
		To(Equal("IntervalClosed: min <= x <= max"))

	Expect(opt.IntervalClosure(100).String()).
		To(Equal("IntervalClosure(100)"))
}
