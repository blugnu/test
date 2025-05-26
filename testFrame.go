package test

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"sync"
)

var stack = runtime.Stack

// goid gets the id of the calling goroutine using runtime.Stack
func goid() uintptr {
	buf := make([]byte, 64)
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
	TestingT
	ref string
}

type testframes struct {
	sync.RWMutex
	frames map[uintptr][]testframe
}

func (tf *testframes) peek() *testframe {
	tf.RLock()
	id := goid()
	stk := tf.frames[id]
	tf.RUnlock()

	if len(stk) == 0 {
		return nil
	}
	return &stk[len(stk)-1]
}

func (tf *testframes) pop() {
	tf.Lock()
	defer tf.Unlock()

	id := goid()
	stk := tf.frames[id]
	stk = stk[:len(stk)-1]
	tf.frames[id] = stk
}

func (tf *testframes) push(t TestingT) {
	tf.Lock()
	defer tf.Unlock()

	id := goid()
	stk := tf.frames[id]
	stk = append(stk, testframe{TestingT: t, ref: CallerFuncName(2)})
	tf.frames[id] = stk
}

var testFrames = testframes{frames: map[uintptr][]testframe{}}

func CallerFilename() string {
	pcs := make([]uintptr, 1)
	n := runtime.Callers(3, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)

	frame, _ := frames.Next()
	_, filename := path.Split(frame.File)
	return filename
}

func CallerFuncName(skip int, maxCallers ...int) string {
	if len(maxCallers) == 0 {
		maxCallers = []int{20}
	}

	pcs := make([]uintptr, maxCallers[0])
	n := runtime.Callers(2+skip, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)

	_, _ = frames.Next()
	frame, _ := frames.Next()

	el := strings.Split(frame.Func.Name(), "/")
	return el[len(el)-1]
}

func TestFrame(skip int) TestingT {
	if f := testFrames.peek(); f != nil {
		return f.TestingT
	}
	return nil
}

func With(t TestingT) {
	testFrames.push(t)
	if t != nil {
		t.Cleanup(func() {
			testFrames.pop()
		})
	}
}
