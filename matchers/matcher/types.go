package matcher

type ForAny interface {
	Match(any, ...any) bool
}

type ForType[T any] interface {
	Match(T, ...any) bool
}
