package jsonx

import "testing"

func TestMarshal(t *testing.T) {
	t.Run("should generate null when data is nil", func(t *testing.T) {
		expected := []byte(`null`)
		actual, err := Marshal(nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(actual) != string(expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}
	})
}
