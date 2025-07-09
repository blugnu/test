package test_test

import (
	"fmt"
	"testing"

	. "github.com/blugnu/test"
)

func TestCleanup(t *testing.T) {
	With(t)

	Run(
		Test("is a no-op when called with a nil func", func() {
			// this is a coverage test for the nil func case; there is no
			// way to verify anything other than that it (implicitly) does
			// not panic
			Cleanup(nil)
		}))

	Run(Test("runs cleanup functions at the end of a test", func() {
		out, err := Record(func() {
			Run(Test("subtest 1", func() {
				Cleanup(func() {
					fmt.Println("cleanup 1")
				})
			}))

			Run(Test("subtest 2", func() {
				Cleanup(func() {
					fmt.Println("cleanup 2")
				})
			}))
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"cleanup 1",
			"cleanup 2",
		}))
	}),
	)
}
