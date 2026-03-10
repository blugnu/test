package test

import "fmt"

type Outcome int

const (
	Passed Outcome = iota
	Failed
	Panicked
)

func (to Outcome) String() string {
	switch to {
	case Passed:
		return "test.Passed"
	case Failed:
		return "test.Failed"
	case Panicked:
		return "test.Panicked"
	default:
		return fmt.Sprintf("test.Outcome(%d)", to)
	}
}
