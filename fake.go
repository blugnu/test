package test

// Fake is a generic type that can be used to fake a function call returning
// some result type R and/or an error. It is useful for creating fakes for 
// interface implmentations by composing an anonymous Fake[R] into a struct
// implementing the desired method to be faked by returning the required values
// from the Fake.
//
// By using an anonymous field of type Fake[R] in a struct, the Reset() method
// of the Fake is promoted to the struct, allowing the struct to be reset to its
// zero value. 
//
// When faking an interface method that returns multiple result values (in
// addition to an error), use a struct type with fields for each of the result
// values.
//
// When faking an interface method that returns only result values and no error
// simply ignore the Err field.  To fake a method that returns only an error,
// it is recommended to specify a result type of any and ignore the Result field.
//
// # Example
//
// 	type MyInterface interface {
// 		MyMethod() (int, error)
// 	}
//
// 	type MyFake struct {
// 		Fake[int]
// 	}
//
// 	func (f *MyFake) MyMethod() (int, error) {
// 		return f.Result, f.Err
// 	}
type Fake[R any] struct {
	Result R
	Err error
}

// Reset sets the Result and Err fields to their zero values.
func (f *Fake[R]) Reset() {
	f.Result = *new(R)
	f.Err = nil
}