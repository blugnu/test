package test

import (
	"os"
	"testing"
)

func TestCaptureOutput(t *testing.T) {
	// ARRANGE
	writeOutput := func() {
		os.Stdout.WriteString("to stdout (1)\n")
		os.Stdout.WriteString("to stdout (2)\n")
		os.Stderr.WriteString("to stderr (1)\n")
		os.Stderr.WriteString("to stderr (2)\n")
	}

	// ACT
	stdout, stderr := CaptureOutput(t, func(t *testing.T) {
		writeOutput()
	})

	// ASSERT
	Equal(t, "to stdout (1)", stdout.s[0])
	Equal(t, "to stdout (2)", stdout.s[1])
	Equal(t, "to stderr (1)", stderr.s[0])
	Equal(t, "to stderr (2)", stderr.s[1])
}

func TestCapturedOutput(t *testing.T) {
	// ARRANGE
	sut := CapturedOutput{"test", []string{}}

	testcases := []struct {
		name    string
		sut     []string
		fn      func(*testing.T)
		outcome any
		output  any
	}{
		{name: "IsEmpty, no output", fn: func(t *testing.T) { sut.IsEmpty(t) }, outcome: ShouldPass},
		{name: "IsEmpty, with output", sut: []string{"some output"}, fn: func(t *testing.T) { sut.IsEmpty(t) },
			outcome: ShouldFail,
			output: []string{
				"wanted: (no output)",
				"got   : BEGIN CAPTURE [test]",
				"  some output",
				"--END CAPTURE",
			},
		},
		{name: "Contains with invalid arg", fn: func(t *testing.T) { sut.Contains(t, 42) },
			outcome: ExpectPanic(ErrInvalidArgument),
		},
		{name: "Contains no output, want nil", fn: func(t *testing.T) { sut.Contains(t, nil) },
			outcome: ShouldPass,
		},
		{name: "Contains no output, want string", fn: func(t *testing.T) { sut.Contains(t, "some output") },
			outcome: ShouldFail,
			output: []string{
				"wanted: \"some output\"",
				"got   : (no output)",
			},
		},
		{name: "Contains no output, want strings (0)", fn: func(t *testing.T) { sut.Contains(t, []string{}) },
			outcome: ShouldPass,
		},
		{name: "Contains no output, want strings (1)", fn: func(t *testing.T) { sut.Contains(t, []string{"some output"}) },
			outcome: ShouldFail,
			output: []string{
				"wanted: \"some output\"",
				"got   : (no output)",
			},
		},
		{name: "Contains no output, want strings (>1)", fn: func(t *testing.T) { sut.Contains(t, []string{"some output", "more output"}) },
			outcome: ShouldFail,
			output: []string{
				"wanted:",
				"some output",
				"more output",
				"got   : (no output)",
			},
		},
		{name: "Contains output, want empty", sut: []string{"some output"}, fn: func(t *testing.T) { sut.Contains(t, "") },
			outcome: ShouldFail,
			output: []string{
				"wanted: (no output)",
				"got   : BEGIN CAPTURE [test]",
				"  some output",
				"--END CAPTURE",
			},
		},
		{name: "Contains output with wanted string", sut: []string{"   line 1"}, fn: func(t *testing.T) { sut.Contains(t, "line 1") },
			outcome: ShouldPass,
		},
		{name: "Contains output with wanted strings", sut: []string{"line 1", "line 2"}, fn: func(t *testing.T) { sut.Contains(t, []string{"line 1", "line 2"}) },
			outcome: ShouldPass,
		},
		{name: "Contains output with wanted strings having empty lines", sut: []string{"line 1", "  ", "line 2"}, fn: func(t *testing.T) { sut.Contains(t, []string{"line 1", "", "line 2"}) },
			outcome: ShouldPass,
		},
		{name: "Contains output missing string", sut: []string{"line 1", "line 2"}, fn: func(t *testing.T) { sut.Contains(t, "line 3") },
			outcome: ShouldFail,
			output: []string{
				"wanted: \"line 3\"",
				"got   : BEGIN CAPTURE [test]",
				"  line 1",
				"  line 2",
				"--END CAPTURE",
			},
		},
		{name: "Contains output missing empty line in wanted", sut: []string{"line 1", "line 2"}, fn: func(t *testing.T) { sut.Contains(t, []string{"line 1", "", "line 2"}) },
			outcome: ShouldFail,
			output: []string{
				"wanted:",
				"line 1",
				"",
				"line 2",
				"got   : BEGIN CAPTURE [test]",
				"  line 1",
				"  line 2",
				"--END CAPTURE",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			sut.s = tc.sut

			// ACT
			stdout, _ := Helper(t, func(st *testing.T) {
				tc.fn(st)
			}, tc.outcome)

			// HERE BE DRAGONS: Contains is used to test itself!

			// ASSERT
			stdout.Contains(t, tc.output)
		})
	}
}
