package opt

func IsSet(opts []any, opt any) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}
