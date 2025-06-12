package test

import "github.com/blugnu/test/test"

type restorable[T any] struct {
	ptr *T
}

// ReplacedBy is a method on an Original value that allows you to replace
// the original variable with a new value for the duration of a test.
//
// It returns a function that, when called, will restore the original value
// of the variable; a call to this function should be deferred to ensure that
// the original value is restored even if the test panics or fails.
func (tf restorable[T]) ReplacedBy(r T) func(restorable[T]) {
	og := *tf.ptr
	*tf.ptr = r

	// returns a closure over the original variable address and value
	// that will restore the original value when called
	//
	// the (unused) parameter ensures that the returned function is
	// compatible with the Restore function signature
	return func(restorable[T]) { *tf.ptr = og }
}

// Restore is used to restore the original value of a variable
// after it has been temporarily changed for a test using
// the Original function and its ReplacedBy method.
//
// It should be called in a defer statement to ensure that the
// original value is restored even if the test panics or fails.
func Restore[T any](fn func(restorable[T])) {
	// This function is used to restore the original value of a variable
	// after it has been temporarily changed for a test.
	// It is intended to be called in a defer statement.
	fn(restorable[T]{})
}

// Original is used to create an override for a variable of type T.
//
// The argument v is a pointer to the variable that will be overridden.
// The function returns a value providing a ReplacedBy method,
// which can be used to replace the original value with a new one.
func Original[T any](v *T) restorable[T] {
	if v == nil {
		test.T().Helper()
		test.Invalid("test.Original: cannot create an override for a nil pointer")
	}
	return restorable[T]{ptr: v}
}
