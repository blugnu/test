package test

import (
	"errors"
	"fmt"
)

// MockFn is a generic type that can be used to mock a function returning some result
// (type R) and/or an error, with a value (type A) to capture arguments.
//
// MockFn has two modes of operation:
//
//   - Expected Calls: calls are configured using a fluent configuration API starting
//     with an ExpectCall method, optionally configuring any expected arguments and
//     te result and/or error to be returned.  In this mode, the mock will fail to meet
//     expectations if arguments recorded with actual calls do not match the arguments
//     for the corresponding expected call.
//
//   - Mapped Results: results for a given set of arguments are configured using the
//     WhenCalledWith method.  In this mode, the mock will fail to meet expectations if
//     all configured argument:result combinations are not used.
//
// An implementation of the function being mocked must be provided by the test code to
// record calls to the mock function and to return the configured result and/or error
// using either ExpectedResult (expected calls) or ResultFor (mapped results) methods.
//
// The type can be used in a variety of ways to suit the requirements of the mock in a
// particular test scenario.
//
// # TL;DR: Simple Fakes
//
// When the following conditions apply, consider using the simpler test.Fake[R] type:
//
// - the number and order of calls to the mocked function is not significant
// - the arguments used to call the mocked function are not significant
// - the return value (and/or error) of the mocked function is consistent across all calls
//
// # Mocking Return Values (Fake)
//
// The R type parameter is used to specify the return type of the function being mocked.
// The WillReturn method can be used to specify the return value that the mock function
// should return for an expected call.
//
// When mocking a function that returns only an error, specify a result type of any and
// ignore the Result field.
//
// When mocking a function that returns multiple result values (in addition to an error),
// use a struct type to specify the return type, with fields for each of the result values.
//
// # Testing and Recording Arguments (Spy)
//
// MockFn does not implement the function being mocked; the test code must provide the
// implementation, which should record each call to the MockFn using RecordCall(), specifying
// any arguments.
//
// If the method being mocked accepts multiple multiple arguments, they may be captured
// using a struct type for the A type parameter, with fields for each of the arguments
// to be captured.
//
// Similarly, when mocking an interface method that returns multiple result values
// (in addition to an error), a struct type may be used for the R type parameter,
// with fields for each of the result values.
//
// When mocking a method which accepts no arguments, or if not interested in the
// arguments, then you may use type any for the A type parameter, or consider
// using the simpler test.Fake[R] type instead.
//
// Similarly, when mocking function that returns only result values and no error
// simply ignore the Err field.  To fake a method that returns only an error, specify
// a result type of any and ignore the Result field.
//
// # Example
//
//	type MyInterface interface {
//		MyMethod() (int, error)
//		MyMethodWithArgs(id string, opt bool) (int, error)
//	}
//
//	type myMock struct {
//		myMethod         test.Fake[int]
//		myMethodWithArgs test.MockFn[struct{ID string; Opt bool}, int]
//	}
//
//	func (mock *myMock) MyMethod() (int, error) {
//		return mock.myMethod.Result, mock.myMethod.Err
//	}
//
//	func (mock *myMock) MyMethodWithArgs(id string, opt bool) (int, error) {
//		return mock.myMethodWithArgs.CalledWith(struct{ID string; Opt bool}{ID: id, Opt: opt})
//	}
type MockFn[A comparable, R any] struct {
	// actual is a slice of all calls made to the mock function by the code under test
	actual []*mockFnCall[A, R]

	// expectations is a slice of all expected calls to the mock function; the slice is nil
	// if the mock function is not configured for expected calls.
	expectations []*mockFnCall[A, R]

	// expected is the next expected call to the mock function; the value is nil if there are no
	// more expected calls
	expected *mockFnCall[A, R]

	// idxExpected is the index of the next expected call in the expectations slice; the value
	// is 0 if no expected calls have yet been recorded and -1 if all expected calls have been
	// recorded
	idxExpected int

	// responses is a map of arguments to the Fake[R] values configured for those arguments;
	// the map is nil if the mock function is not configured for mapped results
	responses map[A]*Fake[R]

	// errs is a slice of errors recorded during the test; the slice is nil if no errors have
	// been recorded
	errs []error
}

