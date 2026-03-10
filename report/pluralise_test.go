package report_test

import (
	"testing"

	"github.com/blugnu/test/report"
)

func TestPluralise(t *testing.T) {
	t.Run("singular", func(t *testing.T) {
		if got := report.Pluralise("key", 1); got != "key" {
			t.Errorf("Pluralise(\"key\", 1) = %q, want \"key\"", got)
		}
	})

	t.Run("default plural", func(t *testing.T) {
		if got := report.Pluralise("key", 0); got != "keys" {
			t.Errorf("Pluralise(\"key\", 0) = %q, want \"keys\"", got)
		}
		if got := report.Pluralise("key", 2); got != "keys" {
			t.Errorf("Pluralise(\"key\", 2) = %q, want \"keys\"", got)
		}
	})

	t.Run("custom plural", func(t *testing.T) {
		if got := report.Pluralise("child", 2, "children"); got != "children" {
			t.Errorf("Pluralise(\"child\", 2, \"children\") = %q, want \"children\"", got)
		}
	})

	t.Run("too many plurals", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Pluralise did not panic when given too many plural forms")
			}
		}()
		report.Pluralise("key", 2, "keys", "keyses")
	})
}
