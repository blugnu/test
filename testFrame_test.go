package test

import (
	"errors"
	"testing"
)

func Test_goid(t *testing.T) {
	With(t)

	Run("returns the id of the calling goroutine", func() {
		id := goid()
		Expect(id).ToNot(Equal(uintptr(0)))
	})

	Run("returns different ids for different goroutines", func() {
		var id1, id2 uintptr
		await(func() { id1 = goid() })
		await(func() { id2 = goid() })
		Expect(id1).ToNot(Equal(id2))
	})

	Run("when unable to get stack", func() {
		// NOTE: it is important to setup the panic handler BEFORE replacing
		// the stack function to avoid interfering with the identification of
		// T() for the panic handler test itself.
		defer Expect(Panic(ErrNoStack)).DidOccur()

		defer AfterUsing(&stack, func([]byte, bool) int {
			return 0 // simulate failure to get stack
		})()

		_ = goid()
	})

	Run("when unable to parse stack", func() {
		// NOTE: it is important to setup the panic handler BEFORE replacing
		// the stack function to avoid interfering with the identification of
		// T() for the panic handler test itself.
		defer Expect(Panic(ErrUnexpectedStackFormat)).DidOccur()
		t := T()
		t.Cleanup(AfterUsing(&stack, func([]byte, bool) int {
			return 1 // simulate failure to parse stack
		}))

		_ = goid()
	})
}

func ExpectErrNoTestFrame(fn func()) {
	// we deliberately nil the test frame to simulate a test frame not being
	// set. This also means that we cannot use the test frame to check for
	// errors, since it is nil so we cache the test frame in a local variable
	// before nilling it.
	t := T()

	// make sure we restore the test frame when done
	defer With(t)
	defer func() {
		t.Helper()

		if r := recover(); r != nil {
			if err, ok := r.(error); !ok || !errors.Is(err, ErrNoTestFrame) {
				t.Errorf("expected ErrNoTestFrame, got %v", r)
			}
			return // we got the expected panic
		}
		t.Errorf("expected panic, got nil")
	}()

	t.Helper()
	With(nil)
	fn()
}
