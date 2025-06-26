package matcher

type ForAny interface {
	Match(got any, opts ...any) bool
}

type ForType[T any] interface {
	Match(got T, opts ...any) bool
}
