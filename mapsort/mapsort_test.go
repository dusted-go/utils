package mapsort

import (
	"strconv"
	"testing"
)

func Test_Keys(t *testing.T) {
	m := map[string]int{"B": 2, "a": 1, "D": 4, "c": 3}
	actual := Keys(m)
	expected := []string{"B", "D", "a", "c"}

	for i, item := range expected {
		if actual[i] != item {
			t.Errorf("Expected: %v, Actual: %v", expected, actual)
			return
		}
	}
}

type Foo struct {
	Value int
}

func Test_KeysByValue(t *testing.T) {
	m := map[string]Foo{"B": {Value: 2}, "a": {Value: 1}, "D": {Value: 4}, "c": {Value: 3}}
	actual := KeysByValue(m, func(foo Foo) string { return strconv.Itoa(foo.Value) })
	expected := []string{"a", "B", "c", "D"}

	for i, item := range expected {
		if actual[i] != item {
			t.Errorf("Expected: %v, Actual: %v", expected, actual)
			return
		}
	}
}
