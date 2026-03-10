package report

// Pluralise returns singular if n is 1, otherwise plural[0] or singular + "s"
// if no plural form is provided.
//
// panics if more than one plural form is provided.
func Pluralise(singular string, n int, plural ...string) string {
	switch {
	case n == 1:
		return singular
	case len(plural) == 0:
		return singular + "s"
	case len(plural) == 1:
		return plural[0]
	default:
		panic("Pluralise: too many plural forms provided (at most one allowed)")
	}
}
