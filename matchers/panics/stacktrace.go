package panics

import (
	"strings"

	"github.com/blugnu/test/opt"
)

func StackTrace(raw []byte, opts ...any) []string {
	if len(raw) == 0 || opt.IsSet(opts, opt.StackTrace(false)) {
		return nil
	}

	trace := strings.Split(string(raw), "\n")
	result := make([]string, 0, len(trace))

	panicStart := 0
	for i, line := range trace {
		if strings.HasPrefix(strings.TrimSpace(line), "panic(") {
			panicStart = i
			break
		}
	}

	const skipPanicLines = 2
	for i := panicStart + skipPanicLines; i < len(trace); i++ {
		if line := trace[i]; line != "" {
			result = append(result, "  "+line)
		}
	}

	return result
}
