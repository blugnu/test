package test

import (
	"sync"
	"testing"
)

func await(f func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	wg.Wait()
}

func TestT(t *testing.T) {
	With(t)

	Run("GetT", func() {
		Run("when test frame set", func() {
			t1 := GetT()
			Expect(t1).IsNotNil()
			Expect(t1.Name()).To(Equal("TestT/GetT/when_test_frame_set"))
		})

		Run("when test frame not set", func() {
			defer Expect(Panic(ErrNoTestFrame)).DidOccur()

			With(nil)
			_ = GetT()
		})
	})

	Run("T", func() {
		Run("when test frame set", func() {
			t1 := T()
			Expect(t1).IsNotNil()
			Expect(t1.Name()).To(Equal("TestT/T/when_test_frame_set"))
		})

		Run("when test frame not set", func() {
			defer Expect(Panic(ErrNoTestFrame)).DidOccur()

			With(nil)
			_ = T()
		})
	})
}

func TestExampleTestRunner(t *testing.T) {
	With(t)

	sut := ExampleTestRunner{}

	Run("no-ops coverage", func() {
		out, log := Record(func() {
			sut.Cleanup(nil)
			sut.Fail()
			sut.Helper()
			sut.Parallel()
			sut.Run("no-op", nil)
			sut.Setenv("no-op", "no-op")
			sut.SkipNow()
		})
		Expect(out).IsNil()
		Expect(log).IsNil()
	})

	Run("non-fatal errors", func() {
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
		out, err := Record(func() {
			await(func() {
				sut.Fatal("fatal error")
				sut.Error("this is not reached")
			})
		})

		Expect(err).IsNil()
		Expect(out).To(EqualSlice([]string{
			"fatal error",
		}))
	})

	Run("fatalf error", func() {
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
}
