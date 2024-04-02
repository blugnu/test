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
	// ARRANGE
	var i resettableInt = 42
	var s resettableString = "hello"

	// ACT
	Reset(&i, &s)

	// ASSERT
	That(t, i).Equals(0)
	That(t, s).Equals("")
}
