package opt

import "fmt"

// IntervalClosure defines the closure types interval options used in
// ordered comparisons:
//
//	Closed    // [min, max] or: min <= x <= max
//	Open      // (min, max) or: min < x < max
//	OpenMin   // (min, max] or: min < x <= max
//	OpenMax   // [min, max) or: min <= x < max
type IntervalClosure int

const (
	// IntervalClosed means the interval is closed on both ends: [min, max]: min <= x <= max
	IntervalClosed IntervalClosure = iota

	// IntervalOpenMin means the interval is open on the minimum end: (min, max]: min < x <= max
	IntervalOpenMin

	// IntervalOpenMax means the interval is open on the maximum end: [min, max): min <= x < max
	IntervalOpenMax

	// IntervalOpen means the interval is open on both ends: (min, max): min < x < max
	IntervalOpen
)

func (i IntervalClosure) String() string {
	switch i {
	case IntervalClosed:
		return "IntervalClosed: min <= x <= max"
	case IntervalOpen:
		return "IntervalOpen: min < x < max"
	case IntervalOpenMin:
		return "IntervalOpenMin: min < x <= max"
	case IntervalOpenMax:
		return "IntervalOpenMax: min <= x < max"
	}

	return fmt.Sprintf("IntervalClosure(%d)", i)
}
