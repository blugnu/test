package testframe

import (
	"testing"
)

type cleaner struct {
	cleanupCalled bool
}

func (c *cleaner) Cleanup(func()) {
	c.cleanupCalled = true
}

func TestNilFrame(t *testing.T) {
	t.Run("cleanup delegation", func(t *testing.T) {
		c := &cleaner{}
		sut := Nil{T: c}

		sut.Cleanup(nil)

		if !c.cleanupCalled {
			t.Error("expected Cleanup to be called, but it was not")
		}
	})

	t.Run("coverage", func(t *testing.T) {
		// testframe.Nil is required to implement test.TestingT in order to be
		// usable with test.With().  Only the Cleanup implementation is
		// significant; the remainder of the methods are no-ops and are tested
		// only for coverage.
		sut := Nil{T: &cleaner{}}

		if sut.Name() != "NilFrame" {
			t.Errorf("expected Name to return 'NilFrame', got '%s'", sut.Name())
		}

		if !sut.Run("test", nil) {
			t.Error("expected Run to return true")
		}

		sut.Error("test error")
		sut.Errorf("test error %s", "formatted")
		sut.Fail()
		sut.FailNow()
		if sut.Failed() {
			t.Error("expected Failed to return false")
		}

		sut.Fatal("test fatal error")
		sut.Fatalf("test fatal error %s", "formatted")
		sut.Helper()
		sut.Parallel()
		sut.Setenv("key", "value")
		sut.SkipNow()
	})
}
