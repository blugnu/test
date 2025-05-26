package test

import "fmt"

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
// FakeResult[R} is intended for simple cases where a faked function or method is
// returning a consistent result value and/or error for all calls.
//
// For cases requiring such capabilities, consider using test.MockFn[A, R].
//
// # Example: Mocking an interface with a single method
//
//	type MyInterface interface {
//		MyMethod() (int, error)
//	}
//
//	type mockMyMethodInterface struct {
//		FakeResult[int]
//	}
//
//	func (mock *MyFake) MyMethod() (int, error) {
//		return mock.Result, mock.Err
//	}
//
//	// FakeResult[R] implements Reset() so by using an embedded field this is promoted
//	// to the mockMyMethodInterface struct itself.
//
// If mocking an interface with multiple methods, the mockMyMethodInterface
// struct would implement a Reset() method to reset each of the FakeResult[R] fields
// individually.
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
//	type myMethodResult struct {
//		myMethod FakeResult[struct{name string; age int}]
//	}
//
//	func (mock *mockMyInterface) MyMethod() (int, string, error) {
//		return mock.myMethod.Result.name, mock.myMethod.Result.age, mock.myMethod.Err
//	}
//
// # Return Value(s) with No Error
//
// When faking a function that returns only result values and no error, simply ignore
// the Err field.
//
// # Returning Only an Error
//
// To fake a method that returns only an error, specify a result type of any and
// ignore the Result field.
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
// value, and the first error value is used to set the error.
//
// If multiple R or error values are provided, the function will panic.  This
// ensures that a test that is not setup correctly fails immediately without
// yielding potentially incorrect results.
func (fake *FakeResult[R]) Returns(v ...any) {
	resultSet := false
	errSet := false
	for _, r := range v {
		switch r := r.(type) {
		case R:
			if resultSet {
				panic(fmt.Errorf("%w: only one result value (R) may be specified", ErrInvalidOperation))
			}
			resultSet = true
			fake.Result = r
		case error:
			if errSet {
				panic(fmt.Errorf("%w: only one error value may be specified", ErrInvalidOperation))
			}
			errSet = true
			fake.Err = r
		default:
			panic(fmt.Errorf("%w: %T: only values of type %T or error may be specified", ErrInvalidOperation, r, *new(R)))
		}
	}
}
