package assert

import "testing"

func Equal[T comparable](t *testing.T, expected, actual T) {
	if expected != actual {
		t.Errorf("Expected: %v, Actual: %v", expected, actual)
	}
}
