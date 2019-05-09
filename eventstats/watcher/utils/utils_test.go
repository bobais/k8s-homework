package utils

import (
	"testing"
)

func TestIsPositiveStatement(t *testing.T) {
	positive := []string{"t", "true", "True", "y", "yes", "Yes", "1"}
	negative := []string{"f", "false", "False", "n", "no", "No", "0"}
	for _, value := range positive {
		got := IsPositiveStatement(value)
		if !IsPositiveStatement(value) {
			t.Errorf("Result for %s got: %t, want: %t.", value, got, true)
		}
	}

	for _, value := range negative {
		got := IsPositiveStatement(value)
		if IsPositiveStatement(value) {
			t.Errorf("Result for %s got: %t, want: %t.", value, got, false)
		}
	}
}
