package test_test

import (
	"fmt"
	"testing"

	. "github.com/blugnu/test"
)

func TestCleanup(t *testing.T) {
	With(t)

	Run("is a no-op when called with a nil func", func() {
		Cleanup(nil)
	})

	Run("runs cleanup functions at the end of a test", func() {
		out, err := Record(func() {
			Run("subtest 1", func() {
				Cleanup(func() {
					fmt.Println("cleanup 1")
				})
			})
			Run("subtest 2", func() {
				Cleanup(func() {
					fmt.Println("cleanup 2")
				})
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"cleanup 1",
			"cleanup 2",
		}))
	})
}
