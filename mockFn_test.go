package test //nolint:testpackage // tests private functions and types

import (
	"errors"
	"testing"
)

func byref[T any](v T) *T {
	return &v
}

func TestMockFnRecordCall(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "mapped results are configured",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{},
				}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				_, _ = sut.RecordCall(42)
			},
		},
		{scenario: "unexpected call",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result, err := sut.RecordCall()

				// ASSERT
				Expect(result).To(Equal(0))
				Expect(err).Is(ErrUnexpectedCall)
			},
		},
		{scenario: "unexpected arguments",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: byref(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(0)

				// ASSERT
				Expect(result).To(Equal(84))
				Expect(err).Is(ErrUnexpectedArgs)
			},
		},
		{scenario: "further expectations",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: byref(42), result: 84},
						{args: byref(43), result: 86},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(42)

				// ASSERT
				Expect(err).IsNil()
				Expect(result).To(Equal(84))
				Expect(sut.expected).To(Equal(sut.expectations[1]))
				Expect(sut.idxExpected).To(Equal(1))
			},
		},
		{scenario: "no further expectations",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: byref(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(42)

				// ASSERT
				Expect(err).IsNil()
				Expect(result).To(Equal(84))
				Expect(sut.expected).IsNil()
				Expect(sut.idxExpected).To(Equal(-1))
			},
		},
		{scenario: "arguments expected but not recorded",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: byref(42), result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall()

				// ASSERT
				Expect(result).To(Equal(84))
				Expect(err).Is(ErrExpectedArgs)
			},
		},
		{scenario: "call expected regardless of arguments",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{
						{args: nil, result: 84},
					},
				}
				sut.expected = sut.expectations[0]

				// ACT
				result, err := sut.RecordCall(42)

				// ASSERT
				Expect(err).IsNil()
				Expect(result).To(Equal(84))
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnExpectationsWereMet(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "no errors",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Expect(err).IsNil()
			},
		},
		{scenario: "expected calls/one error",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}
				sut.errs = append(sut.errs, errors.New("expected error"))

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Expect(err).Is(ErrExpectationsNotMet)
				Expect(err).Is(sut.errs[0])
			},
		},
		{scenario: "expected calls/multiple errors",
			exec: func() {
				// ARRANGE
				err1 := errors.New("expected error 1")
				err2 := errors.New("expected error 2")
				sut := MockFn[int, int]{}
				sut.errs = append(sut.errs, err1)
				sut.errs = append(sut.errs, err2)

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Expect(err).Is(ErrExpectationsNotMet)
				Expect(err).Is(err1)
				Expect(err).Is(err2)
			},
		},
		{scenario: "mapped results/unused results",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{42: {Result: 84}},
				}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Expect(err).Is(ErrExpectationsNotMet)
				Expect(err).Is(ErrResultNotUsed)
			},
		},
		{scenario: "mapped results/all used",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{42: {Result: 84}},
					actual:    []*mockFnCall[int, int]{{args: byref(42), result: 84}},
				}

				// ACT
				err := sut.ExpectationsWereMet()

				// ASSERT
				Expect(err).IsNil()
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnExpectCall(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "mapped results configured",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{},
				}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT + ASSERT
				sut.ExpectCall()
			},
		},
		{scenario: "valid configuration",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result := sut.ExpectCall()

				// ASSERT
				Expect(result).To(DeepEqual(&mockFnCall[int, int]{}))
				Expect(sut.expectations).To(DeepEqual([]*mockFnCall[int, int]{result}))
				Expect(sut.expected).To(Equal(result))
			},
		},
		{scenario: "multiple expectations",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				result1 := sut.ExpectCall().WithArgs(1)
				result2 := sut.ExpectCall().WithArgs(2)

				// ASSERT
				Expect(result1).To(DeepEqual(&mockFnCall[int, int]{args: byref(1)}))
				Expect(result2).To(DeepEqual(&mockFnCall[int, int]{args: byref(2)}))
				Expect(sut.expectations).To(DeepEqual([]*mockFnCall[int, int]{{args: byref(1)}, {args: byref(2)}}))
				Expect(sut.expected).To(Equal(result1))
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnReset(t *testing.T) {
	With(t)

	// ARRANGE
	sut := MockFn[int, int]{
		responses:    map[int]*FakeResult[int]{42: {Result: 84}},
		expectations: []*mockFnCall[int, int]{{args: byref(42), result: 84}},
		expected:     &mockFnCall[int, int]{args: byref(42), result: 84},
		actual:       []*mockFnCall[int, int]{{args: byref(42), result: 84}},
		idxExpected:  1,
		errs:         []error{errors.New("expected error")},
	}

	// ACT
	sut.Reset()

	// ASSERT
	Expect(sut).To(DeepEqual(MockFn[int, int]{}))
}

func TestMockFnResultFor(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "expected calls are configured",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					expectations: []*mockFnCall[int, int]{{args: byref(42), result: 84}},
				}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				sut.ResultFor(42)
			},
		},
		{scenario: "no result configured for arguments",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{42: {Result: 84}},
				}
				defer Expect(Panic(ErrNoResultForArgs)).DidOccur()

				// ACT
				sut.ResultFor(-1)
			},
		},
		{scenario: "valid configuration",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{
					responses: map[int]*FakeResult[int]{42: {Result: 84}},
				}

				// ACT
				result := sut.ResultFor(42)

				// ASSERT
				Expect(result).To(Equal(FakeResult[int]{Result: 84}))
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnWhenCalledWith(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "combined with expected calls",
			exec: func() {
				// ARRANGE + ASSERT
				sut := MockFn[int, int]{}
				sut.expectations = []*mockFnCall[int, int]{{args: byref(42), result: 84}}

				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				sut.WhenCalledWith(42)
			},
		},
		{scenario: "duplicate arguments",
			exec: func() {
				// ARRANGE + ASSERT
				sut := MockFn[int, int]{}
				sut.responses = map[int]*FakeResult[int]{42: {Result: 84}}

				defer Expect(Panic(ErrInvalidArgument)).DidOccur()

				// ACT
				sut.WhenCalledWith(42)
			},
		},
		{scenario: "valid configuration",
			exec: func() {
				// ARRANGE
				sut := MockFn[int, int]{}

				// ACT
				sut.WhenCalledWith(42).Returns(84)

				// ASSERT
				Expect(sut.responses).To(EqualMap(map[int]*FakeResult[int]{42: {Result: 84}}))
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnCallWillReturn(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "returns value",
			exec: func() {
				// ARRANGE
				sut := mockFnCall[int, int]{}

				// ACT
				sut.WillReturn(42)

				// ASSERT
				Expect(sut).To(Equal(mockFnCall[int, int]{result: 42}))
			},
		},
		{scenario: "returns error",
			exec: func() {
				// ARRANGE
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")

				// ACT
				sut.WillReturn(err)

				// ASSERT
				Expect(sut).To(Equal(mockFnCall[int, int]{err: err}))
			},
		},
		{scenario: "returns value and error",
			exec: func() {
				// ARRANGE
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")

				// ACT
				sut.WillReturn(42, err)

				// ASSERT
				Expect(sut).To(Equal(mockFnCall[int, int]{result: 42, err: err}))
			},
		},
		{scenario: "multiple return values",
			exec: func() {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				sut.WillReturn(42, 42)
			},
		},
		{scenario: "multiple errors",
			exec: func() {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				err := errors.New("expected error")
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				sut.WillReturn(err, err)
			},
		},
		{scenario: "invalid type",
			exec: func() {
				// ARRANGE + ASSERT
				sut := mockFnCall[int, int]{}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				sut.WillReturn("invalid type")
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}

func TestMockFnCallWithArgs(t *testing.T) {
	With(t)

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func()
	}{
		{scenario: "valid configuration",
			exec: func() {
				// ARRANGE
				sut := &mockFnCall[int, int]{}

				// ACT
				result := sut.WithArgs(42)

				// ASSERT
				Expect(sut).To(DeepEqual(&mockFnCall[int, int]{args: byref(42)}))
				Expect(result).To(Equal(sut))
			},
		},
		{scenario: "multiple arguments configured",
			exec: func() {
				// ARRANGE + ASSERT
				sut := &mockFnCall[int, int]{}
				defer Expect(Panic(ErrInvalidOperation)).DidOccur()

				// ACT
				_ = sut.WithArgs(42)
				_ = sut.WithArgs(42)
			},
		},
	}
	for _, tc := range testcases {
		Run(tc.scenario, func() {
			tc.exec()
		})
	}
}
