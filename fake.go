package test

import (
	"fmt"

	"github.com/blugnu/test/test"
)

// FakeResult[R] is a generic type that can be used to help fake a function returning
// some result type R and/or an error. It can be useful for creating simple fakes
// for functions or interface methods.
//
// The type does not provide any mechanism for providing an implementation of a
// function, but facilitates a consistent pattern for simple fakes and minimises
// the amount of boilerplate code required to create them.
//
// # Limitations
//
// No mechanism is provided for configuring expected calls, capturing or testing
// for expected arguments, or returning different values for multiple calls.
// FakeResult[R] is for simple cases where a fake function or method consistently
// returns a specific result.
//
// For cases requiring more advanced capabilities, consider using test.MockFn[A, R].
//
// # Example: Mocking an interface with a single method
//
//	type MyInterface interface {
//		MyMethod() (int, error)
//	}
//
//	type fakeMyMethodInterface struct {
//		FakeResult[int]
//	}
//
//	func (fake *MyFake) MyMethod() (int, error) {
//		return fake.Result, fake.Err
//	}
//
//	// FakeResult[R] implements Reset(); embedding promotes the Reset method
//	// to the fakeMyMethodInterface struct itself.
//
// If mocking an interface with multiple methods, the fakeMyMethodInterface
// struct would implement an explicit Reset() method to reset all FakeResult[R]
// fields individually.
//
// # Returning Multiple Values
//
// When faking an interface method that returns multiple result values (in
// addition to an error), use a struct type with fields for each of the result
// values.
//
//	type MyInterface interface {
//		MyMethod() (int, string, error)
//	}
//
//	type fakeMyMethod struct {
//		fakeMyMethodFn FakeResult[struct{name string; age int}]
//	}
//
//	func (fake *fakeMyMethod) MyMethod() (int, string, error) {
//		fn := fake.fakeMyMethodFn
//		return fn.Result.age, fn.Result.name, fn.Err
//	}
//
// # Return Value(s) with No Error
//
// When faking a function that returns only result values and no error, simply ignore
// the Err field.
//
// # Returning Only an Error
//
// When faking a function that only returns an error, use FakeResult[error].  This
// provides a FakeResult with a Result type of error in addition to the Err field; the
// Result field should be ignored.
//
// The use of the error type as the Result type clarifies the intent of the fake and
// avoids confusiion. Consider:
//
//	func (s MyStruct) SomeMethod() error {
//		return s.SomeMethodFn.Err
//	}
//
// It is clear at this point that the SomeMethodFn field is a FakeResult where only the
// error is relevant, and the Result field is ignored.  This is a common pattern when
// faking methods that return only an error.
//
// Now consider the possibilities when declaring the SomeMethodFn field:
//
// s.SomeMethodFn := FakeResult[any]{}     // does this fake a function returning any or is the Result ignored?
// s.SomeMethodFn := FakeResult[error]{}   // this fake is clearly for a function returning (only) an error
//
// The SomeMethodFn field is a FakeResult; since the Result field is ignored, the type
// parameter R can be any type, but using FakeResult[error] makes it clear that the
// fake is for a function that returns an error, even though the Result field, to which the
// type parameter 'error' relates, is ignored.
type FakeResult[R any] struct {
	Result R
	Err    error
}

// Reset resets the fake to its zero value.
func (fake *FakeResult[R]) Reset() {
	*fake = FakeResult[R]{}
}

// Returns sets the result value and/or error to be returned by the fake.
//
// The first R value in the variadic parameter list is used to set the result
// value, and the first error value is used to set the error.  nil values are
// ignored.  A value of any other type will cause the current test to fail as
// invalid.
//
// R and error values may be specified in any order, but for symmetry it is
// recommended that they are specified in the order they will be returned.
//
// e.g. when faking a function:
//
//	func MyFunc() (int, error)
//
//	fakeMyFuncFn := &FakeResult[int]{}
//
//	fakeMyFuncFn.Returns(42, nil) // fakes a function returning (42, nil)
//
// If multiple R or error values are provided, the function will fail the
// current test as invalid.
func (fake *FakeResult[R]) Returns(v ...any) {
	test.T().Helper()

	resultSet := false
	errSet := false
	for _, r := range v {
		switch r := r.(type) {
		case nil:
			// nil is ignored; this allows for faking functions that return
			// nilable results and errors in a satisfying and symmetrical way,
			// e.g. fake.Returns(value, nil)

		case error:
			if errSet {
				test.Error(ErrInvalidOperation, "only one error value may be specified")
			}
			errSet = true
			fake.Err = r

		case R:
			if resultSet {
				test.Error(ErrInvalidOperation, fmt.Sprintf("only one result value (%T) may be specified", *new(R)))
			}
			resultSet = true
			fake.Result = r

		default:
			test.Error(ErrInvalidOperation, fmt.Sprintf("only values of type %T or error (or nil) may be specified", *new(R)))
		}
	}
}
