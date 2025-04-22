package test

import (
	"fmt"
	"testing"
)

func TestByRef(t *testing.T) {
	With(t)

	// ACT
	addr := ByRef("some literal")

	// ASSERT
	Expect(*addr).To(Equal("some literal"))
}

func ExampleByRef() {
	s := ByRef("value")
	fmt.Println(*s)

	// Output:
	// value
}
