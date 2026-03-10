package report

// Empty returns a string representation of an empty value, optionally including a type
// name for clarity. This can be used in failure reports to indicate that a value was
// expected to be empty.
//
// If a type name is provided, the returned string will be in the format "<empty typeName>".
// If no type name is provided, the returned string will simply be "<empty>".
func Empty(typeName ...string) string {
	if len(typeName) > 0 {
		return "<empty " + typeName[0] + ">"
	}

	return "<empty>"
}