// mockFnCall represents a call to a mock function.  It is used both to configure expected
// calls and to record actual calls.
type mockFnCall[A comparable, R any] struct {
	// args is a pointer to the arguments associated with the call.
	//
	// - For an expected call: the value is nil if the call does not have any arguments or
	//   if the arguments are not significant;
	//
	// - For an actual call: the value is nil if there were no arguments or arguments were
	//   not recorded.
	args *A

	// result is the result associated with the call.
	//
	// - For an expected call: the value(s) to be returned as the result when the expected call
	// is made;
	//
	// - For a recorded calls: the value is not used and is the zero value of the R type.
	result R

	// err is the error associated with the call.
	//
	// - For an expected call: the error to be returned when this expected call is made;
	//
	// - For a recorded call: any error determined when the call is recorded.
	err error
}

// WillReturn configures the result and/or error to be returned by the mock function for an
// expected call.
//
// The method accepts variadic arguments of any type; the first value of type R is used as
// the result to be returned by the mock function, and the first value of type error is used
// as the error to be returned by the mock function.
//
// If multiple values of type R or error are provided the function will panic. This ensures
// that the mock function is configured correctly and that a test that is not setup correctly
// fails immediately without yielding potentially incorrect results.
func (mock *mockFnCall[A, R]) WillReturn(v ...any) {
	resultSet := false
	errSet := false
	for _, r := range v {
		switch r := r.(type) {
		case R:
			if resultSet {
				panic(fmt.Errorf("%w: only one result value may be configured", ErrInvalidOperation))
			}
			mock.result = r
			resultSet = true
		case error:
			if errSet {
				panic(fmt.Errorf("%w: only one error value may be configured", ErrInvalidOperation))
			}
			mock.err = r
			errSet = true
		default:
			panic(fmt.Errorf("%w: %T: only values of type %T or error may be specified", ErrInvalidOperation, r, *new(R)))
		}
	}
}

// WithArgs configures the arguments associated with an expected call to the mock function.
//
// The method accepts a single argument of type A, which is used to configure the arguments
// for the expected call.  If args are already configured for the mock function, the function
// will panic with an ErrInvalidOperation error.
func (mock *mockFnCall[A, R]) WithArgs(args A) *mockFnCall[A, R] {
	if mock.args != nil {
		panic(fmt.Errorf("%w: arguments already configured", ErrInvalidOperation))
	}

	mock.args = &args
	return mock
}

// RecordCall is used by a mock implementation to record a call to a mock function,
// optionally testing that arguments match those expected and returning the result and
// error configured for the expected call.
//
// # returns
//
// Returns mock.expected.Result for an expected call; if there is no expected call
// the zero value of the R type parameter is returned with ErrUnexpectedCall.
//
// If the arguments do not match the expected call, mock.expected.Result is
// returned with ErrUnexpectedArgs.
//
// # errors
//
//	ErrUnexpectedCall       the call has no corresponding expected call
//	ErrUnexpectedArgs       the arguments do not match the expected call
//	mock.expected.Err	    the error specified for the expected call (if any)
//
// # example
//
//	type MyInterface interface {
//		MyMethod(int) (int, error)
//	}
//
//	type myMock struct {
//		myMethod test.MockFn[int, int]
//	}
//
//	func (mock *myMock) MyMethod(arg int) (int, error) {
//		return mock.myMethod.RecordCall(arg)
//	}
//
// # Multiple Return Values
//
// When a mocked function returns multiple values (in addition to an error), the
// return value will typically be a struct with fields for each of the result values.
//
// The struct fields must be returned as individual return values by the mock
// implementation:
//
//	func (mock *myMock) IntDiv(num, div int) (struct{Result, Remainder int}, error) {
//		result, err := mock.intDiv.RecordCall(any)
//		return result.Result, result.Remainder, err
//	}
func (mock *MockFn[A, R]) RecordCall(args ...A) (returns R, err error) {
	if mock.responses != nil {
		panic(fmt.Errorf("%w: mock function is configured for mapped results; use <fn>.ResultFor()", ErrInvalidOperation))
	}

	actual := &mockFnCall[A, R]{}
	if len(args) > 0 {
		actual.args = &args[0]
	}

	if mock.expected == nil {
		actual.err = ErrUnexpectedCall
		err = fmt.Errorf("%w: with args: %v", ErrUnexpectedCall, args)
		mock.actual = append(mock.actual, actual)
		mock.errs = append(mock.errs, err)
		return
	}

	err = mock.expected.err
	returns = mock.expected.result

	switch {
	case mock.expected.args == nil:
		// NO-OP - no arguments of interest were configured for the expected call

	case actual.args == nil:
		actual.err = ErrExpectedArgs
		err = fmt.Errorf("%w:\n  expected: %v\n  got     : nil (no args recorded)", ErrExpectedArgs, mock.expected.args)
		mock.errs = append(mock.errs, err)

	case *mock.expected.args != *actual.args:
		actual.err = ErrUnexpectedArgs
		err = fmt.Errorf("%w:\n  expected: %v\n  got     : %v", ErrUnexpectedArgs, mock.expected.args, actual.args)
		mock.errs = append(mock.errs, err)
	}

	mock.actual = append(mock.actual, actual)
	mock.idxExpected++
	mock.expected = nil

	if mock.idxExpected < len(mock.expectations) {
		mock.expected = mock.expectations[mock.idxExpected]
	} else {
		mock.idxExpected = -1
	}

	return
}

