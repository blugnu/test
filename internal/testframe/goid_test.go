package testframe

import (
	"errors"
	"sync"
	"testing"
)

func await(f func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	wg.Wait()
}

// these tests must use standard library testing.T to avoid an import cycle :(
func Test_goid(t *testing.T) {
	t.Run("returns the id of the calling goroutine", func(t *testing.T) {
		id := goid()
		if id == uintptr(0) {
			t.Errorf("expected goroutine id to be non-zero, got %d", id)
		}
	})

	t.Run("returns different ids for different goroutines", func(t *testing.T) {
		var id1, id2 uintptr
		await(func() { id1 = goid() })
		await(func() { id2 = goid() })
		if id1 == id2 {
			t.Errorf("expected different goroutine ids, got %d and %d", id1, id2)
		}
	})

	t.Run("when unable to get stack", func(t *testing.T) {
		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				if !errors.Is(err, ErrNoStack) {
					t.Errorf("expected panic with ErrNoStack, got: %v", err)
				}
			} else {
				t.Errorf("expected panic with error, got: %v", r)
			}
		}()

		og := stack
		defer func() { stack = og }()
		stack = func([]byte, bool) int { return 0 } // simulate failure to parse stack

		_ = goid()
	})

	t.Run("when unable to parse stack", func(t *testing.T) {
		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				if !errors.Is(err, ErrUnexpectedStackFormat) {
					t.Errorf("expected panic with ErrUnexpectedStackFormat, got: %v", err)
				}
			} else {
				t.Errorf("expected panic with error, got: %v", r)
			}
		}()

		og := stack
		defer func() { stack = og }()
		stack = func([]byte, bool) int { return 1 } // simulate failure to parse stack

		_ = goid()
	})
}
