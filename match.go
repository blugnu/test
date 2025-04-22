package test

type Matcher[T any] interface {
	Match(T, ...any) bool
}
