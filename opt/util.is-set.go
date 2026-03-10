package opt

import "slices"

func IsSet(opts []any, opt any) bool {
	return slices.Contains(opts, opt)
}
