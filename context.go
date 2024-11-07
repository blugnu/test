package test

import (
	"context"
	"testing"
)

// ContextIndicator[T, I] returns a test.AnyTest[I] (constructed using test.That[I]), to facilitate
// testing of an indicator value returned by a function to extract some value from a context.Context
// that returns a value and an indicator value.
//
// # params
//
//	t *testing.T                       // the testing.T instance
//	ctx context.Context                // the context to extract the value from
//	fn func(context.Context) (T, I)    // the function to extract the value and indicator value
//	opts ...any                        // optional test options (will be passed to the That[I] test
//	                                   // returned for testing the captured indicator)
//
// # returns
//
//	AnyTest[I]     // an AnyTest[I] for the indicator returned by the function
//
// # usage
//
// This function provides convenience for testing functions that return a value of type T and an
// indicator value of type I from a context.Context. The types T and I will be inferred from the
// signature of the supplied function; the caller need only provide the function and the context.
//
// Additional options may be specified which are passed to the That[I] test function that is returned
// for testing the value.
//
// The type T value is ignored and discarded by this function; only the indicator is captured
// for testing; to test the T value, use test.ContextValue().
//
// # example
//
// The following example demonstrates how to test the indicator returned by a context function:
//
//	 func ContextKey(ctx context.Context) (string, bool) {
//		 if v := context.Value(ctx, "key"); v != nil {
//			 return v.(string), true
//		 }
//		 return "", false
//	 }
//
//	 func TestContextKey(t *testing.T) {
//		 // ARRANGE
//		 ctx := context.WithValue(context.Background(), "key", "value")
//
//		 // ACT
//		 test.ContextIndicator(t, ctx, ContextKey).Equals(true)
//	 }
func ContextIndicator[T any, I any](
	t *testing.T,
	ctx context.Context,
	fn func(context.Context) (T, I),
	opts ...any,
) AnyTest[I] {
	t.Helper()
	_, i := fn(ctx)
	return That(t, i, opts...)
}

// ContextValue[T, I] returns a test.AnyTest[T] (constructed using test.That[T]), to facilitate
// testing of a value returned by a function to extract some value from a context.Context that
// also returns an indicator value, type I.
//
// # params
//
//	t *testing.T                       // the testing.T instance
//	ctx context.Context                // the context to extract the value from
//	fn func(context.Context) (T, I)    // the function to extract the value and indicator value
//	opts ...any                        // optional test options (will be passed to the That[T] test
//	                                   // returned for testing the captured value)
//
// # returns
//
//	AnyTest[T]     // an AnyTest[T] for the value returned by the function
//
// # usage
//
// This function provides convenience for testing functions that return a value of type T and an
// indicator value of type I from a context.Context. The types T and I will be inferred from the
// signature of the supplied function; the caller need only provide the function and the context.
//
// Additional options may be specified which are passed to the That test function that is returned
// for testing the value.
//
// The type I value is ignored and discarded by this function; only the T value is captured for
// testing. The indicator may be tested using test.ContextIndicator().
//
// # note
//
// This function is not useful for, and does not support testing of, context functions that return
// only a value with no indicator.  For such functions the returned value can be tested directly:
//
//	 func ContextKey(ctx context.Context) string {
//		 if v := context.Value(ctx, "key"); v != nil {
//			 return v.(string)
//		 }
//		 return ""
//	 }
//
//	 func TestContextKey(t *testing.T) {
//		 // ARRANGE + ACT
//		 ctx := context.WithValue(context.Background(), "key", "value")
//
//		 // ASSERT
//		 test.That(t, ContextKey(ctx)).Equals("value")
//	 }
//
// # example
//
// The following example demonstrates how to test the value returned by a context function that
// returns a value and an indicator:
//
//	 func ContextKey(ctx context.Context) (string, bool) {
//		 if v := context.Value(ctx, "key"); v != nil {
//			 return v.(string), true
//		 }
//		 return "", false
//	 }
//
//	 func TestContextKey(t *testing.T) {
//		 // ARRANGE
//		 ctx := context.WithValue(context.Background(), "key", "value")
//
//		 // ACT
//		 test.ContextValue(t, ctx, ContextKey).Equals("value")
//	 }
func ContextValue[T any, I any](t *testing.T, ctx context.Context, fn func(context.Context) (T, I), opts ...any) AnyTest[T] {
	t.Helper()
	v, _ := fn(ctx)
	return That(t, v, opts...)
}
