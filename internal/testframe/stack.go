package testframe

import (
	"fmt"
	"runtime"
	"sync"
)

// stack is a variable that holds the runtime.Stack function.
// It can be replaced in tests to simulate different stack trace behaviors.
var stack = runtime.Stack

// goid gets the id of the calling goroutine using runtime.Stack
//
// The function relies on the format of the stack trace returned by runtime.Stack,
// expecting the first line to contain "goroutine <id> ", where <id> is the goroutine id.
//
// If the stack trace does not match this format, it will panic with ErrUnexpectedStackFormat.
func goid() uintptr {
	const maxFrames = 64

	buf := make([]byte, maxFrames)
	n := stack(buf, false)
	if n <= 0 {
		panic(ErrNoStack)
	}

	var id uintptr
	if _, err := fmt.Sscanf(string(buf[:n]), "goroutine %d", &id); err != nil {
		panic(fmt.Errorf("%w: %w", ErrUnexpectedStackFormat, err))
	}
	return id
}

type testframe struct {
	T any
	// ref string
}

var stacks = struct {
	sync.RWMutex
	frames map[uintptr][]testframe
}{
	frames: make(map[uintptr][]testframe),
}

// MustPeek retrieves the top element of the current test frame stack without
// removing it.
//
// The function will panic with ErrNoTestFrame if the top element cannot be type
// asserted to T or if the stack is empty.
//
// The function is safe to call concurrently and will not modify the stack.
func MustPeek[T any]() T {
	if t, ok := Peek[T](); ok {
		return t
	}
	panic(ErrNoTestFrame)
}

// Peek retrieves the top element of the current test frame stack without removing it.
//
// The function checks if the top element can be type asserted to T and returns it with
// true if successful. If the stack is empty or the type assertion fails, a zero-value T
// is returned with false.
//
// The function is safe to call concurrently and will not modify the stack.
func Peek[T any]() (T, bool) {
	stacks.RLock()
	defer stacks.RUnlock()

	id := goid()
	stk := stacks.frames[id]
	if len(stk) == 0 {
		return *new(T), false
	}

	f := &stk[len(stk)-1]
	switch t := f.T.(type) {
	case Nil:
		return *new(T), false
	case T:
		return t, true
	default:
		return *new(T), false
	}
}

// Pop removes the top element from the current test frame stack. The element
// is not returned.  If you need to access the top element before removing it,
// use Peek() before Pop().
//
// If the stack is empty the function will panic with ErrEmptyStack.
//
// The function is safe to call concurrently.
func Pop() {
	stacks.Lock()
	defer stacks.Unlock()

	id := goid()
	stk := stacks.frames[id]

	if len(stk) == 0 {
		panic(ErrEmptyStack)
	}

	stk = stk[:len(stk)-1]
	stacks.frames[id] = stk
}

// Push adds a new test frame to the current goroutine's stack.
//
// The function takes a single argument of any type, which is stored in the
// testframe struct. The function retrieves the current goroutine's id using
// goid() and appends the new test frame to the stack associated with that id.
//
// The function is safe to call concurrently and will not modify the stack
// of other goroutines.
func Push(t any) {
	stacks.Lock()
	defer stacks.Unlock()

	id := goid()
	stk := stacks.frames[id]
	stk = append(stk, testframe{T: t})
	stacks.frames[id] = stk
}
