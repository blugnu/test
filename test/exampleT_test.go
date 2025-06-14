package test_test

import (
	"sync"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestExampleT(t *testing.T) {
	With(t)

	await := func(f func()) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
		wg.Wait()
	}

	Run("no-ops coverage", func() {
		sut := test.ExampleT()
		out, log := Record(func() {
			sut.Cleanup(nil)
			sut.Helper()
			sut.Parallel()
			sut.Run("no-op", nil)
			sut.Setenv("no-op", "no-op")

			_ = sut.Name()
		})
		Expect(out).IsNil()
		Expect(log).IsNil()
	})

	Run("skipping after an error", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.Error("first error")
				sut.SkipNow()
				sut.Error("second error")
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"first error",
		}))

		Expect(sut.Failed()).To(BeTrue())
	})

	Run("skipping before any error", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.SkipNow()
				sut.Error("first error")
			})
		})

		Expect(err).IsNil()
		Expect(out).Should(BeEmptyOrNil())

		Expect(sut.Failed()).To(BeFalse())
	})

	Run("failing", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.Fail()
			})
		})

		Expect(err).IsNil()
		Expect(out).Should(BeEmptyOrNil())

		Expect(sut.Failed()).To(BeTrue())
	})

	Run("non-fatal errors", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.Error("error")
				sut.Errorf("errorf %s", "formatted")
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"error",
			"errorf formatted",
		}))
	})

	Run("fatal error", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.Fatal("fatal error")
				sut.Fatal("this is not reached")
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"fatal error",
		}))
	})

	Run("fatalf error", func() {
		sut := test.ExampleT()
		out, err := Record(func() {
			await(func() {
				sut.Fatalf("fatal error %s", "formatted")
				sut.Error("this is not reached")
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"fatal error formatted",
		}))
	})

	Run("failed", func() {
		sut := test.ExampleT()
		Expect(sut.Failed()).To(BeFalse())

		sut.Error()
		Expect(sut.Failed()).To(BeTrue())
	})
}
