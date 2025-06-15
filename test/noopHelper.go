package test

// a noopHelper is a no-op implementation of the Helper interface.
//
// It is used to provide a default implementation of the Helper interface
// when no test frame is available, allowing the T() function to return
// a valid Helper interface without panicking.
//
// This enables helper functions wishing to report an invalid test or test
// error to do so without needing to check if a valid test frame is available:
//
//	func MyHelperFunction() {
//	   // ... evaluate pre-conditions for helper ...
//	   if !preCondition {
//	      test.T().Helper()                         // will not panic if no test frame is available
//	      test.Invalid("pre-conditions not met")    // will panic to report the invalid test
//	   }
//	   // ... continue with helper logic ...
//	}
type noopHelper struct{}

func (noopHelper) Helper() { /* NO-OP */ }

// coverage is a function that exists solely to ensure that the noopHelper is
// not optimized away by the compiler which impacts test coverage metrics
func coverage() {
	noopHelper{}.Helper()
}
