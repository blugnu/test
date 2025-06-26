package testframe //nolint: testpackage // tests rely on access to private functions

import (
	"errors"
	"testing"
)

func TestMustPeek(t *testing.T) {
	og := stacks.frames
	defer func() {
		stacks.frames = og
	}()

	t.Run("empty stack", func(t *testing.T) {
		// ensure we start with no stacks
		stacks.frames = map[uintptr][]testframe{}

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic, got nil")
			} else if err, ok := r.(error); !ok || !errors.Is(err, ErrNoTestFrame) {
				t.Errorf("expected panic with ErrNoTestFrame, got: %v", r)
			}
		}()

		MustPeek[*testing.T]()
	})

	t.Run("valid type", func(t *testing.T) {
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: t}},
		}

		result := MustPeek[*testing.T]()

		if result != t {
			t.Errorf("expected to peek item from stack got %v", result)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: t}},
		}

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic, got nil")
			} else if err, ok := r.(error); !ok || !errors.Is(err, ErrNoTestFrame) {
				t.Errorf("expected panic with ErrNoTestFrame, got: %v", r)
			}
		}()

		MustPeek[string]()
	})
}

func TestPeek(t *testing.T) {
	og := stacks.frames
	defer func() {
		stacks.frames = og
	}()

	t.Run("empty stack", func(t *testing.T) {
		// ensure we start with no stacks
		stacks.frames = map[uintptr][]testframe{}

		result, ok := Peek[*testing.T]()

		if result != nil {
			t.Errorf("expected nil test frame, got %v", result)
		}

		if ok != false {
			t.Errorf("expected false, got true")
		}
	})

	// for these tests we need a valid item in the stack; each stack can
	// hold items of any type, so we use an int for simplicity
	//
	// t.Run runs each sub-test in its own goroutine, so we need to init
	// the stack in each sub-test

	t.Run("valid type", func(t *testing.T) {
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: 42}},
		}

		result, ok := Peek[int]()

		if ok != true {
			t.Error("expected true, got false")
		}

		if result != 42 {
			t.Errorf("expected to peek item from stack got %v", result)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: 42}},
		}

		result, ok := Peek[string]()

		if result != "" {
			t.Errorf("expected 0, got %v", result)
		}

		if ok != false {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("Nil sentinel", func(t *testing.T) {
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: Nil{}}},
		}

		result, ok := Peek[*testing.T]()

		if result != nil {
			t.Errorf("expected nil test frame, got %v", result)
		}

		if ok != false {
			t.Errorf("expected false, got true")
		}
	})
}

func TestPop(t *testing.T) {
	og := stacks.frames
	defer func() {
		stacks.frames = og
	}()

	t.Run("empty stack", func(t *testing.T) {
		// ensure we start with no stacks
		stacks.frames = map[uintptr][]testframe{}

		// ensure we panic if we try to pop without a stack
		defer func() {
			r := recover()
			if err, ok := r.(error); ok {
				if !errors.Is(err, ErrEmptyStack) {
					t.Errorf("expected panic with ErrEmptyStack, got: %v", err)
				}
				return // we got the expected panic
			} else if r != nil {
				t.Errorf("expected panic with error, got: %v", r)
				return
			}
			t.Errorf("expected panic, got nil")
		}()

		Pop()
	})

	t.Run("pop from stack", func(t *testing.T) {
		// ensure we start with an item in the stack
		stacks.frames = map[uintptr][]testframe{
			goid(): {{T: t}},
		}
		// ensure we can pop the item from the stack
		Pop()
		if len(stacks.frames[goid()]) != 0 {
			t.Errorf("expected stack to be empty after pop, got %d items", len(stacks.frames[goid()]))
		}
	})
}

func TestPush(t *testing.T) {
	og := stacks.frames
	defer func() {
		stacks.frames = og
	}()

	// ensure we start with no stacks
	stacks.frames = map[uintptr][]testframe{}

	Push(t)
	if len(stacks.frames) != 1 {
		t.Errorf("expected 1 stack, got %d", len(stacks.frames))
	}
}
