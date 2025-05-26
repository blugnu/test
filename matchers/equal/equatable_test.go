package equal_test

type equatable struct {
	result bool
}

// Equal implements the Equal method for the equable type. This implementation
// always equals true, regardless of the value of the other equable instance,
// which may be used to verify that the Equal method is called and that the
// comparison has not been made using == or reflect.DeepEqual.
func (e equatable) Equal(other equatable) bool {
	return true
}
