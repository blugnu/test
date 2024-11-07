package test

import "fmt"

// Fake[R] is a generic type that can be used to fake a function returning some
// result type R and/or an error. It is useful for creating simple fakes for
// functions or interface methods.
//
// Fake[R] provides no mechanism for configuring expected calls, capturing or
// testing for expected arguments, or returning different values for multiple
// calls. It is intended for simple cases where a function or method is faked,
// returning a fixed result value and/or error.
//
// For more complex cases, consider using test.MockFn[A, R] instead.
//
// # example
//
//	type MyInterface interface {
//		MyMethod() (int, error)
//	}
//
//	type mockMyMethodInterface struct {
//		Fake[int]
//	}
//
//	func (mock *MyFake) MyMethod() (int, error) {
//		return mock.Result, mock.Err
//	}
//
//	// there is no need to implement (mock *MyFake) Reset() as the Reset() method
//	// of the embedded Fake[int] will be promoted to the mockMyMethodInterface struct.
//
// # Resetting a Fake
//
// A mock implementation should implement a Reset() method to reset any all
// all Fake[R] (and any MockFn[A, R]) fields.
//
// When mocking an interface with a single method, an anonymous field may be
// used in a struct implementing the interface, allowing the Reset() method of
// the Fake[R] to be promoted to the struct itself.
//
// # Returning Multiple Values
//
// When faking an interface method that returns multiple result values (in
// addition to an error), use a struct type with fields for each of the result
// values.
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
type Fake[R any] struct {
	Result R
	Err    error
}

// Reset resets the Fake to its zero value.
func (fake *Fake[R]) Reset() {
	*fake = Fake[R]{}
}

// Returns sets the result value and/or error to be returned by the fake.
//
// The first R value in the variadic parameter list is used to set the result
// value, and the first error value is used to set the error.
//
// If multiple R or error values are provided, the function will panic.  This
// ensures that a test that is not setup correctly fails immediately without
// yielding potentially incorrect results.
func (fake *Fake[R]) Returns(v ...any) {
	resultSet := false
	errSet := false
	for _, r := range v {
		switch r := r.(type) {
		case R:
			if resultSet {
				panic(fmt.Errorf("%w: multiple result values provided", ErrInvalidOperation))
			}
			resultSet = true
			fake.Result = r
		case error:
			if errSet {
				panic(fmt.Errorf("%w: multiple error values provided", ErrInvalidOperation))
			}
			errSet = true
			fake.Err = r
		default:
			panic(fmt.Errorf("%w: %T: only values of type %T or error may be specified", ErrInvalidOperation, r, *new(R)))
		}
	}
}
