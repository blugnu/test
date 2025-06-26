package test //nolint: testpackage // examples for test package reference the test package in expected output

import (
	"context"

	"github.com/blugnu/test/test"
)

func ExampleHaveContextKey() {
	test.Example()

	type key int
	ctx := context.WithValue(context.Background(), key(57), "varieties")

	// these tests will pass
	Expect(ctx).To(HaveContextKey(key(57)))
	Expect(ctx).ToNot(HaveContextKey(key(58)))

	// this test will fail
	Expect(ctx).To(HaveContextKey(key(58)))

	// Output:
	// expected key: test.key(58)
	//   key not present in context
}

func ExampleHaveContextValue() {
	// this is needed to make the example work; this would be usually
	// be `With(t)` where `t` is the *testing.T
	test.Example()

	type key int
	ctx := context.WithValue(context.Background(), key(57), "varieties")

	// these tests will pass
	Expect(ctx).To(HaveContextValue(key(57), "varieties"))
	Expect(ctx).ToNot(HaveContextValue(key(56), "varieties"))
	Expect(ctx).ToNot(HaveContextValue(key(57), "flavours"))

	// this test will fail
	Expect(ctx).To(HaveContextValue(key(57), "flavours"))

	// Output:
	// context value: test.key(57)
	//   expected: "flavours"
	//   got     : "varieties"
}
