package testcase

type Controller[T any] struct {
	// idx is the index of the test case in the list of test cases
	idx int

	// name is the name of the test case, used to identify it in the test output
	// and in any Before/After scaffolding functions
	name string

	// data is the test case data
	data T

	// debug is used to indicate whether the test case is being debugged
	debug bool

	// skip is used to indicate whether the test case should be skipped
	skip bool

	// parallel is used to indicate whether the test case should be run in parallel
	parallel bool
}

func NewController[T any](data T, index int, name string) Controller[T] {
	return Controller[T]{
		idx:   index,
		name:  NameOrDefault(data, name, index),
		data:  data,
		debug: IsDebugging(data),
		skip:  IsSkipping(data),
	}
}
