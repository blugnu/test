package test

import (
	"testing"
)

func TestAddressOf(t *testing.T) {
	// ACT
	addr := AddressOf("some literal")

	// ASSERT
	That(t, *addr).Equals("some literal")
}
