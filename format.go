package test

type CustomOneLineReport interface {
	Format() string
}
type CustomReport interface {
	Format() []string
}
type CustomFormatExpectedAndGot interface {
	Expected() any
	Format(any, any) []string
}
type Expected interface {
	Expected() any
}
type FormatExpectedAndGot interface {
	Expected() any
	Format(any) string
}
type OneLineExpected interface {
	Expected() any
	OneLineError()
}

type CustomOneLineReportFunc func() string

func (f CustomOneLineReportFunc) Format() string {
	return f()
}

func getFormatter(opts ...any) any {
	for _, opt := range opts {
		switch opt := opt.(type) {
		case
			CustomOneLineReport,
			CustomReport,
			CustomFormatExpectedAndGot,
			Expected,
			FormatExpectedAndGot,
			OneLineExpected:
			return opt
		default:
			continue
		}
	}
	return nil
}
