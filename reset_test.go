package test

import "testing"

type resettableInt int

func (r *resettableInt) Reset() {
	*r = 0
}

type resettableString string

func (r *resettableString) Reset() {
	*r = ""
}

func TestReset(t *testing.T) {
	With(t)

	// ARRANGE
	var i resettableInt = 42
	var s resettableString = "hello"

	// ACT
	Reset(&i, &s)

	// ASSERT
	Expect(i).To(Equal[resettableInt](0))
	Expect(s).To(Equal[resettableString](""))
}
