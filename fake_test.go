package test

import "testing"

func TestFake(t *testing.T) {
	// ARRANGE
	fake := Fake[int]{Result: 42}

	// ACT
	fake.Reset()

	// ASSERT
	That(t, fake.Result).Equals(0)
}
