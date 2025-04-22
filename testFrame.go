package test

import (
	"runtime"
)

type frameCache map[uintptr]TestRunner

func (fc frameCache) Get(id uintptr) (TestRunner, bool) {
	t, ok := fc[id]
	return t, ok
}

func (fc frameCache) Remove(id uintptr) {
	delete(fc, id)
}

func (fc frameCache) Set(id uintptr, t TestRunner) {
	fc[id] = t
}

var testFrames = frameCache{}

func CallerFrameId(maxCallers ...int) uintptr {
	if len(maxCallers) == 0 {
		maxCallers = []int{20}
	}

	pcs := make([]uintptr, maxCallers[0])
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	if _, more := frames.Next(); !more {
		return 0
	}

	frame, _ := frames.Next()
	return frame.Entry
}

func TestFrame(maxCallers ...int) (TestRunner, uintptr) {
	if len(maxCallers) == 0 {
		maxCallers = []int{20}
	}

	pcs := make([]uintptr, maxCallers[0])
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		t, ok := testFrames.Get(frame.Entry)
		switch {
		case !more:
			return nil, 0
		case ok:
			return t, frame.Entry
		}
	}
}

func With(t TestRunner) uintptr {
	frameId := CallerFrameId()
	if frameId == 0 {
		panic("With: unable to get frame ID")
	}
	testFrames.Set(frameId, t)
	return frameId
}
