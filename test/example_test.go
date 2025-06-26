package test_test

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

func TestExample(t *testing.T) {
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

	runExample := func(fn func(TestingT)) (TestingT, []string, []string) {
		test.Example()
		defer testframe.Pop()

		// execute the passed function with the T(), which is the
		// ExampleTestRunner we pushed with test.Example()
		t := T()
		out, log := Record(func() { fn(t) })
		return t, out, log
	}

	Run("no-ops coverage", func() {
		_, out, log := runExample(func(sut TestingT) {
			sut.Cleanup(nil)
			sut.Helper()
			sut.Parallel()
			sut.Run("no-op", nil)
			sut.Setenv("no-op", "no-op")
		})

		Expect(out).IsNil()
		Expect(log).IsNil()
	})

	Run("name", func() {
		t, out, log := runExample(func(sut TestingT) {})
		Expect(out).IsNil()
		Expect(log).IsNil()
		Expect(t.Name()).To(Equal("ExampleTestRunner"))
	})

	Run("running a subtest", func() {
		_, out, log := runExample(func(sut TestingT) {
			sut.Run("subtest", func(t *testing.T) {
				fmt.Println("sub test ran OK")
			})
		})
		Expect(out).To(EqualSlice([]string{
			"sub test ran OK",
		}))
		Expect(log).IsNil()
	})

	Run("skipping after an error", func() {
		t, out, err := runExample(func(sut TestingT) {
			await(func() {
				sut.Error("first error")
				sut.SkipNow()
				sut.Error("second error")
			})
		})

		Expect(t.Failed()).To(BeTrue())
		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"first error",
		}))
	})

	Run("skipping before any error", func() {
		t, out, err := runExample(func(sut TestingT) {
			await(func() {
				sut.SkipNow()
				sut.Error("first error")
			})
		})

		Expect(err).IsNil()
		Expect(out).Should(BeEmptyOrNil())
		Expect(t.Failed()).To(BeFalse())
	})

	Run("failing", func() {
		t, out, err := runExample(func(sut TestingT) {
			await(func() {
				sut.Fail()
			})
		})

		Expect(err).IsNil()
		Expect(out).Should(BeEmptyOrNil())
		Expect(t.Failed()).To(BeTrue())
	})

	Run("non-fatal errors", func() {
		_, out, err := runExample(func(sut TestingT) {
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
		_, out, err := runExample(func(sut TestingT) {
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
		_, out, err := runExample(func(sut TestingT) {
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
}
