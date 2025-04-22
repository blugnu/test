package test

import (
	"errors"
	"fmt"
)

type ErrorTests struct {
	got error
}

func (et ErrorTests) Is(target error) {
	T().Helper()
	Expect(errors.Is(et.got, target)).To(BeTrue(), CustomOneLineReportFunc(func() string {
		return fmt.Sprintf("expected error: %v\ngot           : %v", target, et.got)
	}))
}

func ExpectError(got any, opts ...any) ErrorTests {
	t, _ := TestFrame()
	if t == nil {
		panic("ExpectError: no test frame; did you forget to call With(t)?")
	}

	t.Helper()

	err, ok := got.(error)
	if !ok && got != nil {
		t.Fatalf("expected error, got %T", got)
	}

	return ErrorTests{got: err}
}