// ExpectedResults returns an error if any expectations were not met; otherwise nil.
//
// This method is typically called at the end of a test to ensure that all expected
// calls were made.
//
// If the mock function is configured for mapped results, this method will return an
// error if any results were not used.
//
// # errors
//
//	ErrExpectationsNotMet      // one or more expectations were not met; the error is
//	                           // joined with errors for each unmet expectation
func (mock *MockFn[A, R]) ExpectationsWereMet() error {
	if mock.responses != nil {
		mock.errs = nil
		args := map[A]struct{}{}
		for arg := range mock.responses {
			args[arg] = struct{}{}
		}
		for _, called := range mock.actual {
			delete(args, *called.args)
		}
		if len(args) > 0 {
			for arg := range args {
				mock.errs = append(mock.errs, fmt.Errorf("%w: %v", ErrResultNotUsed, arg))
			}
			return fmt.Errorf("%w: %w", ErrExpectationsNotMet, errors.Join(mock.errs...))
		}
	}
	if len(mock.errs) > 0 {
		return fmt.Errorf("%w: %w", ErrExpectationsNotMet, errors.Join(mock.errs...))
	}
	return nil
}

func (mock *MockFn[A, R]) ExpectCall() *mockFnCall[A, R] {
	if mock.responses != nil {
		panic(fmt.Errorf("%w: cannot combine expected calls with mapped results", ErrInvalidOperation))
	}

	ex := &mockFnCall[A, R]{}
	mock.expectations = append(mock.expectations, ex)
	if mock.expected == nil {
		mock.expected = ex
	}
	return ex
}

// Reset sets the mock function to its zero value (no errors, no expected or recorded calls
// and no mapped results).
func (mock *MockFn[A, R]) Reset() {
	*mock = MockFn[A, R]{}
}

// ResultFor returns the Fake[R] value configured for the specified arguments. This method
// is called by a mock implementation to return the result and/or error configured for a
// specific set of arguments.
//
// # errors
//
// In the event of an error, the function will panic with one of the following errors:
//
//	ErrInvalidOperation     // when the mock function is configured for expected calls
//
//	ErrNoResultForArgs      // when no result is configured for the specified arguments
func (mock *MockFn[A, R]) ResultFor(args A) Fake[R] {
	if mock.responses == nil {
		panic(fmt.Errorf("%w: mock function is configured for expected calls; use <fn>.ExpectedResult()", ErrInvalidOperation))
	}

	if result, ok := mock.responses[args]; ok {
		return *result
	}
	panic(ErrNoResultForArgs)
}

// WhenCalledWith is used to configure the result for a specific set of arguments. This
// is useful when configuring a test where the result of a mocked function call depends
// on the arguments passed to the function but the arguments themselves, the number of
// times the function is called, or the order in which calls to the function are made are
// not significant.
//
// Only one results can be configured for a given set of arguments.
//
// If the arguments, the number and/or order of calls are significant or if different
// results are required to be mocked for different calls with the same arguments, use
// ExpectCall to configure expected calls instead.
//
// Mocked results mapped to arguments may not be combined with expected calls.
//
// # returns
//
// Returns a Fake[R] value that can be used to configure the result and/or error for the
// specified arguments.
//
// # errors
//
// In the event of an error, the function will panic with one of the following errors:
//
//	ErrInvalidOperation     // when the mock function already has one or more expected calls
//	                        // configured
//
//	ErrInvalidArgument      // when a result is already configured for the specified arguments
func (mock *MockFn[A, R]) WhenCalledWith(args A) *Fake[R] {
	if len(mock.expectations) > 0 {
		panic(fmt.Errorf("%w: cannot combine mapped results with expected calls", ErrInvalidOperation))
	}

	switch {
	case mock.responses == nil:
		mock.responses = make(map[A]*Fake[R])
	default:
		if _, ok := mock.responses[args]; ok {
			panic(fmt.Errorf("%w: result already configured for args: %v", ErrInvalidArgument, args))
		}
	}

	r := &Fake[R]{}
	mock.responses[args] = r

	return r
}
