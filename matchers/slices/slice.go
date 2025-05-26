package slices

type slice[T any] []T

func (s slice[T]) appendToTestReport(r []string, p string, opts ...any) []string {
	return AppendToReport(r, s, p, opts...)
}

// containsItems returns true if the receiver contains all of the elements
// in the target slice.  The order of the elements in the target slice
// does not matter.  The target slice may contain duplicates; the receiver
// matches if the receiver contains at least as many duplicates of any
// duplicated elements in the target slice.
func (s slice[T]) containsItems(target []T, cmp func(any, any) bool) bool {
	// make a copy of the target slice
	cpy := make([]T, len(target))
	copy(cpy, target)

	// check if the receiver contains all elements in the target slice
	for _, item := range s {
		for i, target := range cpy {
			if cmp(item, target) {
				cpy = append(cpy[:i], cpy[i+1:]...)
				break
			}
		}
	}

	return len(cpy) == 0
}

// containsSlice returns true if the receiver contains at least one occurence
// of all elements in the target slice in the same order as the target slice.
func (s slice[T]) containsSlice(target []T, cmp func(any, any) bool) bool {
	ix := 0
	for {
		if ix = s.next(target[0], ix, cmp); ix == -1 {
			return false
		}
		if ix+len(target) > len(s) {
			return false
		}
		for j, target := range target {
			if !cmp(s[ix+j], target) {
				goto nextOccurence
			}
		}
		return true
	nextOccurence:
		ix++
	}
}

func (s slice[T]) next(target T, startAt int, cmp func(any, any) bool) int {
	// startAt is the index of the first element to check
	// cmp is the comparison function
	// target is the element to find
	// return -1 if not found
	for i := startAt; i < len(s); i++ {
		if cmp(s[i], target) {
			return i
		}
	}
	return -1
}
