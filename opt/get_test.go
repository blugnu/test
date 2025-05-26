package opt

import "testing"

func TestGet(t *testing.T) {
	var (
		result string
		ok     bool
	)

	t.Run("Get[string] when options contain a string", func(t *testing.T) {
		result, ok = Get[string]([]any{0, 1, "test", true})
		if !ok {
			t.Errorf("ok: expected true, got false")
		}
		if result != "test" {
			t.Errorf("result: expected 'test', got '%s'", result)
		}
	})

	t.Run("Get[string] when options does not contain a string", func(t *testing.T) {
		result, ok = Get[string]([]any{0, 1, true})
		if ok {
			t.Errorf("ok: expected false, got true")
		}
		if result != "" {
			t.Errorf("result: expected '', got '%s'", result)
		}
	})

	t.Run("Get[string] when options contain a value of custom string type", func(t *testing.T) {
		type CustomStringType string

		result, ok = Get[string]([]any{0, 1, CustomStringType("test"), true})
		if ok {
			t.Errorf("ok: expected false, got true")
		}
		if result != "" {
			t.Errorf("result: expected '', got '%s'", result)
		}
	})
}
