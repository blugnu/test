package test

import "testing"

func TestType(t *testing.T) {
	// ARRANGE
	type returns struct {
		string
		bool
	}
	testcases := []struct {
		name string
		got  any
		returns
		outcome HelperResult
		output  any
	}{
		{name: "string", got: "foo", returns: returns{"foo", true},
			outcome: ShouldPass,
		},
		{name: "int", got: 42,
			outcome: ShouldFail,
			output: []string{
				"wanted: string",
				"got   : int",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			var (
				got string
				ok  bool
			)
			stdout, _ := Helper(t, func(st *testing.T) {
				got, ok = Type[string](st, tc.got)
			}, tc.outcome)

			// ASSERT
			Equal(t, tc.returns, returns{got, ok})
			stdout.Contains(t, tc.output)
		})
	}
}
