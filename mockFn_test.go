package test

import (
	"errors"
	"testing"
)

func TestMockFnRecordCall(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "mapped results are configured",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{},
				}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				_, _ = sut.RecordCall(42)
			},
		},
		{scenario: "unexpected call",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result, err := sut.RecordCall()

				// ASSERT
				That(t, result).Equals(0)
				Error(t, err).Is(ErrUnexpectedCall)
			},
		},
		{scenario: "unexpected arguments",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: AddressOf(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(0)

				// ASSERT
				That(t, result).Equals(84)
				Error(t, err).Is(ErrUnexpectedArgs)
			},
		},
		{scenario: "further expectations",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: AddressOf(42), result: 84},
						{args: AddressOf(43), result: 86},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(42)

				// ASSERT
				Error(t, err).IsNil()
				That(t, result).Equals(84)
				That(t, sut.expected).Equals(sut.expectations[1])
				That(t, sut.idxExpected).Equals(1)
			},
		},
		{scenario: "no further expectations",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: AddressOf(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(42)

				// ASSERT
				Error(t, err).IsNil()
				That(t, result).Equals(84)
				That(t, sut.expected).IsNil()
				That(t, sut.idxExpected).Equals(-1)
			},
		},
		{scenario: "arguments expected but not recorded",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: AddressOf(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall()

				// ASSERT
				That(t, result).Equals(84)
				Error(t, err).Is(ErrExpectedArgs)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnExpectationsWereMet(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "no errors",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Error(t, err).IsNil()
			},
		},
		{scenario: "expected calls/one error",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}
				sut.errs = append(sut.errs, errors.New("expected error"))

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Error(t, err).Is(ErrExpectationsNotMet)
				Error(t, err).Is(sut.errs[0])
			},
		},
		{scenario: "expected calls/multiple errors",
			exec: func(t *testing.T) {
				// ARRANGE
				err1 := errors.New("expected error 1")
				err2 := errors.New("expected error 2")
				sut := MockFn[int, int]{}
				sut.errs = append(sut.errs, err1)
				sut.errs = append(sut.errs, err2)

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Error(t, err).Is(ErrExpectationsNotMet)
				Error(t, err).Is(err1)
				Error(t, err).Is(err2)
			},
		},
		{scenario: "mapped results/unused results",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{42: {Result: 84}},
				}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Error(t, err).Is(ErrExpectationsNotMet)
				Error(t, err).Is(ErrResultNotUsed)
			},
		},
		{scenario: "mapped results/all used",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{42: {Result: 84}},
					actual:    []*mockFnCall[int, int]{{args: AddressOf(42), result: 84}},
				}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Error(t, err).IsNil()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnExpectCall(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "mapped results configured",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{},
				}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT + ASSERT
				sut.ExpectCall()
			},
		},
		{scenario: "valid configuration",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result := sut.ExpectCall()

				// ASSERT
				That(t, result).Equals(&mockFnCall[int, int]{})
				That(t, sut.expectations).Equals([]*mockFnCall[int, int]{result})
				That(t, sut.expected).Equals(result)
			},
		},
		{scenario: "multiple expectations",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result1 := sut.ExpectCall().WithArgs(1)
				result2 := sut.ExpectCall().WithArgs(2)

				// ASSERT
				That(t, result1).Equals(&mockFnCall[int, int]{args: AddressOf(1)})
				That(t, result2).Equals(&mockFnCall[int, int]{args: AddressOf(2)})
				That(t, sut.expectations).Equals([]*mockFnCall[int, int]{{args: AddressOf(1)}, {args: AddressOf(2)}})
				That(t, sut.expected).Equals(result1)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnReset(t *testing.T) {
	// ARRANGE
	sut := MockFn[int, int]{
		responses:    map[int]*Fake[int]{42: {Result: 84}},
		expectations: []*mockFnCall[int, int]{{args: AddressOf(42), result: 84}},
		expected:     &mockFnCall[int, int]{args: AddressOf(42), result: 84},
		actual:       []*mockFnCall[int, int]{{args: AddressOf(42), result: 84}},
		idxExpected:  1,
		errs:         []error{errors.New("expected error")},
	}

	// ACT
	sut.Reset()

	// ASSERT
	That(t, sut).Equals(MockFn[int, int]{})
}

func TestMockFnResultFor(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "expected calls are configured",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{{args: AddressOf(42), result: 84}},
				}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.ResultFor(42)
			},
		},
		{scenario: "no result configured for arguments",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{42: {Result: 84}},
				}
				defer ExpectPanic(ErrNoResultForArgs).Assert(t)

				// ACT
				sut.ResultFor(-1)
			},
		},
		{scenario: "valid configuration",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*Fake[int]{42: {Result: 84}},
				}

				// ACT
				result := sut.ResultFor(42)

				// ASSERT
				That(t, result).Equals(Fake[int]{Result: 84})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnWhenCalledWith(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "combined with expected calls",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := MockFn[int, int]{}
				sut.expectations = []*mockFnCall[int, int]{{args: AddressOf(42), result: 84}}

				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.WhenCalledWith(42)
			},
		},
		{scenario: "duplicate arguments",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := MockFn[int, int]{}
				sut.responses = map[int]*Fake[int]{42: {Result: 84}}

				defer ExpectPanic(ErrInvalidArgument).Assert(t)

				// ACT
				sut.WhenCalledWith(42)
			},
		},
		{scenario: "valid configuration",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				sut.WhenCalledWith(42).Returns(84)

				// ASSERT
				That(t, sut.responses).Equals(map[int]*Fake[int]{42: {Result: 84}})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnCallWillReturn(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "returns value",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := mockFnCall[int, int]{}

				// ACT
				sut.WillReturn(42)

				// ASSERT
				That(t, sut).Equals(mockFnCall[int, int]{result: 42})
			},
		},
		{scenario: "returns error",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")

				// ACT
				sut.WillReturn(err)

				// ASSERT
				That(t, sut).Equals(mockFnCall[int, int]{err: err})
			},
		},
		{scenario: "returns value and error",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")

				// ACT
				sut.WillReturn(42, err)

				// ASSERT
				That(t, sut).Equals(mockFnCall[int, int]{result: 42, err: err})
			},
		},
		{scenario: "multiple return values",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.WillReturn(42, 42)
			},
		},
		{scenario: "multiple errors",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.WillReturn(err, err)
			},
		},
		{scenario: "invalid type",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.WillReturn("invalid type")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestMockFnCallWithArgs(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "valid configuration",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := &mockFnCall[int, int]{}

				// ACT
				result := sut.WithArgs(42)

				// ASSERT
				That(t, sut).Equals(&mockFnCall[int, int]{args: AddressOf(42)})
				Value(t, result).Equals(sut)
			},
		},
		{scenario: "multiple arguments configured",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := &mockFnCall[int, int]{}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				_ = sut.WithArgs(42)
				_ = sut.WithArgs(42)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
