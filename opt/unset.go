package opt

// Unset removes a specific option from a slice of options
func Unset(opts []any, opt any) []any {
	result := make([]any, 0, len(opts))

	for _, o := range opts {
		if o == opt {
			continue
		}
		result = append(result, o)
	}

	return result
}
